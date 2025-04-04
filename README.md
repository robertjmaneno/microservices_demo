# Microservices Demo

This repository demonstrates a basic microservices architecture implemented in Go.  
It simulates an e-commerce system with the following services:

- Product Service
- Order Service
- API Gateway

---

<details>
<summary>🧱 Architecture Overview</summary>

**Product Service**  
- Manages a product catalog  
- Uses PostgreSQL for data storage  
- Exposes a REST API  

**Order Service**  
- Handles order placement  
- Fetches product data via the API Gateway  
- Caches product data for resilience  

**API Gateway**  
- Acts as a reverse proxy  
- Routes client requests to the appropriate microservice  

</details>

---

<details>
<summary>🔧 Prerequisites</summary>

You must have the following installed:

- Go version 1.16 or higher  
- PostgreSQL version 12 or higher  
- Git
- Visual Studio Code

</details>

---

<details>
<summary>🚀 Setup Instructions</summary>

1. Clone the repository  

git clone https://github.com/robertjmaneno/microservices_demo.git
cd microservices-demo



2. Start PostgreSQL server and create databases  

psql -U postgres

CREATE DATABASE products_db;
CREATE DATABASE orders_db;
\q

markdown
Copy
Edit

Ensure the PostgreSQL user `postgres` has the password `Admin@123`, or update the code accordingly.

3. Install Go dependencies  

go mod init microservices-demo
go get github.com/lib/pq
go get github.com/gorilla/mux
go get github.com/swaggo/http-swagger
go get github.com/swaggo/swag

markdown
Copy
Edit

4. Start each service in separate terminal windows  

- Product Service  

go run product-service.go



- Order Service  

go run order-service.go



- API Gateway  

go run api-gateway.go



</details>

---

<details>
<summary>📁 Project Structure</summary>

microservices-demo/
├── api-gateway.go
├── product-service.go
├── order-service.go
├── README.md


</details>

---

<details>
<summary>📌 Service Ports</summary>

- Product Service: http://localhost:8082  
- Order Service: http://localhost:8083  
- API Gateway: http://localhost:8080  

</details>

---

<details>
<summary>🧪 Testing the Services</summary>

**1. Get All Products via API Gateway**

curl http://localhost:8080/products



Expected Output:  
[
{id: 1, name: Laptop, price: 999.99},
{id: 2, name: Phone, price: 499.99}
]



**2. Place an Order**

curl -X POST http://localhost:8080/orders -H Content-Type: application/json -d {product_id: 1}



Expected Output:  
Order placed for Laptop (ID: 1)!


**3. Use Swagger UI**

- Product Service Swagger: http://localhost:8082/swagger/index.html  
- Order Service Swagger: http://localhost:8083/swagger/index.html  

**4. Test Caching / Resilience**

Stop Product Service (`Ctrl + C`)  
Place another order:  

curl -X POST http://localhost:8080/orders -H Content-Type: application/json -d {product_id: 2}



Expected Output:  
Order placed for Phone (ID: 2) using cached data!



**5. API Gateway Health Check**

curl http://localhost:8080/health



Expected Output:  
API Gateway is healthy



</details>

---

<details>
<summary>⚠️ Troubleshooting</summary>

- Make sure PostgreSQL server is running  
- Ensure databases are created correctly  
- Check credentials match code  
- Ensure ports 8080, 8082, 8083 are free  
- Run `go get` if any dependency errors occur  

</details>

---

<details>
<summary>📄 License</summary>

This project is licensed under the MIT License.  

</details>

---

<details>
<summary>🤝 Contributing</summary>

- Fork the repository  
- Create issues for bugs or features  
- Submit pull requests with improvements  

</details>

---

<details>
<summary>🔐 Credentials</summary>

Default PostgreSQL credentials in use:

- Username: postgres  
- Password: Admin@123

Update the code if your setup is different.  

