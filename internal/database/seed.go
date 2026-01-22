package database

import (
	"seal-of-ash/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func SeedUsers() {
	var count int64
	ActiveDB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []models.User{
		{Username: "kasir", Role: "user"},
		{Username: "forensic", Role: "forensic"},
		{Username: "admin", Role: "admin"},
	}

	for _, u := range users {
		hash, _ := bcrypt.GenerateFromPassword([]byte("password"), 12)
		u.Password = string(hash)
		ActiveDB.Create(&u)
	}
}
