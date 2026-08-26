package routes

import (
	"fmt"
	"net/http"

	"opendonasi-backend/handlers"
	"opendonasi-backend/middleware"
	"opendonasi-backend/seed"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, jwtSecret []byte) {

	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	campaignHandler := handlers.NewCampaignHandler(db)
	donationHandler := handlers.NewDonationHandler(db, jwtSecret)
	fundraiserHandler := handlers.NewFundraiserHandler(db)
	reportHandler := handlers.NewReportHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "OpenDonasi API is running!"})
		})

		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		api.GET("/seed-users", func(c *gin.Context) {
			usersCreated, campaignsCreated, err := seed.SeedUsersAndCampaigns(db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d users and %d campaigns seeded successfully", usersCreated, campaignsCreated)})
		})

		
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			
			protected.GET("/profile", userHandler.GetProfile)

			
			protected.GET("/donations/me", donationHandler.GetMyDonations)

			
			protected.GET("/admin/donations", donationHandler.GetAdminDonations)
			protected.GET("/users", userHandler.GetUsers)
			protected.DELETE("/users/:id", userHandler.DeleteUser)
			protected.GET("/admin/campaigns", campaignHandler.GetAdminCampaigns)
			protected.PATCH("/admin/campaigns/:id/status", campaignHandler.UpdateCampaignStatus)
			protected.DELETE("/admin/campaigns/:id", campaignHandler.DeleteCampaign)
			protected.GET("/transactions", donationHandler.GetTransactions)
			protected.PATCH("/admin/transactions/:id/cairkan", donationHandler.CairkanTransaction)

			
			protected.GET("/fundraiser/donations", fundraiserHandler.GetFundraiserDonations)
			protected.GET("/fundraiser/funds/summary", fundraiserHandler.GetFundraiserSummary)
			protected.GET("/fundraiser/campaigns", fundraiserHandler.GetFundraiserCampaigns)
			protected.POST("/fundraiser/campaigns", fundraiserHandler.CreateFundraiserCampaign)
			protected.POST("/campaigns/:id/photo", campaignHandler.UploadCampaignPhoto)
			protected.POST("/fundraiser/reports", reportHandler.CreateReport)
		}

		
		api.GET("/stats", campaignHandler.GetStatistics)
		api.GET("/statistics", campaignHandler.GetStatistics)
		api.GET("/campaigns", campaignHandler.GetCampaigns)
		api.GET("/campaigns/:id", campaignHandler.GetCampaign)
		api.POST("/donations", donationHandler.CreateDonation)
		api.GET("/campaigns/:id/donations", donationHandler.GetCampaignDonations)
	}

	
	admin := r.Group("/admin")
	{
		admin.GET("/login", adminHandler.LoginPage)
		admin.POST("/login", adminHandler.LoginWeb)
		admin.GET("/logout", adminHandler.LogoutWeb)

		protectedAdmin := admin.Group("/")
		protectedAdmin.Use(middleware.AdminAuthMiddleware())
		{
			protectedAdmin.GET("/dashboard", adminHandler.DashboardPage)
			protectedAdmin.GET("/verify-campaign", adminHandler.VerifyCampaignPage)
			protectedAdmin.POST("/verify-campaign/:id/approve", adminHandler.ApproveCampaign)
			protectedAdmin.POST("/verify-campaign/:id/reject", adminHandler.RejectCampaign)

			protectedAdmin.GET("/verify-donation", adminHandler.VerifyDonationPage)
			protectedAdmin.POST("/verify-donation/:id/approve", adminHandler.ApproveDonation)
			protectedAdmin.POST("/verify-donation/:id/reject", adminHandler.RejectDonation)

			protectedAdmin.GET("/transfer-dana", adminHandler.TransferDanaPage)
			protectedAdmin.POST("/transfer-dana/:id/transfer", adminHandler.ProcessTransfer)

			protectedAdmin.GET("/verify-report", adminHandler.VerifyReportPage)
			protectedAdmin.POST("/verify-report/:id/approve", adminHandler.ApproveReport)
			protectedAdmin.POST("/verify-report/:id/reject", adminHandler.RejectReport)

			protectedAdmin.GET("/manage-users", adminHandler.ManageUsersPage)
			protectedAdmin.POST("/manage-users/:id/delete", adminHandler.DeleteUserWeb)
		}
	}
}
