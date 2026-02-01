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

	"github.com/spf13/viper"
)

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

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
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Dependency Injection
	repo := repository.NewPostgresCategoryRepository(db) // Switched to Postgres
	svc := service.NewCategoryService(repo)
	h := handler.NewCategoryHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	addr := ":" + config.Port
	fmt.Println("Server is running on http://localhost" + addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
