package main

import (
	"fmt"
	"os"

	v1 "github.com/rinsuki/sea-go/api/v1"
	"github.com/rinsuki/sea-go/db"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	fmt.Println("Sea Go Backend")

	_ = db.GetConnection()

	r := gin.Default()
	v1.RegisterToRouter(r.Group("/api/v1"))

	r.Run(os.Getenv("PORT"))
}
