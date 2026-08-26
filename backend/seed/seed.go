package seed

import (
	"opendonasi-backend/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsersAndCampaigns(db *gorm.DB) (int, int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return 0, 0, err
	}

	users := []models.User{
		{Name: "Admin Utama", Email: "admin@test.com", Password: string(hash), Role: models.RoleAdmin},
		{Name: "Donatur Setia", Email: "donatur@test.com", Password: string(hash), Role: models.RoleDonatur},
		{Name: "Penggalang Dana", Email: "penggalang@test.com", Password: string(hash), Role: models.RolePenggalang},
	}

	var usersCreated int
	var penggalangID uint
	for _, u := range users {
		var existing models.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err != nil {
			if db.Create(&u).Error == nil {
				usersCreated++
				if u.Role == models.RolePenggalang {
					penggalangID = u.ID
				}
			}
		} else if existing.Role == models.RolePenggalang {
			penggalangID = existing.ID
		}
	}

	var campaignsCreated int
	if penggalangID != 0 {
		campaigns := []models.Campaign{
			{
				Title:          "Bantu Korban Banjir Bandang Garut",
				Category:       "Banjir",
				Description:    "Banjir bandang melanda puluhan desa. Ribuan warga mengungsi dan membutuhkan bantuan logistik segera.",
				TargetDana:     100000000,
				DanaTerkumpul:  75000000,
				LokasiKejadian: "Garut, Jawa Barat",
				Status:         models.StatusAktif,
				PenggalangID:   penggalangID,
			},
			{
				Title:          "Darurat Gempa Cianjur — Bangun Kembali",
				Category:       "Gempa Bumi",
				Description:    "Ratusan rumah hancur akibat gempa. Bantuan donasi ini akan disalurkan untuk membangun hunian sementara warga.",
				TargetDana:     500000000,
				DanaTerkumpul:  250000000,
				LokasiKejadian: "Cianjur, Jawa Barat",
				Status:         models.StatusAktif,
				PenggalangID:   penggalangID,
			},
			{
				Title:          "Bantuan Logistik Erupsi Gunung Semeru",
				Category:       "Gunung Meletus",
				Description:    "Pengungsi membutuhkan masker, obat-obatan, dan logistik bahan pokok selama masa tanggap darurat.",
				TargetDana:     150000000,
				DanaTerkumpul:  150000000,
				LokasiKejadian: "Lumajang, Jawa Timur",
				Status:         models.StatusTargetTercapai,
				PenggalangID:   penggalangID,
			},
		}

		for _, camp := range campaigns {
			var existing models.Campaign
			if err := db.Where("title = ?", camp.Title).First(&existing).Error; err != nil {
				if db.Create(&camp).Error == nil {
					campaignsCreated++
				}
			}
		}
	}

	return usersCreated, campaignsCreated, nil
}
