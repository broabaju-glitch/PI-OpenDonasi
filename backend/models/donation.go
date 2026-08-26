package models

import "gorm.io/gorm"

type DonationStatus string

const (
	DonationMenungguVerifikasi DonationStatus = "Menunggu Verifikasi"
	DonationBerhasil           DonationStatus = "Berhasil"
	DonationDitolak            DonationStatus = "Ditolak"
	DonationDicairkan          DonationStatus = "Sudah Dicairkan"
)

type Donation struct {
	gorm.Model
	Amount         float64        `json:"amount"`
	BuktiTransfer  string         `json:"bukti_transfer"`
	Status         DonationStatus `json:"status" gorm:"default:'Menunggu Verifikasi'"`
	WhatsappNumber string         `json:"whatsapp_number"`
	IsAnonymous    bool           `json:"is_anonymous" gorm:"default:false"`
	DonorName      string         `json:"donor_name"` // Name provided during donation (may differ from user name)

	CampaignID uint     `json:"campaign_id"`
	Campaign   Campaign `json:"campaign" gorm:"foreignKey:CampaignID"`

	DonaturID *uint `json:"donatur_id" gorm:"column:donatur_id"` // nullable because optional login
	Donatur   User  `json:"donatur" gorm:"foreignKey:DonaturID"`
}
