//go:build ignore

package main

import (
	"cateogry-api/database"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func main() {
	// Setup Configuration (Same logic as main.go to ensure consistency)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Println("Error reading config file:", err)
		}
	}

	dbConn := viper.GetString("DB_CONN")
	if dbConn == "" {
		fmt.Println("❌ Error: DB_CONN environment variable is not set.")
		return
	}

	// Mask password for display
	displayConn := dbConn
	if parts := strings.Split(dbConn, ":"); len(parts) > 2 {
		// Very basic masking assuming postgres://user:pass@host format
		if atIndex := strings.LastIndex(dbConn, "@"); atIndex != -1 {
			// Find the last colon before @
			if passStart := strings.LastIndex(dbConn[:atIndex], ":"); passStart != -1 {
				displayConn = dbConn[:passStart+1] + "****" + dbConn[atIndex:]
			}
		}
	}
	fmt.Println("Testing connection to:", displayConn)

	// Attempt Connection
	db, err := database.InitDB(dbConn)
	if err != nil {
		fmt.Println("\n❌ DATABASE CONNECTION FAILED")
		fmt.Println("Error Details:", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("\n✅ DATABASE CONNECTED SUCCESSFULLY")
}
