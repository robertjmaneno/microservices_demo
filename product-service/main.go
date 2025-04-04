package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/swag"

	_ "product-service/docs"
)

// @title Product Service API
// @version 1.0
// @description API for managing products in a microservices demo.
// @host localhost:8082
// @BasePath /
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// GetProducts godoc
// @Summary Get list of products
// @Description Retrieve all products from the database
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} Product
// @Failure 500 {string} string "Internal Server Error"
// @Router /products [get]
func getProducts(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, price FROM products")
		if err != nil {
			log.Printf("Database query error: %v", err)
			http.Error(w, "Database query failed", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
				log.Printf("Error scanning rows: %v", err)
				http.Error(w, "Error scanning rows", http.StatusInternalServerError)
				return
			}
			products = append(products, p)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func main() {
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=Admin@123 dbname=products_db sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (id SERIAL PRIMARY KEY, name TEXT, price FLOAT)`)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	_, err = db.Exec(`INSERT INTO products (name, price) VALUES ('Laptop', 999.99) ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Error inserting Laptop: %v", err)
	}
	_, err = db.Exec(`INSERT INTO products (name, price) VALUES ('Phone', 499.99) ON CONFLICT DO NOTHING`)
	if err != nil {
		log.Printf("Error inserting Phone: %v", err)
	}

	http.HandleFunc("/products", getProducts(db))

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Product Service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
