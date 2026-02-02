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
	baseURL := "http://localhost:8080/api/v1"

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

	// ... (Category tests above) ...
	fmt.Println("GET /categories/4 (404) - PASS")

	// ================= PRODUCT TESTS =================

	// 7. CREATE Product linked to Category 1
	newProd := map[string]interface{}{
		"name":        "Lego Set",
		"price":       500000,
		"stock":       10,
		"category_id": 1,
	}
	jsonData, _ = json.Marshal(newProd)
	resp, err = http.Post(baseURL+"/products", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating product:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		fmt.Println("Expected 201 Created for Product, got", resp.Status)
		os.Exit(1)
	}
	var createdProd map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&createdProd)
	prodID := int(createdProd["id"].(float64))
	fmt.Println("POST /products - PASS (ID:", prodID, ")")

	// 8. GET Product By ID (Expect Category Name to be present)
	resp, err = http.Get(fmt.Sprintf("%s/products/%d", baseURL, prodID))
	if err != nil {
		fmt.Println("Error getting product:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Expected 200 OK for Product, got", resp.Status)
		os.Exit(1)
	}
	var fetchedProd map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&fetchedProd)

	if fetchedProd["category_name"] == "" || fetchedProd["category_name"] == nil {
		fmt.Println("FAIL: category_name is missing in product response (JOIN failed?)")
		os.Exit(1)
	}
	fmt.Println("GET /products/{id} - PASS (Category Name:", fetchedProd["category_name"], ")")

	// 9. SOFT DELETE Product
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/products/%d", baseURL, prodID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error deleting product:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 204 {
		fmt.Println("Expected 204 No Content for Product Delete, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("DELETE /products/{id} - PASS (Soft Delete)")

	// 10. GET Deleted Product (Should be 404)
	resp, err = http.Get(fmt.Sprintf("%s/products/%d", baseURL, prodID))
	if err != nil {
		fmt.Println("Error getting deleted product:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 404 {
		fmt.Println("Expected 404 Not Found for deleted product, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /products/{id} (404) - PASS")

	// 11. GET Health Check
	resp, err = http.Get(baseURL + "/health")
	if err != nil {
		fmt.Println("Error getting health check:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("Expected 200 OK for health check, got", resp.Status)
		os.Exit(1)
	}
	fmt.Println("GET /health - PASS")

	fmt.Println("ALL TESTS PASSED")
}
