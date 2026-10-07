package database

import (
	"fmt"
	"log"

	"siakad-mini/app/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func hash(p string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	return string(b)
}

func Seed(db *gorm.DB) error {
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count > 0 {
		log.Println("Seeder dilewati: data sudah ada")
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 1 admin
		admin := model.User{Email: "admin@siakad.test", Password: hash("admin12345"), Role: "admin"}
		if err := tx.Create(&admin).Error; err != nil {
			return err
		}

		// 20 mahasiswa (password awal = NIM)
		nama := []string{
			"Rina Putri", "Budi Santoso", "Citra Lestari", "Dimas Pratama", "Eka Wulandari",
			"Fajar Nugroho", "Gita Permata", "Hadi Wijaya", "Indah Sari", "Joko Susilo",
			"Kartika Dewi", "Lukman Hakim", "Maya Anggraini", "Naufal Rizky", "Olivia Putri",
			"Putra Ramadhan", "Qonita Aulia", "Rizal Firmansyah", "Salsabila Nur", "Teguh Prakoso",
		}
		prodi := []string{"Teknik Informatika", "Sistem Informasi", "Teknik Komputer"}
		ipk := []float64{
			3.85, 3.62, 3.45, 3.30, 3.15, 3.05, 3.00, 3.72, 3.50, 3.20,
			2.95, 2.80, 2.75, 2.60, 2.50, 2.99, 2.70, 2.55,
			2.40, 2.10,
		}

		for i := 0; i < 20; i++ {
			nim := fmt.Sprintf("434241%06d", i+1)
			user := model.User{
				Email:    fmt.Sprintf("mhs%02d@siakad.test", i+1),
				Password: hash(nim),
				Role:     "mahasiswa",
			}
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			student := model.Student{
				UserID:      user.ID,
				NIM:         nim,
				Nama:        nama[i],
				Prodi:       prodi[i%3],
				Angkatan:    2022 + i%4,
				IPKTerakhir: ipk[i],
			}
			if err := tx.Create(&student).Error; err != nil {
				return err
			}
		}

		// 10 mata kuliah (2 di antaranya kuota kecil untuk tes kuota penuh)
		courses := []model.Course{
			{KodeMK: "IF101", NamaMK: "Algoritma dan Pemrograman", SKS: 3, Semester: 1, Kuota: 30},
			{KodeMK: "IF102", NamaMK: "Basis Data", SKS: 3, Semester: 2, Kuota: 30},
			{KodeMK: "IF103", NamaMK: "Struktur Data", SKS: 3, Semester: 2, Kuota: 30},
			{KodeMK: "IF104", NamaMK: "Jaringan Komputer", SKS: 3, Semester: 3, Kuota: 25},
			{KodeMK: "IF105", NamaMK: "Pemrograman Web", SKS: 4, Semester: 3, Kuota: 25},
			{KodeMK: "IF106", NamaMK: "Rekayasa Perangkat Lunak", SKS: 3, Semester: 4, Kuota: 25},
			{KodeMK: "IF107", NamaMK: "Pemrograman Back End", SKS: 4, Semester: 4, Kuota: 2},
			{KodeMK: "IF108", NamaMK: "Keamanan Siber", SKS: 3, Semester: 5, Kuota: 2},
			{KodeMK: "IF109", NamaMK: "Kecerdasan Buatan", SKS: 3, Semester: 5, Kuota: 20},
			{KodeMK: "IF110", NamaMK: "Etika Profesi", SKS: 2, Semester: 6, Kuota: 30},
		}
		return tx.Create(&courses).Error
	})
}