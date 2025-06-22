package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	mailsender "videoanalytics/apps/mailSender"
	"videoanalytics/apps/metrics"
	"videoanalytics/apps/models"
	"videoanalytics/apps/utils"
	"videoanalytics/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/clause"
)

func BackUpDataUpload(r *gin.Context) {
	log := new(utils.LoggerId)
	log.SetSid(r.Request)
	log.Log("BackUpDataUpload (+)")
	defer log.Log("BackUpDataUpload (+)")

	metrics.HttpRequests.WithLabelValues("/BackUpData", "GET").Inc()

	lErr := RedisFetchValue(log)
	if lErr != nil {
		if lErr == redis.Nil {
			r.Writer.Header().Set("Content-Type", "application/json")
			r.JSON(http.StatusOK, gin.H{
				"status": "S",
				"Msg":    "No Data Available in Redis",
			})
			return
		}
		metrics.RequestErrors.WithLabelValues("/BackUpData").Inc()
		log.Log("BackUpDataUpload-001", lErr.Error())
		r.JSON(http.StatusBadRequest, gin.H{
			"status": "E",
			"Msg":    "Data Not Inserted Properly",
		})
		return
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.JSON(http.StatusOK, gin.H{
		"status": "S",
		"Msg":    "Data Inserted Successfully",
	})
}

func RedisFetchValue(pLogDetails *utils.LoggerId) error {
	pLogDetails.Log("RedisFetchValue (+)")
	defer pLogDetails.Log("RedisFetchValue (-)")

	lastID := "0"

	for {
		lStreams, lErr := Gredis.XRead(Gctx, &redis.XReadArgs{
			Streams: []string{"ad-clicks", lastID},
			Count:   10, // max messages per fetch(0 to fetch all)
			Block:   5 * time.Second,
		}).Result()
		if lErr != nil {
			if lErr == redis.Nil {
				pLogDetails.Log("RedisFetchValue-001 (No Data Available)", lErr.Error())
				return lErr
			}
			pLogDetails.Log("RedisFetchValue-001 (Redis XREAD error:)", lErr.Error())
			return fmt.Errorf("RedisFetchValue-001 (Redis XREAD error:) %s", lErr.Error())
		}

		if len(lStreams) == 0 || len(lStreams[0].Messages) == 0 {
			pLogDetails.Log("No more data is available. Exiting loop.")
			break
		}
		for _, lMsg := range lStreams[0].Messages {
			var lClickDetails models.AdClickRate
			lastID = lMsg.ID // advance cursor so we don’t reread the same entry

			lRaw, ok := lMsg.Values["data"].(string)
			if !ok {
				pLogDetails.Log("RedisFetchValue-002 (Expected string)")
				return fmt.Errorf("RedisFetchValue-002 (Expected string)")
			}
			lErr := json.Unmarshal([]byte(lRaw), &lClickDetails)
			if lErr != nil {
				pLogDetails.Log("RedisFetchValue-003 (Unmarshal in redis data)", lErr.Error())
				return fmt.Errorf("RedisFetchValue-003 (Unmarshal in redis data) %s", lErr.Error())
				// log.Println("Bad JSON:", err)
				// continue
			} else {
				lErr := database.DB.Table("ad_click_rates").Clauses(clause.OnConflict{DoNothing: true}).Create(&lClickDetails)
				if lErr.Error != nil {
					go mailsender.SendFailureAlert(lClickDetails, lErr.Error)
					pLogDetails.Log("RedisFetchValue -004", lErr.Error.Error())
					return fmt.Errorf("RedisFetchValue-004 %s", lErr.Error.Error())
				} else {
					pLogDetails.Log("Removed Ids :", lMsg.ID)
					Gredis.XDel(Gctx, "ad-clicks", lMsg.ID)
					// metrics.AdClicks.WithLabelValues(strconv.Itoa(int(lClickDetails.AdID))).Inc()
				}
			}
		}
	}
	return nil
}
