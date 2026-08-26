package handlers

import (
	"log"
	"net/http"
	"time"

	"opendonasi-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FundraiserHandler struct {
	DB *gorm.DB
}

func NewFundraiserHandler(db *gorm.DB) *FundraiserHandler {
	return &FundraiserHandler{DB: db}
}

func (h *FundraiserHandler) GetFundraiserDonations(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RolePenggalang {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya penggalang yang dapat mengakses"})
		return
	}

	var donations []models.Donation
	
	h.DB.Preload("Campaign").Preload("Donatur").
		Joins("JOIN campaigns ON campaigns.id = donations.campaign_id").
		Where("campaigns.penggalang_id = ?", userID).
		Order("donations.created_at desc").Find(&donations)

	c.JSON(http.StatusOK, donations)
}

func (h *FundraiserHandler) GetFundraiserSummary(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RolePenggalang {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya penggalang yang dapat mengakses"})
		return
	}

	
	var campaignIDs []uint
	h.DB.Model(&models.Campaign{}).Where("penggalang_id = ?", userID).Pluck("id", &campaignIDs)

	var totalTerkumpul float64
	var totalWithdrawn float64
	var saldoAktif float64

	if len(campaignIDs) > 0 {
	
		h.DB.Model(&models.Donation{}).Where("campaign_id IN ? AND status IN ?", campaignIDs, []models.DonationStatus{models.DonationBerhasil, models.DonationDicairkan}).Select("COALESCE(SUM(amount), 0)").Scan(&totalTerkumpul)

	
		h.DB.Model(&models.Donation{}).Where("campaign_id IN ? AND status = ?", campaignIDs, models.DonationDicairkan).Select("COALESCE(SUM(amount), 0)").Scan(&totalWithdrawn)

		
		h.DB.Model(&models.Donation{}).Where("campaign_id IN ? AND status = ?", campaignIDs, models.DonationBerhasil).Select("COALESCE(SUM(amount), 0)").Scan(&saldoAktif)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_terkumpul": totalTerkumpul,
		"total_pending":   0, // No more pending withdrawals conceptually
		"total_withdrawn": totalWithdrawn,
		"saldo_aktif":     saldoAktif,
		"has_pending":     false,
	})
}

func (h *FundraiserHandler) GetFundraiserCampaigns(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
		return
	}
	userID, ok := userIDRaw.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID tidak dapat dibaca"})
		return
	}


	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RolePenggalang {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya penggalang yang dapat mengakses"})
		return
	}

	var campaigns []models.Campaign
	result := h.DB.Where("penggalang_id = ?", userID).Order("created_at desc").Find(&campaigns)
	if result.Error != nil {
		log.Printf("❌ DB error fetching campaigns for user %d: %v", userID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data campaign"})
		return
	}

	log.Printf("📋 Fetched %d campaigns for penggalang #%d", len(campaigns), userID)

	if campaigns == nil {
		campaigns = []models.Campaign{}
	}
	c.JSON(http.StatusOK, campaigns)
}

func (h *FundraiserHandler) CreateFundraiserCampaign(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RolePenggalang {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya penggalang yang dapat membuat campaign"})
		return
	}

	var input struct {
		Title         string  `json:"title" binding:"required"`
		Category      string  `json:"category" binding:"required"`
		Description   string  `json:"description" binding:"required"`
		TargetDana    float64 `json:"target_dana" binding:"required"`
		AlamatLengkap string  `json:"alamat_lengkap" binding:"required"`
		LinkGmaps     string  `json:"link_gmaps"`
		StartDate     string  `json:"start_date"`
		EndDate       string  `json:"end_date"`
		Rekening      string  `json:"rekening"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data tidak lengkap: " + err.Error()})
		return
	}

	campaign := models.Campaign{
		Title:         input.Title,
		Category:      input.Category,
		Description:   input.Description,
		TargetDana:    input.TargetDana,
		AlamatLengkap: input.AlamatLengkap,
		LinkGmaps:     input.LinkGmaps,
		Rekening:      input.Rekening,
		Status:        models.StatusMenungguVerifikasi,
		PenggalangID:  userID.(uint),
	}

	if input.StartDate != "" {
		if t, err := time.Parse("2006-01-02", input.StartDate); err == nil {
			campaign.StartDate = t
		}
	}
	if input.EndDate != "" {
		if t, err := time.Parse("2006-01-02", input.EndDate); err == nil {
			campaign.EndDate = t
		}
	}

	if err := h.DB.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan campaign"})
		return
	}

	log.Printf("📢 Campaign baru dibuat: '%s' oleh user #%v", campaign.Title, userID)
	c.JSON(http.StatusCreated, gin.H{
		"message":  "Campaign berhasil dibuat! Menunggu verifikasi admin.",
		"campaign": campaign,
	})
}
