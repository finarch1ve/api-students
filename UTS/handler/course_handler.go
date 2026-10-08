package handler

import (
	"strconv"
	"strings"

	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CourseHandler struct {
	DB *gorm.DB
}

type courseView struct {
	ID        uint   `json:"id"`
	KodeMK    string `gorm:"column:kode_mk" json:"kode_mk"`
	NamaMK    string `gorm:"column:nama_mk" json:"nama_mk"`
	SKS       int    `gorm:"column:sks" json:"sks"`
	Semester  int    `json:"semester"`
	Kuota     int    `json:"kuota"`
	Terisi    int    `json:"terisi"`
	SisaKuota int    `gorm:"-" json:"sisa_kuota"`
}

// GET /courses
func (h *CourseHandler) List(c *fiber.Ctx) error {
	q := h.DB.Table("courses AS c").
		Select("c.id, c.kode_mk, c.nama_mk, c.sks, c.semester, c.kuota, COUNT(e.id) AS terisi")

	if ta := strings.TrimSpace(c.Query("tahun_akademik")); ta != "" {
		q = q.Joins("LEFT JOIN enrollments e ON e.course_id = c.id AND e.tahun_akademik = ?", ta)
	} else {
		q = q.Joins("LEFT JOIN enrollments e ON e.course_id = c.id")
	}

	if v := strings.TrimSpace(c.Query("semester")); v != "" {
		s, err := strconv.Atoi(v)
		if err != nil {
			return helper.ValidationFailed(map[string][]string{"semester": {"harus berupa angka"}})
		}
		q = q.Where("c.semester = ?", s)
	}
	if v := strings.TrimSpace(c.Query("search")); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		q = q.Where("(LOWER(c.kode_mk) LIKE ? OR LOWER(c.nama_mk) LIKE ?)", like, like)
	}

	q = q.Group("c.id")
	if c.Query("available") == "true" {
		q = q.Having("COUNT(e.id) < c.kuota")
	}

	var rows []courseView
	if err := q.Order("c.id ASC").Scan(&rows).Error; err != nil {
		return err
	}
	for i := range rows {
		sisa := rows[i].Kuota - rows[i].Terisi
		if sisa < 0 {
			sisa = 0
		}
		rows[i].SisaKuota = sisa
	}
	if rows == nil {
		rows = []courseView{}
	}

	return helper.Success(c, fiber.StatusOK, "Data mata kuliah berhasil diambil", rows)
}