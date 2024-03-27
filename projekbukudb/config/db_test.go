package config_test

import (
	"fmt"
	"ryan-projekbukudb/config"
	"testing"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env not found, using system env")
	}
}

func TestKoneksi(t *testing.T) {
	Init()
	config.OpenDB()
}
