package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"mini_3/app"
	"mini_3/configs"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("env tidak ditemukan, menggunakan env sistem")
	}
	configs.OpenDB(true)
}

/* runner */
func main() {
	Init()
	app.Menu()
}
