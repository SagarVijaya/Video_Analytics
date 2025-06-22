package api

import (
	"fmt"
	"net/http"
	"videoanalytics/apps/metrics"
	"videoanalytics/apps/models"
	"videoanalytics/apps/utils"
	"videoanalytics/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAdsList(r *gin.Context) {
	log := new(utils.LoggerId)
	log.SetSid(r.Request)
	log.Log("GetAdList (+)")
	defer log.Log("GetAdList (-)")

	metrics.HttpRequests.WithLabelValues("/ads", "GET").Inc()

	lResponse, lErr := FetchAdDetails(log)
	if lErr != nil {
		metrics.RequestErrors.WithLabelValues("/ads").Inc()
		log.Log("GetAdList -001", lErr.Error())
		r.JSON(http.StatusInternalServerError, gin.H{
			"status": "E",
			"errMsg": "Data Fetch Error",
		})
		return
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.JSON(http.StatusOK, gin.H{
		"status":  "S",
		"errMsg":  "",
		"adsList": lResponse,
	})
}

func FetchAdDetails(pLogDetails *utils.LoggerId) ([]models.AdDetails, error) {
	pLogDetails.Log("FetchAdDetails (+)")
	defer pLogDetails.Log("FetchAdDetails (-)")
	var lAdDetails []models.AdDetails
	lErr := database.DB.Table("ad_details ad").Select("ad.ID as ID,ad.image_url as image_url,ad.target_url as target_url,created_at").Find(&lAdDetails)
	if lErr.Error != nil {
		pLogDetails.Log("FetchAdDetails - 001", lErr.Error.Error())
		return lAdDetails, fmt.Errorf("FetchAdDetails - 001 %s", lErr.Error.Error())
	} else {
		go UpdateCountValue(pLogDetails, lAdDetails)
	}

	return lAdDetails, nil
}

func UpdateCountValue(pLogDetails *utils.LoggerId, pDetailsCount []models.AdDetails) {
	pLogDetails.Log("UpdateCountValue (+)")
	defer pLogDetails.Log("UpdateCountValue (-)")
	for _, lValue := range pDetailsCount {
		go func() {

			lErr := database.DB.Table("ad_details ad").Where("id =?", lValue.ID).Updates(map[string]interface{}{
				"view_count": gorm.Expr("`view_count` + ?", 1),
			})
			if lErr.Error != nil {
				pLogDetails.Log("UpdateCountValue - 001", lErr.Error.Error())
			}
		}()
	}
}
