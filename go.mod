module github.com/timokae/boot.dev-chirpy

replace github.com/timokae/boot.dev-chirpy-database => ./internal/database
replace github.com/timokae/boot.dev-chirpy-auth => ./internal/auth

go 1.22.2

require github.com/timokae/boot.dev-chirpy-database v0.0.0
require github.com/timokae/boot.dev-chirpy-auth v0.0.0
require golang.org/x/crypto v0.23.0 // indirect
