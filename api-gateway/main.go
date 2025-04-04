package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/mux"
)

func newProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Error parsing proxy URL: %v", err)
	}
	return httputil.NewSingleHostReverseProxy(url)
}

func main() {
	productService := os.Getenv("PRODUCT_SERVICE_URL")
	orderService := os.Getenv("ORDER_SERVICE_URL")

	if productService == "" {
		productService = "http://localhost:8082"
	}
	if orderService == "" {
		orderService = "http://localhost:8083"
	}

	productProxy := newProxy(productService)
	orderProxy := newProxy(orderService)

	router := mux.NewRouter()

	router.PathPrefix("/products").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		productProxy.ServeHTTP(w, r)
	})
	router.PathPrefix("/orders").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orderProxy.ServeHTTP(w, r)
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Gateway is healthy"))
	})

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
