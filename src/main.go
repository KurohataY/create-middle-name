package main

import (
	"example.com/m/v2/function/funcDB"
	"example.com/m/v2/server"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv"
)

func main() {
	funcDB.DbInit()
	server.Init()
}
