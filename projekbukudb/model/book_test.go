package model_test

import (
	"fmt"
	"ryan-projekbukudb/config"
	"ryan-projekbukudb/model"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func Init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("env not found, using system env")
	}
	config.OpenDB()
}

func TestCreateCarSuccess(t *testing.T) {
	Init()

	libraryData := model.Buku{
		ISBN:      "123",
		Pengarang: "Veda",
		Tahun:     1945,
		Judul:     "Ayam",
		Gambar:    "goreng",
		Stok:      1,
	}

	err := libraryData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	//fmt.Println(libraryData)

	config.Mysql.DB.Unscoped().Delete(&libraryData)
}

func TestGetByIdSuccess(t *testing.T) {
	Init()

	libraryData := model.Buku{
		Model: model.Model{
			ID: 1,
		},
	}

	data, err := libraryData.GetById(config.Mysql.DB, libraryData.ID)
	assert.Nil(t, err)

	fmt.Println(data)
}

func TestGetAll(t *testing.T) {
	Init()

	libraryData := model.Buku{
		ISBN:      "2525",
		Pengarang: "Ryan Teja",
		Tahun:     2007,
		Judul:     "a",
		Gambar:    "abc",
		Stok:      25,
	}

	err := libraryData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	res, err := libraryData.GetAll(config.Mysql.DB)
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(res), 1)

	fmt.Println(res)

	config.Mysql.DB.Unscoped().Delete(&libraryData)
}

func TestUpdateByID(t *testing.T) {
	Init()

	libraryData := model.Buku{
		ISBN:      "1",
		Pengarang: "Teja",
		Tahun:     20124,
		Judul:     "bebek",
		Gambar:    "panggang",
		Stok:      200,
	}

	err := libraryData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	libraryData.Judul = "Ayam goreng"

	err = libraryData.UpdateOneByID(config.Mysql.DB, libraryData.ID)
	assert.Nil(t, err)

	config.Mysql.DB.Unscoped().Delete(&libraryData)
}

func TestDeleteByID(t *testing.T) {
	Init()

	libraryData := model.Buku{
		ISBN:      "QWE12",
		Pengarang: "Jajang Suherman",
		Tahun:     2019,
		Judul:     "Spiderman",
		Gambar:    "man.png",
		Stok:      120,
	}

	err := libraryData.Create(config.Mysql.DB)
	assert.Nil(t, err)

	libraryData = model.Buku{
		Model: model.Model{
			ID: libraryData.ID,
		},
	}

	err = libraryData.DeleteByID(config.Mysql.DB, libraryData.ID)
	assert.Nil(t, err)

	config.Mysql.DB.Unscoped().Delete(&libraryData)
}
