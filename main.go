package main

import (
	"cateogry-api/database"
	"cateogry-api/internal/handler"
	"cateogry-api/internal/repository"
	"cateogry-api/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "cateogry-api/docs" // Import generated docs

	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

//	@title			Category & Product API
//	@version		1.0
//	@description	This is a sample server for managing categories and products.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1
//	@schemes	http

func main() {
	// Setup Configuration
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Println("Error reading config file, using environment variables", err)
		}
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// Setup Database
	// Setup Database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Println("Warning: Failed to initialize database connection:", err)
		// We continue so the health check endpoint can report the error
	} else {
		defer db.Close()
	}

	// Health Check
	healthHandler := handler.NewHealthHandler(db)

	// Category Dependency Injection
	categoryRepo := repository.NewPostgresCategoryRepository(db)
	categorySvc := service.NewCategoryService(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(categorySvc)

	// Product Dependency Injection
	productRepo := repository.NewProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productSvc)

	// Transaction Dependency Injection
	transactionRepo := repository.NewTransactionRepository(db)
	transactionSvc := service.NewTransactionService(transactionRepo)
	transactionHandler := handler.NewTransactionHandler(transactionSvc)

	// API Versioning Setup
	v1Mux := http.NewServeMux()
	v1Mux.HandleFunc("/health", healthHandler.Check)
	categoryHandler.RegisterRoutes(v1Mux)
	productHandler.RegisterRoutes(v1Mux)
	transactionHandler.RegisterRoutes(v1Mux)

	// Main Router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))

	// Swagger Setup (Root Level)
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
	))

	// Start Background Cleanup Routine (Every 24 hours, delete records older than 30 days)
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("Running cleanup routine...")
			duration := 30 * 24 * time.Hour
			if err := categoryRepo.CleanUpOldDeleted(duration); err != nil {
				log.Println("Error cleaning up categories:", err)
			}
			if err := productRepo.CleanUpOldDeleted(duration); err != nil {
				log.Println("Error cleaning up products:", err)
			}
			log.Println("Cleanup routine finished.")
		}
	}()

	addr := ":" + config.Port
	fmt.Println("Server is running on http://localhost" + addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
