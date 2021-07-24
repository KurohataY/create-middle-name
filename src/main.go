package main

import (
	"example.com/m/v2/database"
	"example.com/m/v2/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv"
)

func main() {
	database.Init()
	server.Init()
}
