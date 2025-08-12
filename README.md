# Potion Shop

A simple Go + SQLite web app for buying potions, tracking gold, and managing inventory.

## Run with Docker
```bash
docker build -t potion-shop .
docker run -p 8080:8080 -v $(pwd)/data:/app/data potion-shop
```

Visit: http://localhost:8080

## Run locally (requires Go)
```bash
go run main.go
```
Visit: http://localhost:8080
