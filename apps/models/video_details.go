package models

import "time"

type AdDetails struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	ImageURL  string    `json:"image_url" gorm:"image_url"`
	TargetURL string    `json:"target_url" gorm:"target_url"`
	ViewCount int       `json:"viewCount,omitempty" gorm:"ViewCount"`
	CreatedAt time.Time `json:"created_at" gorm:"created_at"`
}

type AdClickRate struct {
	ID             int       `gorm:"primaryKey"`
	AdID           int       `json:"ad_id"`
	ClickedAt      time.Time `json:"clicked_at"`
	IP             string    `json:"ip"`
	PlaybackSecond float64   `json:"playback_second"`
}

type AdMinutesAnalysics struct {
	AdID         int   `json:"ad_id"`
	TotalClicks  int64 `json:"total_clicks"`
	RecentClicks int64 `json:"recent_clicks"`
}

type ResponseDetails struct {
	AdID   int     `json:"ad_id"`
	Clicks int64   `json:"clicks"`
	Impr   int64   `json:"impressions"`
	CTR    float64 `json:"ctr"`
}
