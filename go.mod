module github.com/timokae/boot.dev-chirpy

replace github.com/timokae/boot.dev-chirpy-database => ./internal/database

replace github.com/timokae/boot.dev-chirpy-auth => ./internal/auth

go 1.22.2

require github.com/timokae/boot.dev-chirpy-database v0.0.0

require github.com/timokae/boot.dev-chirpy-auth v0.0.0

require github.com/golang-jwt/jwt/v5 v5.2.1 // indirect

require (
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.23.0 // indirect
)
