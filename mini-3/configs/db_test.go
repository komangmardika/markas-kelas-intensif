package configs_test

import (
	"fmt"
	"github.com/joho/godotenv"
	"mini_3/configs"
	"testing"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env tidak ditemukan, menggunakan env sistem")
	}
}

func TestConnection(t *testing.T) {
	Init()
	configs.OpenDB(false)
}
