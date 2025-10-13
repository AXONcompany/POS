package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
    "gorm.io/gorm"
	"github.com/AXONcompany/POS/internal/domain/product"
	"github.com/AXONcompany/POS/internal/domain/order"
	"github.com/AXONcompany/POS/internal/domain/table"
	"github.com/AXONcompany/POS/internal/domain/user"
)


var DB *gorm.DB

func Connect(host,user,password,dbname,port string)error{
		
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil{
		return fmt.Errorf("error connecting to database %w", err)
	}

	DB = db

	return nil
}


func Migrate()error{
	return DB.AutoMigrate(
		&product.Category{},
		&product.Product{},
		&order.Order{},
		&table.Table{},
		&user.User{},
	)
}