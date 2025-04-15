# Go-Gin-Project

This project is a learning tool for exploring basic CRUD (Create, Read, Update, Delete) operations in Go using the Gin framework. It demonstrates clean architecture, dependency injection with Wire, and modular design patterns for scalable Go applications.

---

## 🛠️ Tech Stack

- **Go** – Programming language  
- **Gin** – Web framework  
- **GORM** – ORM for database operations  
- **Wire** – Dependency injection  
- **PostgreSQL** – Relational database  
- **Docker & Docker Compose** – Containerization  
- **Testify** – Testing toolkit  

---

## ⚙️ Prerequisites

- Go 1.16 or higher  
- PostgreSQL  
- (Optional) Docker & Docker Compose  

---

## 📦 Installation

1. Clone the repository

2. Copy environment config and adjust as needed:

   ```bash
   cp example.env .env
   ```

3. Create the database
---

## ▶️ Running the App

```bash
go run cmd/main.go
```

---

## 📡 API Endpoints

| Method | Endpoint        | Description         |
|--------|------------------|---------------------|
| GET    | `/api/tag`       | List all tags       |
| GET    | `/api/tag/:id`   | Get tag by ID       |
| POST   | `/api/tag`       | Create new tag      |
| PUT    | `/api/tag/:id`   | Update tag by ID    |
| DELETE | `/api/tag/:id`   | Delete tag by ID    |

---

## ✅ Testing

Run unit and integration tests:

```bash
go test ./...
```

With coverage report:

```bash
go test -cover ./...
```

Generate HTML coverage report:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html && xdg-open coverage.html
```





