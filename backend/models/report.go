package models

import "gorm.io/gorm"

type ReportStatus string

const (
	ReportMenungguVerifikasi ReportStatus = "Menunggu Verifikasi"
	ReportDisetujui          ReportStatus = "Disetujui"
)

type Report struct {
	gorm.Model
	Title       string       `json:"title"`
	Description string       `json:"description" gorm:"type:text"`
	Image       string       `json:"image" gorm:"type:varchar(255)"`
	Status      ReportStatus `json:"status" gorm:"default:'Menunggu Verifikasi'"`

	CampaignID uint     `json:"campaign_id"`
	Campaign   Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`
}
