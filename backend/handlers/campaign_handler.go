package handlers

import (
	"fmt"
	"io"
	"log"
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

type CampaignHandler struct {
	DB *gorm.DB
}

func NewCampaignHandler(db *gorm.DB) *CampaignHandler {
	return &CampaignHandler{DB: db}
}

func (h *CampaignHandler) GetAdminCampaigns(c *gin.Context) {
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

	var campaigns []models.Campaign
	if err := h.DB.Preload("Penggalang").Order("created_at desc").Find(&campaigns).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data campaign"})
		return
	}

	c.JSON(http.StatusOK, campaigns)
}

func (h *CampaignHandler) UpdateCampaignStatus(c *gin.Context) {
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

	campaignID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status wajib diisi"})
		return
	}

	statusValue := models.CampaignStatus(input.Status)
	if statusValue != models.StatusAktif && statusValue != models.StatusDitolak && statusValue != models.StatusSelesai && input.Status != "disetujui" && input.Status != "selesai" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid"})
		return
	}


	if input.Status == "disetujui" {
		statusValue = models.StatusAktif
	}
	if input.Status == "selesai" {
		statusValue = models.StatusSelesai
	}

	var campaign models.Campaign
	if err := h.DB.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}

	campaign.Status = statusValue
	if err := h.DB.Save(&campaign).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status berhasil diperbarui"})
}

func (h *CampaignHandler) DeleteCampaign(c *gin.Context) {
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

	campaignID := c.Param("id")

	if err := h.DB.Delete(&models.Campaign{}, campaignID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus campaign"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Campaign berhasil dihapus"})
}

func (h *CampaignHandler) GetStatistics(c *gin.Context) {
	var activeCampaigns, completedCampaigns int64
	var totalDonatur int64

	now := time.Now()

	h.DB.Model(&models.Campaign{}).Where("dana_terkumpul < target_dana AND end_date >= ?", now).Count(&activeCampaigns)


	h.DB.Model(&models.Campaign{}).Where("dana_terkumpul >= target_dana OR end_date < ?", now).Count(&completedCampaigns)


	h.DB.Model(&models.User{}).Where("role = ?", models.RoleDonatur).Count(&totalDonatur)


	var result struct{ Total float64 }
	h.DB.Model(&models.Campaign{}).Select("COALESCE(SUM(dana_terkumpul),0) as total").Scan(&result)

	c.JSON(http.StatusOK, gin.H{
		"active_campaigns":    activeCampaigns,
		"total_funds":         result.Total,
		"total_collected":     result.Total, 
		"total_donors":        totalDonatur,
		"total_donatur":       totalDonatur, 
		"completed_campaigns": completedCampaigns,
	})
}

func (h *CampaignHandler) GetCampaigns(c *gin.Context) {
	var campaigns []models.Campaign
	h.DB.Preload("Penggalang").Find(&campaigns)
	c.JSON(http.StatusOK, campaigns)
}

func (h *CampaignHandler) GetCampaign(c *gin.Context) {
	var campaign models.Campaign
	if err := h.DB.Preload("Penggalang").First(&campaign, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, campaign)
}

func (h *CampaignHandler) UploadCampaignPhoto(c *gin.Context) {
	userID, _ := c.Get("userID")
	campaignID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID campaign tidak valid"})
		return
	}


	var campaign models.Campaign
	if err := h.DB.First(&campaign, campaignID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Campaign tidak ditemukan"})
		return
	}
	if campaign.PenggalangID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Anda tidak memiliki akses ke campaign ini"})
		return
	}


	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File foto tidak ditemukan dalam request"})
		return
	}


	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung. Gunakan JPG, PNG, atau WebP."})
		return
	}


	uploadDir := "public/uploads/campaigns"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori upload"})
		return
	}


	fileName := fmt.Sprintf("campaign_%d_%d%s", campaignID, time.Now().UnixNano(), ext)
	dstPath := filepath.Join(uploadDir, fileName)


	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuka file"})
		return
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file ke disk"})
		return
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menulis file"})
		return
	}


	fotoURL := fmt.Sprintf("/uploads/campaigns/%s", fileName)
	if err := h.DB.Model(&campaign).Update("foto", fotoURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui database"})
		return
	}

	log.Printf("🖼️ Foto campaign #%d diperbarui: %s", campaignID, fotoURL)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Foto campaign berhasil diperbarui.",
		"foto_url": fotoURL,
	})
}
