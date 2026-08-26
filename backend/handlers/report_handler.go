package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"opendonasi-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	DB *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{DB: db}
}

func (h *ReportHandler) CreateReport(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	// Verifikasi role
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil || user.Role != models.RolePenggalang {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya penggalang yang dapat mengunggah laporan"})
		return
	}

	campaignIDStr := c.PostForm("campaign_id")
	title := c.PostForm("title")
	description := c.PostForm("description")

	if campaignIDStr == "" || title == "" || description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Campaign, judul, dan deskripsi wajib diisi"})
		return
	}

	campaignID, err := strconv.ParseUint(campaignIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID Campaign tidak valid"})
		return
	}

	// Verifikasi campaign milik penggalang tersebut
	var campaign models.Campaign
	if err := h.DB.First(&campaign, campaignID).Error; err != nil || campaign.PenggalangID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak berhak mengakses campaign ini"})
		return
	}

	// Handle image upload
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Foto laporan wajib dilampirkan"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Gunakan JPG, PNG, atau WebP."})
		return
	}

	uploadDir := "public/uploads/reports"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori upload"})
		return
	}

	fileName := fmt.Sprintf("report_%d_%d%s", campaignID, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(fileHeader, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file gambar"})
		return
	}

	fotoURL := fmt.Sprintf("/uploads/reports/%s", fileName)

	report := models.Report{
		CampaignID:  uint(campaignID),
		Title:       title,
		Description: description,
		Image:       fotoURL,
		Status:      models.ReportMenungguVerifikasi,
	}

	if err := h.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan laporan ke database"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Laporan berhasil diunggah",
		"report":  report,
	})
}
