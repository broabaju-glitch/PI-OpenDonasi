package models

import "gorm.io/gorm"

type Role string

const (
	RoleAdmin      Role = "admin"
	RolePenggalang Role = "penggalang"
	RoleDonatur    Role = "donatur"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
	// Profile fields
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	BankName    string `json:"bank_name"`
	BankAcc     string `json:"bank_acc"`
	BankAccName string `json:"bank_acc_name"`
}
