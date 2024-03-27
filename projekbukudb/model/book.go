package model

import "gorm.io/gorm"

type Buku struct {
	Model
	ISBN      string `json:"isbn"`
	Pengarang string `json:"pengarang"`
	Tahun     uint   `json:"tahun"`
	Judul     string `json:"judul"`
	Gambar    string `json:"gambar"`
	Stok      uint   `json:"stok"`
}

var ListBuku []Buku

func (lb *Buku) Create(db *gorm.DB) error {
	err := db.Model(Buku{}).Create(&lb).Error
	if err != nil {
		return err
	}

	return nil
}

func (lb *Buku) GetById(db *gorm.DB, id uint) (Buku, error) {
	res := Buku{}

	err := db.Model(Buku{}).Where("id = ?", id).Take(&res).Error
	if err != nil {
		return Buku{}, err
	}

	return res, nil
}

func (lb *Buku) GetAll(db *gorm.DB) ([]Buku, error) {
	res := []Buku{}

	err := db.Model(Buku{}).Find(&res).Error
	if err != nil {
		return []Buku{}, err
	}

	return res, nil
}

func (lb *Buku) UpdateOneByID(db *gorm.DB, id uint) error {
	//err := db.Model(Buku{}).Where("id = ?", lb.Model.ID).Updates(&lb).Error
	err := db.Model(Buku{}).
		Select("isbn", "pengarang", "tahun", "judul", "gambar", "stok").
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"isbn":      lb.ISBN,
			"pengarang": lb.Pengarang,
			"tahun":     lb.Tahun,
			"judul":     lb.Judul,
			"gambar":    lb.Gambar,
			"stok":      lb.Stok,
		}).Error

	if err != nil {
		return err
	}

	return nil
}

func (lb *Buku) DeleteByID(db *gorm.DB, id uint) error {
	err := db.Model(Buku{}).Where("id = ?", id).Delete(&lb).Error
	if err != nil {
		return err
	}
	return nil
}
