package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.User{},
		&ds.Medicine{},
		&ds.Component{},
		&ds.MedicineProduction{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
