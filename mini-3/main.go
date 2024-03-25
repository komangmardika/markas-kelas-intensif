package main

/*
 * author: I Komang Mardika
 * email: komang.mardika@hotmail.com
 * description: mini project 3 for markasbali ft kominfo golang hacker workshop
 *
 */
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
