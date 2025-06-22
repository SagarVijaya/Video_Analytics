package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	mailsender "videoanalytics/apps/mailSender"
	"videoanalytics/apps/metrics"
	"videoanalytics/apps/models"
	"videoanalytics/apps/utils"
	"videoanalytics/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

var Gredis *redis.Client
var Gctx = context.Background()

func StoreClickInfo(r *gin.Context) {
	log := new(utils.LoggerId)
	log.SetSid(r.Request)
	log.Log("StoreClickInfo (+)")
	defer log.Log("StoreClickInfo (-)")

	metrics.HttpRequests.WithLabelValues("/ads/click", "POST").Inc()

	var lRequestData models.AdClickRate
	lErr := r.BindJSON(&lRequestData)
	if lErr != nil {
		metrics.RequestErrors.WithLabelValues("/ads/click").Inc()
		log.Log("StoreClickInfo-001", lErr.Error())
		r.JSON(http.StatusBadRequest, gin.H{
			"status": "E",
			"errMsg": "Invalid Details",
		})
		return
	}

	// 	clickChan := make(chan models.AdClickRate, 1000)

	// go func() {
	// 	for click := range clickChan {
	// 		InsertClickDetails(log, click) // same logic with Redis fallback
	// 	}
	// }()
	go InsertClickDetails(log, lRequestData)

	r.Writer.Header().Set("Content-Type", "application/json")
	r.JSON(http.StatusOK, gin.H{
		"status": "S",
		"errMsg": "",
	})
}

func InsertClickDetails(pLogDetails *utils.LoggerId, pClickInfo models.AdClickRate) {
	pLogDetails.Log("InsertClickDetails (+)")
	defer pLogDetails.Log("InsertClickDetails (-)")
	lErr := database.DB.Table("ad_click_rates").Clauses(clause.OnConflict{DoNothing: true}).Create(&pClickInfo)
	if lErr.Error != nil {
		// Alert Mail Trigger
		go mailsender.SendFailureAlert(pClickInfo, lErr.Error)

		metrics.AdClicks.WithLabelValues(strconv.Itoa(int(pClickInfo.AdID))).Inc()
		pLogDetails.Log("InsertClickDetails -001", lErr.Error.Error())
		// Convert to JSON
		lClickdata, lErr := json.Marshal(pClickInfo)
		if lErr != nil {
			pLogDetails.Log("InsertClickDetails -002", lErr.Error())
			pLogDetails.Log("Error Data:", pClickInfo)
		}

		// Push to Redis stream
		lerr := Gredis.XAdd(Gctx, &redis.XAddArgs{
			Stream: "ad-clicks",
			Values: map[string]interface{}{
				"data": lClickdata,
			},
		}).Err()
		if lerr != nil {
			pLogDetails.Log("InsertClickDetails -003", lerr.Error())
		}
		metrics.RequestErrors.WithLabelValues("/ads/click").Inc()
	}
}
