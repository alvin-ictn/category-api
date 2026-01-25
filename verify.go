//go:build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	baseURL := "http://localhost:8080"

	// Wait for server to start
	time.Sleep(2 * time.Second)

	fmt.Println("Starting Verification...")

	// 1. GET /categories
	resp, err := http.Get(baseURL + "/categories")
	if err != nil {
		fmt.Println("Error getting categories:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Expected 200 OK, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /categories - PASS")

	// 2. GET /categories/1
	resp, err = http.Get(baseURL + "/categories/1")
	if err != nil {
		fmt.Println("Error getting category 1:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Expected 200 OK for cat 1, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /categories/1 - PASS")

	// 3. POST /categories
	newCat := map[string]string{"name": "Toys", "description": "Fun stuff"}
	jsonData, _ := json.Marshal(newCat)
	resp, err = http.Post(baseURL+"/categories", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating category:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		fmt.Println("Expected 201 Created, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("POST /categories - PASS")

	// 4. GET /categories/4 (assuming ID 4)
	resp, err = http.Get(baseURL + "/categories/4")
	if err != nil {
		fmt.Println("Error getting category 4:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Expected 200 OK for cat 4, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /categories/4 - PASS")

	// 5. DELETE /categories/4
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/categories/4", nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error deleting category 4:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		fmt.Println("Expected 204 No Content, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("DELETE /categories/4 - PASS")

	// 6. GET /categories/4 (Should be 404)
	resp, err = http.Get(baseURL + "/categories/4")
	if err != nil {
		fmt.Println("Error getting category 4 again:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		fmt.Println("Expected 404 Not Found for cat 4, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /categories/4 (404) - PASS")

	fmt.Println("ALL TESTS PASSED")
}
