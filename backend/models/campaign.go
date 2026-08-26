package models

import (
	"time"

	"gorm.io/gorm"
)

type CampaignStatus string

const (
	StatusMenungguVerifikasi CampaignStatus = "Menunggu Verifikasi"
	StatusAktif              CampaignStatus = "Aktif"
	StatusDitolak            CampaignStatus = "Ditolak"
	StatusTargetTercapai     CampaignStatus = "Target Tercapai"
	StatusBerakhir           CampaignStatus = "Berakhir"
	StatusDanaDisalurkan     CampaignStatus = "Dana Disalurkan"
	StatusMenungguLaporan    CampaignStatus = "Menunggu Laporan"
	StatusSelesai            CampaignStatus = "Selesai"
)

type Campaign struct {
	gorm.Model
	Title          string         `json:"title"`
	Category       string         `json:"category"`
	Description    string         `json:"description" gorm:"type:text"`
	TargetDana     float64        `json:"target_dana"`
	DanaTerkumpul  float64        `json:"dana_terkumpul"`
	DanaDisalurkan float64        `json:"dana_disalurkan"`
	LokasiKejadian string         `json:"lokasi_kejadian"`
	AlamatLengkap  string         `json:"alamat_lengkap"`
	LinkGmaps      string         `json:"link_gmaps"`
	Foto           string         `json:"foto"`
	StartDate      time.Time      `json:"start_date"`
	EndDate        time.Time      `json:"end_date"`
	Rekening       string         `json:"rekening"` // specific target bank account for this campaign if differs
	Status         CampaignStatus `json:"status" gorm:"default:'Menunggu Verifikasi'"`

	PenggalangID uint `json:"penggalang_id"`
	Penggalang   User `json:"penggalang" gorm:"foreignKey:PenggalangID"`
}
