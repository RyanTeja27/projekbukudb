package main

import (
	"bufio"
	"fmt"
	"os"
	"ryan-projekbukudb/config"
	"ryan-projekbukudb/model"
	"sync"
	"time"

	"github.com/MasterDimmy/go-cls"
	"github.com/go-pdf/fpdf"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("env not found, using system env")
	}
}

func TambahBuku(db *gorm.DB) {
	isbn := ""
	JudulB := ""
	PengarangB := ""
	gambar := ""
	var TahunTerbit uint
	var stock uint

	fmt.Println("")
	fmt.Println("Tambahkan Buku")
	fmt.Println("")
	draftBuku := []model.Buku{}

	for {

		fmt.Print("Masukan Kode ISBN : ")
		_, err := fmt.Scanln(&isbn)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Print("Masukan Judul Buku : ")
		_, err = fmt.Scanln(&JudulB)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		fmt.Print("Masukan Nama Pengarang : ")
		_, err = fmt.Scanln(&PengarangB)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
		}

		fmt.Print("Masukan Tahun Terbit : ")
		_, err = fmt.Scanln(&TahunTerbit)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Print("Masukan Gambar Buku : ")
		_, err = fmt.Scanln(&gambar)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Print("Masukan  Jumlah Stock : ")
		_, err = fmt.Scanln(&stock)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		draftBuku = append(draftBuku, model.Buku{
			ISBN:      isbn,
			Judul:     JudulB,
			Pengarang: PengarangB,
			Gambar:    gambar,
			Tahun:     TahunTerbit,
			Stok:      stock,
		})

		pilihanTambahBuku := 0
		fmt.Println("Ketik 1 untuk tambah buku lain, ketik 0 untuk selesai")
		_, err = fmt.Scanln(&pilihanTambahBuku)
		if err != nil {
			fmt.Println("Terjadi Error: ", err)
			return
		}

		if pilihanTambahBuku == 0 {
			break
		}
	}

	fmt.Println("Menambah Buku...")

	_ = os.Mkdir("books", 0777)

	ch := make(chan model.Buku)

	wg := sync.WaitGroup{}

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go SimpanBuku(ch, &wg, i, db)
	}

	for _, bukuTersimpan := range draftBuku {
		ch <- bukuTersimpan
	}

	close(ch)

	wg.Wait()

	fmt.Println("Berhasil Menambahkan Buku")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func SimpanBuku(ch <-chan model.Buku, wg *sync.WaitGroup, noPustakawan int, db *gorm.DB) {

	for bukuTersimpan := range ch {
		if err := db.Create(&bukuTersimpan).Error; err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		fmt.Printf("Pustakawan No %d Memproses Kode Buku : %s!\n", noPustakawan, bukuTersimpan.Judul)
	}
	wg.Done()
}

func LihatDaftarBuku(ch <-chan string, chBuku chan model.Buku, wg *sync.WaitGroup, db *gorm.DB) {
	var buku model.Buku
	for bukuISBN := range ch {
		if err := db.Where("isbn = ?", bukuISBN).Find(&buku).Error; err != nil {
			fmt.Println("Terjadi error res:", err)
			continue
		}

		chBuku <- buku
	}
	wg.Done()
}

func LihatList(db *gorm.DB) {
	fmt.Println("")
	fmt.Println("Lihat List Buku")
	fmt.Println("")
	ListBuku := model.Buku{}

	res, err := ListBuku.GetAll(db)
	if err != nil {
		fmt.Println("Terjadi error:", err)
		return
	}

	wg := sync.WaitGroup{}

	ch := make(chan string)
	chBuku := make(chan model.Buku, len(model.ListBuku))

	jumlahPustakawan := 5

	for i := 0; i < jumlahPustakawan; i++ {
		wg.Add(1)
		go LihatDaftarBuku(ch, chBuku, &wg, db)
	}

	for _, fileBuku := range res {
		ch <- fileBuku.ISBN
	}

	close(ch)

	wg.Wait()

	close(chBuku)

	if len(model.ListBuku) < 1 {
		fmt.Println("===== Tidak ada buku =====")
	}

	for i, v := range model.ListBuku {
		i++
		fmt.Printf("%d. ID : %d, ISBN : %s, Penulis : %s, Tahun : %d, Judul : %s, Gambar : %s, Stok : %d\n", i, v.ID, v.ISBN, v.Pengarang, v.Tahun, v.Judul, v.Gambar, v.Stok)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func DetailBuku(db *gorm.DB, ID uint) {
	fmt.Println("")
	fmt.Println("Detail Buku")
	fmt.Println("")

	ListBuku := model.Buku{}

	buku, err := ListBuku.GetById(db, ID)
	if err != nil {
		fmt.Println("Terjadi error:", err)
		return
	}
	fmt.Printf("Kode ISBN : %s\n", buku.ISBN)
	fmt.Printf("Judul Buku : %s\n", buku.Judul)
	fmt.Printf("Pengarang Buku : %s\n", buku.Pengarang)
	fmt.Printf("Tahun Terbit Buku : %d\n", buku.Tahun)
	fmt.Printf("Gambar : %s\n", buku.Gambar)
	fmt.Printf("Stock Buku : %d\n", buku.Stok)

}

func HapusBuku(db *gorm.DB, ID uint) {

	var isiBuku bool
	for _, buku := range model.ListBuku {
		if buku.ID == ID {
			isiBuku = true
			err := buku.DeleteByID(db, ID)
			if err != nil {
				fmt.Println("Terjadi error:", err)
				return
			}

			fmt.Print("\n")
			fmt.Println("Buku Berhasil Dihapus")
			break
		}
	}

	if !isiBuku {

		fmt.Print("\n")
		fmt.Println("Kode Buku Salah Atau Tidak Ada")
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func GeneratedPdfBuku(db *gorm.DB) {

	LihatList(db)
	fmt.Println("===== Membuat Daftar Buku =====")
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)

	for i, buku := range model.ListBuku {
		bukuText := fmt.Sprintf(
			"Buku #%d:\nID : %d\nISBN : %s\nPenulis : %s\nTahun : %d\nJudul : %s\nGambar : %s\nStok : %d\n",
			i+1, buku.ID, buku.ISBN,
			buku.Pengarang, buku.Tahun, buku.Judul, buku.Gambar,
			buku.Stok)

		pdf.MultiCell(0, 10, bukuText, "0", "L", false)
		pdf.Ln(5)
	}

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("daftar_buku_%s.pdf",
			time.Now().Format("2006-01-02 15-04-05")))

	if err != nil {
		fmt.Println("Terjadi Error: ", err)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func EditBuku(db *gorm.DB, ID uint) {

	DetailBuku(db, ID)

	fmt.Println("")
	fmt.Println("Edit Buku")
	fmt.Println("")

	var buku model.Buku

	err := db.Where("id = ?", ID).First(&buku).Error
	if err != nil {
		fmt.Println("Terjadi kesalahan:", err)
		return
	}

	fmt.Print("Masukan Kode ISBN : ")
	_, err = fmt.Scanln(&buku.ISBN)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Judul Buku : ")
	_, err = fmt.Scanln(&buku.Judul)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukan Nama Pengarang : ")
	_, err = fmt.Scanln(&buku.Pengarang)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	fmt.Print("Masukkan Gambar : ")
	_, err = fmt.Scanln(&buku.Gambar)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukkan Stok : ")
	_, err = fmt.Scanln(&buku.Stok)
	if err != nil {
		fmt.Println("Terjadi Kesalahan : ", err)
		return
	}

	fmt.Print("Masukan Tahun Terbit : ")
	_, err = fmt.Scanln(&buku.Tahun)
	if err != nil {
		fmt.Println("Terjadi Error : ", err)
	}

	err = buku.UpdateOneByID(db, buku.ID)
	if err != nil {
		fmt.Println("Terjadi Eror", err)
		return
	}

	fmt.Print("\nBuku Berhasil Di Edit")
}

func main() {

	Init()
	config.OpenDB()

	PilihanBuku := 0

	cls.CLS()
	fmt.Println("")
	fmt.Println("Aplikasi Manajemen Daftar Buku Perpustakaan")
	fmt.Println("")
	fmt.Println("Silahkan Pilih : ")
	fmt.Println("1. Tambah Buku")
	fmt.Println("2. Liat List Buku")
	fmt.Println("3. Detail Buku")
	fmt.Println("4. Ubah/Edit Buku")
	fmt.Println("5. Hapus Buku")
	fmt.Println("6. Print buku pdf")
	fmt.Println("7. Keluar")
	fmt.Println("")

	fmt.Print("Masukan Pilihan : ")
	_, err := fmt.Scanln(&PilihanBuku)
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	switch PilihanBuku {
	case 1:
		TambahBuku(config.Mysql.DB)
	case 2:
		LihatList(config.Mysql.DB)
	case 3:
		var pilihanDetail uint
		LihatList(config.Mysql.DB)
		fmt.Print("Masukkan Kode Buku : ")
		_, err := fmt.Scanln(&pilihanDetail)
		if err != nil {
			fmt.Println("Terjadi Error : ", err)
			return
		}
		DetailBuku(config.Mysql.DB, pilihanDetail)
	case 4:
		var pilihanEdit uint
		LihatList(config.Mysql.DB)
		fmt.Print("Masukkan Kode Buku Yang akan diedit : ")
		_, err := fmt.Scanln(&pilihanEdit)
		if err != nil {
			fmt.Println("Terjadi Error : ", err)
			return
		}
		EditBuku(config.Mysql.DB, pilihanEdit)
	case 5:
		var pilihanHapus uint
		LihatList(config.Mysql.DB)
		fmt.Print("masukkan kode yang akan dihapus : ")
		_, err := fmt.Scanln(&pilihanHapus)
		if err != nil {
			fmt.Println("Terjadi error: ", err)
			return
		}
		HapusBuku(config.Mysql.DB, pilihanHapus)
	case 6:
		GeneratedPdfBuku(config.Mysql.DB)
	case 7:
		fmt.Println("\nSelesai")
		os.Exit(0)
	default:
		fmt.Println("\ntidak ada opsi")
	}

	main()
}
