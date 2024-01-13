package main

import (
	"log"
	"os"

	"R_I_P_labs/internal/app/ds"
	"R_I_P_labs/internal/app/dsn"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

//func main() {
//	_ = godotenv.Load()
//	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//
//	// Migrate the schema
//	err = db.AutoMigrate(
//		&ds.User{},
//		&ds.Medicine{},
//		&ds.Component{},
//		&ds.MedicineProduction{},
//	)
//	if err != nil {
//		panic("cant migrate db")
//	}
//}

func main() {
	_ = godotenv.Load()

	// Настройка уровня логирования
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // Вывод в консоль
		logger.Config{
			SlowThreshold:             time.Second, // Порог для определения медленных запросов
			LogLevel:                  logger.Info, // Уровень логирования
			IgnoreRecordNotFoundError: true,        // Игнорировать ошибку "record not found"
			Colorful:                  true,        // Цветное логирование
		},
	)

	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{
		Logger: newLogger, // Установка созданного логгера
	})
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
