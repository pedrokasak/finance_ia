package database

import (
	"finance-ia/internal/domain/user"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestConnection(databaseURL string) bool {

	db, err := Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get generic database object: %v", err)
		return false
	}
	defer sqlDB.Close()

	var version string
	if err := sqlDB.QueryRow("SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)
	return true
}

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &user.User{},
      
    )
}
