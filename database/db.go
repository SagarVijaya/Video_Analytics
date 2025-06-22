package database

import (
	"database/sql"
	"fmt"
	"log"
	"videoanalytics/apps/models"
	"videoanalytics/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect Database initializes the DB and runs migrations
func ConnectDatabase() error {
	cf := config.GetConfig()
	fmt.Println(cf)

	lCreateDatabase := fmt.Sprintf("%s:%s@tcp(%s:%d)/?parseTime=true",
		cf.Database.User,
		cf.Database.Pass,
		cf.Database.Host,
		cf.Database.Port,
	)
	lCreateDbSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", cf.Database.Name)
	lErr := CreateDatabase(lCreateDatabase, lCreateDbSql)
	if lErr != nil {
		return lErr
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cf.Database.User,
		cf.Database.Pass,
		cf.Database.Host,
		cf.Database.Port,
		cf.Database.Name,
	)

	// Open GORM DB
	db, lErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if lErr != nil {
		log.Fatalf("Failed to connect to GORM database: %v", lErr)
		return lErr
	}

	// Remove the existing table
	// lErr = db.Migrator().DropTable(&models.AdDetails{}, &models.AdClickRate{})
	// if lErr != nil {
	// 	log.Fatalf("Failed to remove tables for auto migrate tables: %v", lErr)
	// 	return lErr
	// }
	// Run migrations
	lErr = db.AutoMigrate(&models.AdDetails{}, &models.AdClickRate{})
	if lErr != nil {
		log.Fatalf("Failed to auto migrate tables: %v", lErr)
		return lErr
	}

	DB = db

	lErr = AlterQueryExecution()
	if lErr != nil {
		log.Fatalf("Failed to create index: %v", lErr)
		return lErr
	}

	lErr = InsertDefaultValues()
	if lErr != nil {
		log.Fatalf("Failed to insert value to table: %v", lErr)
		return lErr
	}
	log.Println("GORM database setup completed successfully")
	return nil
}

// create database if not exists
func CreateDatabase(dsn string, createSQL string) error {
	sqlDB, lErr := sql.Open("mysql", dsn)
	if lErr != nil {
		log.Fatalf("Failed to connect to MariaDB server for DB creation: %v", lErr)
		return lErr
	}
	defer sqlDB.Close()

	if lErr = sqlDB.Ping(); lErr != nil {
		log.Fatalf("Ping failed: %v", lErr)
		return lErr
	}

	if _, lErr = sqlDB.Exec(createSQL); lErr != nil {
		log.Fatalf("Failed to create DB: %v", lErr)
		return lErr
	}

	return nil
}
func AlterQueryExecution() error {
	lErr := DB.Exec(`ALTER TABLE ad_click_rates ADD UNIQUE INDEX uniq_click (ad_id, ip, clicked_at);`).Error
	if lErr != nil {
		return lErr
	}
	lErr1 := DB.Exec(`ALTER TABLE ad_click_rates ADD INDEX idx_ad_time (ad_id, clicked_at);`).Error
	if lErr1 != nil {
		return lErr1
	}
	return nil
}

func InsertDefaultValues() error {
	lErr := DB.Exec(`INSERT INTO ad_details (image_url, target_url, created_at, view_count) VALUES
		('https://cdn.example.com/ads/ad1.jpg', 'https://example.com/product1', NOW(), 0),
		('https://cdn.example.com/ads/ad2.jpg', 'https://example.com/product2', NOW(), 0),
		('https://cdn.example.com/ads/ad3.jpg', 'https://example.com/product3', NOW(), 0),
		('https://cdn.example.com/ads/ad4.jpg', 'https://example.com/product4', NOW(), 0),
		('https://cdn.example.com/ads/ad5.jpg', 'https://example.com/product5', NOW(), 0);
	`).Error
	if lErr != nil {
		return lErr
	}
	return nil
}
