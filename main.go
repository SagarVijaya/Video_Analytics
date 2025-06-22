package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"videoanalytics/apps/api"
	"videoanalytics/apps/metrics"
	"videoanalytics/config"
	"videoanalytics/database"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func init() {
	config.LoadGlobalConfig(".env")
	api.Gredis = redis.NewClient(&redis.Options{
		Addr: config.GetConfig().Redis.Port,
	})
	log.Println("Gredis", api.Gredis)
}

func main() {

	// Create the log directory if it doesn't exist
	err := os.MkdirAll("log", os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	logFolderName := "./log/log" + time.Now().Format("02012006.15.04.05") + ".log"
	file, err := os.OpenFile(logFolderName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	multiOut := io.MultiWriter(file, os.Stdout)
	// make the default log package use it, with timestamps
	log.SetOutput(multiOut)

	err = database.ConnectDatabase()
	if err != nil {
		log.Fatal("Database connection error:", err)
	}

	metrics.StartMetricsServer()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	go router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "This page is not available,Please check the url",
		})
	})

	router.GET("/ads", api.GetAdsList)
	router.GET("/BackUpData", api.BackUpDataUpload)
	routerGroup := router.Group("/ads")
	routerGroup.POST("/click", api.StoreClickInfo)
	routerGroup.GET("/analytics", api.GetAdAnalyticDetails)
	port := config.GetConfig().Server.Port
	log.Println("Server running on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
