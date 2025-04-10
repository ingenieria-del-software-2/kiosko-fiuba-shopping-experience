package main

import (
	"fmt"
	"log"
	"os"

	cartmodel "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/infrastructure/postgresql"
	checkoutmodel "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/infrastructure/postgresql"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to the database with GORM
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DbHost, cfg.DbPort, cfg.DbUser, cfg.DbPassword, cfg.DbName, cfg.DbSslMode,
	)

	log.Printf("Connecting to database: %s", dsn)

	// Open connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Running database migrations...")

	// Auto migrate the models
	err = db.AutoMigrate(
		&cartmodel.CartModel{},
		&checkoutmodel.ShippingAddressModel{},
		&checkoutmodel.ShippingMethodModel{},
		&checkoutmodel.CheckoutModel{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	log.Println("Migration process completed successfully")
	os.Exit(0)
}
