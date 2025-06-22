package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"videoanalytics/apps/metrics"
	"videoanalytics/apps/models"
	"videoanalytics/apps/utils"
	"videoanalytics/database"

	"github.com/gin-gonic/gin"
)

func GetAdAnalyticDetails(r *gin.Context) {
	log := new(utils.LoggerId)
	log.SetSid(r.Request)
	log.Log("GetAdAnalyticDetails (+)")
	defer log.Log("GetAdAnalyticDetails (-)")

	metrics.HttpRequests.WithLabelValues("/ads/analytics", "GET").Inc()

	var lResponse any

	lFromStr := r.Request.URL.Query().Get("from")
	lToStr := r.Request.URL.Query().Get("to")

	lMinutesStr := r.Request.URL.Query().Get("minutes")

	if (lFromStr == "" || lToStr == "") && (lMinutesStr == "") {
		log.Log("GetAdAnalyticDetails-001 (Time Frame is required)")
		metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
		r.JSON(http.StatusBadRequest, gin.H{
			"status": "E",
			"errMsg": "Time Frame is required",
		})
		return
	}
	if lMinutesStr != "" {
		lMinutes, lErr := strconv.Atoi(lMinutesStr)
		if lErr != nil {
			metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
			log.Log("GetAdAnalyticDetails-002")
			r.JSON(http.StatusBadRequest, gin.H{
				"status": "E",
				"errMsg": "Time Frame is invalid",
			})
			return
		}
		lResponse, lErr = GetMinutesAnalytic(log, lMinutes)
		if lErr != nil {
			metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
			log.Log("GetAdAnalyticDetails-003", lErr.Error())
			r.JSON(http.StatusInternalServerError, gin.H{
				"status": "E",
				"errMsg": "Data Fetch Error",
			})
			return
		}
	} else {
		lFrom, lErr := time.Parse(time.RFC3339, lFromStr)
		if lErr != nil {
			metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
			log.Log("GetAdAnalyticDetails-004")
			r.JSON(http.StatusBadRequest, gin.H{
				"status": "E",
				"errMsg": "Start Time Frame is invalid",
			})
			return
		}
		lTo, lErr := time.Parse(time.RFC3339, lToStr)
		if lErr != nil {
			metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
			log.Log("GetAdAnalyticDetails-004", lErr.Error())
			r.JSON(http.StatusBadRequest, gin.H{
				"status": "E",
				"errMsg": "End Time Frame is invalid",
			})
			return
		}
		lResponse, lErr = GetAnalyticDetails(log, lFrom, lTo)
		log.Log("lResponse", lResponse)
		if lErr != nil {
			metrics.RequestErrors.WithLabelValues("/ads/analytics").Inc()
			log.Log("GetAdAnalyticDetails-005", lErr.Error())
			r.JSON(http.StatusInternalServerError, gin.H{
				"status": "E",
				"errMsg": "Data Fetch Error",
			})
			return
		}
	}
	r.Writer.Header().Set("Content-Type", "application/json")
	r.JSON(http.StatusOK, gin.H{
		"status":       "S",
		"errMsg":       "",
		"adsAnalytics": lResponse,
	})
}

func GetAnalyticDetails(pLogDetails *utils.LoggerId, pFrom time.Time, pTo time.Time) ([]models.ResponseDetails, error) {
	pLogDetails.Log("GetAnalyticDetails (+)")
	defer pLogDetails.Log("GetAnalyticDetails (-)")

	var lAdAnalytics []models.ResponseDetails
	log.Println("from", pFrom, pTo)
	lErr := database.DB.Table("ad_details a").Select("a.id AS ad_id,COUNT(ac.id) AS clicks,view_count as impressions").Joins("LEFT JOIN ad_click_rates ac on ac.ad_id = a.id and ac.clicked_at BETWEEN ? AND ?", pFrom, pTo).Group("a.id").Find(&lAdAnalytics)
	if lErr.Error != nil {
		pLogDetails.Log("GetAnalyticDetails -001", lErr.Error.Error())
		return lAdAnalytics, fmt.Errorf("GetAnalyticDetails-001 %s", lErr.Error.Error())
	} else {
		for index, lDetails := range lAdAnalytics {
			// impressions := int64(100)
			if lDetails.Impr == 0 {
				lAdAnalytics[index].CTR = 0 // or nil / omit field / -1 to indicate “undefined”
				continue
			}
			ctr := float64(lDetails.Clicks) / float64(lDetails.Impr)
			lAdAnalytics[index].CTR = ctr
			lAdAnalytics[index].Impr = lDetails.Impr
			log.Println("lAdAnalytics", lDetails.CTR, lDetails.Clicks, lDetails.Impr)
		}
	}
	log.Println("lAdAnalytics", lAdAnalytics)

	return lAdAnalytics, nil
}

func GetMinutesAnalytic(pLogDetails *utils.LoggerId, pMinutes int) ([]models.AdMinutesAnalysics, error) {
	pLogDetails.Log("GetMinutesAnalytic (+)")
	defer pLogDetails.Log("GetMinutesAnalytic (-)")

	var lAdAnalytics []models.AdMinutesAnalysics

	lErr := database.DB.Table("ad_click_rates").Select("ad_id as ac_id,COUNT(*) AS total_clicks,SUM(clicked_at >= NOW() - INTERVAL ? MINUTE) AS recent_clicks", pMinutes).Group("ad_id").Find(&lAdAnalytics)
	if lErr.Error != nil {
		pLogDetails.Log("GetMinutesAnalytic -001", lErr.Error.Error())
		return lAdAnalytics, fmt.Errorf("GetMinutesAnalytic-001 %s", lErr.Error.Error())
	}
	return lAdAnalytics, nil
}
