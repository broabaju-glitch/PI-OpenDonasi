package handlers

import (
	"log"
	"net/http"

	"opendonasi-backend/middleware"
	"opendonasi-backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	adminEmail    = "admin@opendonasi.com"
	adminPassword = "admin123"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) LoginPage(c *gin.Context) {
	cookie, err := c.Cookie(middleware.CookieName)
	if err == nil && cookie == middleware.CookieValue {
		c.Redirect(http.StatusFound, "/admin/dashboard")
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Error": "",
	})
}

func (h *AdminHandler) LoginWeb(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	if email != adminEmail || password != adminPassword {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"Error": "Email atau password salah. Gunakan: admin@opendonasi.com / admin123",
		})
		return
	}

	c.SetCookie(
		middleware.CookieName,
		middleware.CookieValue,
		3600*24,
		"/admin",
		"",
		false,
		true,
	)

	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (h *AdminHandler) LogoutWeb(c *gin.Context) {
	c.SetCookie(middleware.CookieName, "", -1, "/admin", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}

func (h *AdminHandler) DashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{})
}

func (h *AdminHandler) VerifyCampaignPage(c *gin.Context) {
	var campaigns []models.Campaign
	h.DB.Where("status = ?", models.StatusMenungguVerifikasi).Order("created_at desc").Find(&campaigns)
	c.HTML(http.StatusOK, "verify_campaign.html", gin.H{
		"Campaigns": campaigns,
	})
}

func (h *AdminHandler) ApproveCampaign(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Campaign{}).Where("id = ?", id).Update("status", models.StatusAktif)
	log.Printf("✅ Campaign #%s approved (status → Aktif)", id)
	c.Redirect(http.StatusFound, "/admin/verify-campaign")
}

func (h *AdminHandler) RejectCampaign(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Campaign{}).Where("id = ?", id).Update("status", "Ditolak")
	log.Printf("❌ Campaign #%s rejected (status → Ditolak)", id)
	c.Redirect(http.StatusFound, "/admin/verify-campaign")
}

func (h *AdminHandler) VerifyDonationPage(c *gin.Context) {
	var donations []models.Donation
	h.DB.Preload("Campaign").Where("status = ?", models.DonationMenungguVerifikasi).Order("created_at desc").Find(&donations)
	c.HTML(http.StatusOK, "verify_donation.html", gin.H{
		"Donations": donations,
	})
}

func (h *AdminHandler) ApproveDonation(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Donation{}).Where("id = ?", id).Update("status", models.DonationBerhasil)

	var donation models.Donation
	if h.DB.First(&donation, id).Error == nil {
		h.DB.Model(&models.Campaign{}).Where("id = ?", donation.CampaignID).
			Update("dana_terkumpul", gorm.Expr("dana_terkumpul + ?", donation.Amount))
	}

	log.Printf("✅ Donation #%s approved (status → Berhasil)", id)
	c.Redirect(http.StatusFound, "/admin/verify-donation")
}

func (h *AdminHandler) RejectDonation(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Donation{}).Where("id = ?", id).Update("status", models.DonationDitolak)
	log.Printf("❌ Donation #%s rejected (status → Ditolak)", id)
	c.Redirect(http.StatusFound, "/admin/verify-donation")
}

func (h *AdminHandler) TransferDanaPage(c *gin.Context) {
	var campaigns []models.Campaign
	h.DB.Preload("Penggalang").
		Where("status IN ?", []string{
			string(models.StatusTargetTercapai),
			string(models.StatusBerakhir),
		}).
		Order("created_at desc").Find(&campaigns)
	c.HTML(http.StatusOK, "transfer_dana.html", gin.H{
		"Campaigns": campaigns,
	})
}

func (h *AdminHandler) ProcessTransfer(c *gin.Context) {
	id := c.Param("id")
	var campaign models.Campaign
	if h.DB.First(&campaign, id).Error == nil {
		h.DB.Model(&campaign).Updates(map[string]interface{}{
			"status":          models.StatusMenungguLaporan,
			"dana_disalurkan": campaign.DanaTerkumpul,
		})
	}
	log.Printf("💸 Campaign #%s: Dana disalurkan (Rp%.0f)", id, campaign.DanaTerkumpul)
	c.Redirect(http.StatusFound, "/admin/transfer-dana")
}

func (h *AdminHandler) VerifyReportPage(c *gin.Context) {
	var reports []models.Report
	h.DB.Preload("Campaign").Where("status = ?", models.ReportMenungguVerifikasi).Order("created_at desc").Find(&reports)
	c.HTML(http.StatusOK, "verify_report.html", gin.H{
		"Reports": reports,
	})
}

func (h *AdminHandler) ApproveReport(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Report{}).Where("id = ?", id).Update("status", models.ReportDisetujui)

	var report models.Report
	if h.DB.First(&report, id).Error == nil {
		h.DB.Model(&models.Campaign{}).Where("id = ?", report.CampaignID).Update("status", models.StatusSelesai)
	}

	log.Printf("✅ Report #%s approved → Campaign selesai", id)
	c.Redirect(http.StatusFound, "/admin/verify-report")
}

func (h *AdminHandler) RejectReport(c *gin.Context) {
	id := c.Param("id")
	h.DB.Model(&models.Report{}).Where("id = ?", id).Update("status", "Ditolak")
	log.Printf("❌ Report #%s rejected", id)
	c.Redirect(http.StatusFound, "/admin/verify-report")
}

func (h *AdminHandler) ManageUsersPage(c *gin.Context) {
	var users []models.User
	h.DB.Order("created_at desc").Find(&users)

	var totalPenggalang, totalDonatur int64
	h.DB.Model(&models.User{}).Where("role = ?", models.RolePenggalang).Count(&totalPenggalang)
	h.DB.Model(&models.User{}).Where("role = ?", models.RoleDonatur).Count(&totalDonatur)

	c.HTML(http.StatusOK, "manage_users.html", gin.H{
		"Users":           users,
		"TotalUsers":      len(users),
		"TotalPenggalang": totalPenggalang,
		"TotalDonatur":    totalDonatur,
	})
}

func (h *AdminHandler) DeleteUserWeb(c *gin.Context) {
	id := c.Param("id")
	h.DB.Delete(&models.User{}, id)
	log.Printf("🗑️ User #%s deleted", id)
	c.Redirect(http.StatusFound, "/admin/manage-users")
}
