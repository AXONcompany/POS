package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/AXONcompany/POS/internal/infrastructure/perisistence/postgres"
)


func loadVariable(key string)string{
	err := godotenv.Load()
	if err != nil{
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main(){

	DB_NAME := loadVariable("DB_NAME")
	DB_PASS := loadVariable("DB_PASSWORD")
	DB_HOST := loadVariable("DB_HOST")
	DB_USER := loadVariable("DB_USER")
	DB_PORT := loadVariable("DB_PORT")

	err := postgres.Connect(DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)
	err2 := postgres.Migrate()

	if err != nil{
		panic("error connecting to database"+err.Error())
	}

	if err2 != nil{
		panic("error connecting to database" + err2.Error())
	}


}