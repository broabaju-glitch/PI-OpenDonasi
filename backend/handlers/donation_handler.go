package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"opendonasi-backend/models"
	"opendonasi-backend/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type DonationHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewDonationHandler(db *gorm.DB, jwtSecret []byte) *DonationHandler {
	return &DonationHandler{DB: db, JWTSecret: jwtSecret}
}

func (h *DonationHandler) CreateDonation(c *gin.Context) {
	var input struct {
		CampaignID     uint    `form:"campaign_id" binding:"required"`
		Amount         float64 `form:"amount" binding:"required"`
		DonorName      string  `form:"donor_name"`
		WhatsappNumber string  `form:"whatsapp_number" binding:"required"`
		IsAnonymous    bool    `form:"is_anonymous"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data donasi tidak lengkap: " + err.Error()})
		return
	}

	var buktiPath string
	file, err := c.FormFile("bukti_transfer")
	if err == nil && file != nil {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" {
			uploadDir := "public/uploads/donations"
			os.MkdirAll(uploadDir, 0755)
			fileName := fmt.Sprintf("donasi_%d%s", time.Now().UnixNano(), ext)
			dst := filepath.Join(uploadDir, fileName)
			if err := c.SaveUploadedFile(file, dst); err == nil {
				buktiPath = "/uploads/donations/" + fileName
			}
		}
	}

	var campaign models.Campaign
	if err := h.DB.First(&campaign, input.CampaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}

	err = h.DB.Transaction(func(tx *gorm.DB) error {
		var donaturID *uint
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return h.JWTSecret, nil
			})
			if err == nil && token.Valid {
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					id := uint(claims["user_id"].(float64))
					donaturID = &id
				}
			}
		}

		donation := models.Donation{
			Amount:         input.Amount,
			CampaignID:     input.CampaignID,
			DonorName:      input.DonorName,
			WhatsappNumber: input.WhatsappNumber,
			IsAnonymous:    input.IsAnonymous,
			Status:         models.DonationBerhasil,
			DonaturID:      donaturID,
			BuktiTransfer:  buktiPath,
		}

		if donaturID != nil {
			donation.DonaturID = donaturID
			donation.IsAnonymous = false
		}

		if err := tx.Create(&donation).Error; err != nil {
			return err
		}

		campaign.DanaTerkumpul += input.Amount

		if campaign.DanaTerkumpul >= campaign.TargetDana && campaign.Status == models.StatusAktif {
			campaign.Status = models.StatusTargetTercapai
		}

		if err := tx.Save(&campaign).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses donasi"})
		return
	}

	whatsappMsg := services.BuildDonationMessage(campaign.Title, input.Amount)
	go services.SendFonnteWhatsApp(input.WhatsappNumber, whatsappMsg)

	log.Printf("✅ Donation successful: Rp%.0f for campaign '%s'", input.Amount, campaign.Title)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Donasi berhasil diproses! Bar progres otomatis naik.",
	})
}

func (h *DonationHandler) GetCampaignDonations(c *gin.Context) {
	var donations []models.Donation
	h.DB.Where("campaign_id = ?", c.Param("id")).Order("created_at desc").Find(&donations)

	type PublicDonation struct {
		ID          uint                  `json:"id"`
		DonorName   string                `json:"donor_name"`
		Amount      float64               `json:"amount"`
		IsAnonymous bool                  `json:"is_anonymous"`
		Status      models.DonationStatus `json:"status"`
		CreatedAt   string                `json:"created_at"`
	}

	var publicDonations []PublicDonation
	for _, d := range donations {
		name := d.DonorName
		if d.IsAnonymous {
			name = "Hamba Allah"
		}
		publicDonations = append(publicDonations, PublicDonation{
			ID:          d.ID,
			DonorName:   name,
			Amount:      d.Amount,
			IsAnonymous: d.IsAnonymous,
			Status:      d.Status,
			CreatedAt:   d.CreatedAt.Format("2006-01-02"),
		})
	}

	c.JSON(http.StatusOK, publicDonations)
}

func (h *DonationHandler) GetMyDonations(c *gin.Context) {
	userID, _ := c.Get("userID")
	var donations []models.Donation
	h.DB.Preload("Campaign").Preload("Campaign.Penggalang").Where("donatur_id = ?", userID).Order("created_at desc").Find(&donations)
	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetAdminDonations(c *gin.Context) {
	userID, _ := c.Get("userID")
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengakses"})
		return
	}

	var donations []models.Donation
	h.DB.Preload("Campaign").Preload("Donatur").Order("created_at desc").Find(&donations)
	c.JSON(http.StatusOK, donations)
}

func (h *DonationHandler) GetTransactions(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	var currentUser models.User
	if err := h.DB.First(&currentUser, userID).Error; err != nil || currentUser.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengakses"})
		return
	}

	var transactions []models.Donation
	if err := h.DB.Preload("Campaign").Preload("Donatur").Order("created_at desc").Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data transaksi"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *DonationHandler) CairkanTransaction(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	var currentUser models.User
	if err := h.DB.First(&currentUser, userID).Error; err != nil || currentUser.Role != models.RoleAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengakses"})
		return
	}

	transactionID := c.Param("id")

	var transaction models.Donation
	if err := h.DB.First(&transaction, transactionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaksi tidak ditemukan"})
		return
	}

	if transaction.Status != models.DonationBerhasil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaksi belum berhasil atau sudah dicairkan"})
		return
	}

	transaction.Status = models.DonationDicairkan
	if err := h.DB.Save(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Dana berhasil dicairkan"})
}
