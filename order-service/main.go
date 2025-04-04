package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"

	_ "order-service/docs"
)

// @title Order Service API
// @version 1.0
// @description API for managing orders in a microservices demo.
// @host localhost:8083
// @BasePath /
type OrderRequest struct {
	ProductID int `json:"product_id"`
}

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Order struct {
	ID        int `json:"id"`
	ProductID int `json:"product_id"`
}

// ProductCache struct to hold cached products
type ProductCache struct {
	Products []Product
	mu       sync.RWMutex
}

var cache = &ProductCache{}

func updateCache(products []Product) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.Products = products
	log.Println("Updated product cache")
}

func getCachedProduct(productID int) (Product, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	for _, p := range cache.Products {
		if p.ID == productID {
			return p, true
		}
	}
	return Product{}, false
}

// PlaceOrder godoc
// @Summary Place a new order
// @Description Create an order by checking product availability
// @Tags orders
// @Accept json
// @Produce json
// @Param order body OrderRequest true "Order request"
// @Success 200 {string} string "Order placed for {product_name} (ID: {product_id})!"
// @Failure 400 {string} string "Invalid order"
// @Failure 404 {string} string "Product not found"
// @Failure 503 {string} string "Product service unavailable and no cache available"
// @Router /orders [post]
func placeOrder(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var order OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid order", http.StatusBadRequest)
			return
		}

		resp, err := http.Get("http://localhost:8080/products")
		if err == nil && resp.StatusCode == http.StatusOK {

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				http.Error(w, "Error reading products", http.StatusInternalServerError)
				return
			}

			var products []Product
			if err := json.Unmarshal(body, &products); err != nil {
				http.Error(w, "Error parsing products", http.StatusInternalServerError)
				return
			}

			updateCache(products)

			for _, p := range products {
				if p.ID == order.ProductID {
					_, err := db.Exec("INSERT INTO orders (product_id) VALUES ($1)", order.ProductID)
					if err != nil {
						log.Printf("Error inserting order: %v", err)
						http.Error(w, "Error saving order", http.StatusInternalServerError)
						return
					}
					fmt.Fprintf(w, "Order placed for %s (ID: %d)!", p.Name, order.ProductID)
					return
				}
			}
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		log.Println("Product Service unavailable, falling back to cache")
		if cachedProduct, found := getCachedProduct(order.ProductID); found {
			_, err := db.Exec("INSERT INTO orders (product_id) VALUES ($1)", order.ProductID)
			if err != nil {
				log.Printf("Error inserting order: %v", err)
				http.Error(w, "Error saving order", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "Order placed for %s (ID: %d)", cachedProduct.Name, order.ProductID)
			return
		}

		http.Error(w, "Product service unavailable and no cache available", http.StatusServiceUnavailable)
	}
}

func main() {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=Admin@123 dbname=orders_db sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS orders (id SERIAL PRIMARY KEY, product_id INT)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	http.HandleFunc("/orders", placeOrder(db))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Order Service running on :8083")
	log.Fatal(http.ListenAndServe(":8083", nil))
}
