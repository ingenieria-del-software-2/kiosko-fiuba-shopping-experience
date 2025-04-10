package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	cartmodel "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/cart/infrastructure/postgresql"
	checkoutmodel "github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/checkout/infrastructure/postgresql"
	"github.com/ingenieria-del-software-2/kiosko-fiuba-shopping-experience/internal/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Default shipping methods to seed the database
var defaultShippingMethods = []checkoutmodel.ShippingMethodModel{
	{
		ID:            uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Name:          "Standard Shipping",
		Description:   "Delivery in 3-5 business days",
		Price:         5.99,
		EstimatedDays: 5,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Name:          "Express Shipping",
		Description:   "Delivery in 1-2 business days",
		Price:         12.99,
		EstimatedDays: 2,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Name:          "Same Day Delivery",
		Description:   "Delivery within 24 hours (select areas only)",
		Price:         19.99,
		EstimatedDays: 1,
		Active:        true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}

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

	// Seed shipping methods if they don't exist
	log.Println("Checking if shipping methods need to be seeded...")

	var count int64
	db.Model(&checkoutmodel.ShippingMethodModel{}).Count(&count)

	if count == 0 {
		log.Println("Seeding default shipping methods...")
		if err := db.Create(&defaultShippingMethods).Error; err != nil {
			log.Fatalf("Failed to seed shipping methods: %v", err)
		}
		log.Println("Shipping methods seeded successfully")
	} else {
		log.Println("Shipping methods already exist, skipping seed")
	}

	log.Println("Migration process completed successfully")
	os.Exit(0)
}
