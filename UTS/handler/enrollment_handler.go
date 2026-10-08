package handler

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"siakad-mini/app/model"
	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type EnrollmentHandler struct {
	DB *gorm.DB
}

type CreateEnrollmentRequest struct {
	CourseID      uint   `json:"course_id" validate:"required"`
	TahunAkademik string `json:"tahun_akademik" validate:"required"`
}

var tahunAkademikRe = regexp.MustCompile(`^(\d{4})/(\d{4})-(Ganjil|Genap)$`)

func validTahunAkademik(s string) bool {
	m := tahunAkademikRe.FindStringSubmatch(s)
	if m == nil {
		return false
	}
	a, _ := strconv.Atoi(m[1])
	b, _ := strconv.Atoi(m[2])
	return b == a+1
}

// POST /enrollments
func (h *EnrollmentHandler) Create(c *fiber.Ctx) error {
	var req CreateEnrollmentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Unprocessable("Body request tidak valid")
	}
	req.TahunAkademik = strings.TrimSpace(req.TahunAkademik)

	errs := helper.ValidateStruct(req)
	if errs == nil {
		errs = map[string][]string{}
	}
	if _, bad := errs["tahun_akademik"]; !bad && !validTahunAkademik(req.TahunAkademik) {
		errs["tahun_akademik"] = []string{"format harus seperti 2026/2027-Ganjil"}
	}
	if len(errs) > 0 {
		return helper.ValidationFailed(errs)
	}

	userID, _ := c.Locals("user_id").(uint)
	var result fiber.Map

	err := h.DB.Transaction(func(tx *gorm.DB) error {
		var st model.Student
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).First(&st).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return helper.Forbidden("Data mahasiswa tidak ditemukan")
			}
			return err
		}

		var course model.Course
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&course, req.CourseID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return helper.ValidationFailed(map[string][]string{"course_id": {"mata kuliah tidak ditemukan"}})
			}
			return err
		}

		var n int64
		if err := tx.Model(&model.Enrollment{}).
			Where("student_id = ? AND course_id = ? AND tahun_akademik = ?", st.ID, course.ID, req.TahunAkademik).
			Count(&n).Error; err != nil {
			return err
		}
		if n > 0 {
			return helper.Conflict("Mata kuliah ini sudah pernah diambil pada tahun akademik yang sama")
		}

		if err := tx.Model(&model.Enrollment{}).
			Where("course_id = ? AND tahun_akademik = ?", course.ID, req.TahunAkademik).
			Count(&n).Error; err != nil {
			return err
		}
		if int(n) >= course.Kuota {
			return helper.Unprocessable("Kuota mata kuliah sudah penuh")
		}

		var total int64
		if err := tx.Table("enrollments AS e").
			Select("COALESCE(SUM(c.sks), 0)").
			Joins("JOIN courses c ON c.id = e.course_id").
			Where("e.student_id = ? AND e.tahun_akademik = ?", st.ID, req.TahunAkademik).
			Scan(&total).Error; err != nil {
			return err
		}
		batas := helper.MaxSKS(st.IPKTerakhir)
		if int(total)+course.SKS > batas {
			sisa := batas - int(total)
			if sisa < 0 {
				sisa = 0
			}
			return helper.Unprocessable(fmt.Sprintf(
				"Total SKS melebihi batas %d SKS (IPK %.2f). SKS terpakai %d, sisa %d SKS, sedangkan mata kuliah ini %d SKS",
				batas, st.IPKTerakhir, total, sisa, course.SKS))
		}

		e := model.Enrollment{StudentID: st.ID, CourseID: course.ID, TahunAkademik: req.TahunAkademik}
		if err := tx.Create(&e).Error; err != nil {
			return err
		}

		result = fiber.Map{
			"id":             e.ID,
			"student_id":     e.StudentID,
			"course_id":      e.CourseID,
			"tahun_akademik": e.TahunAkademik,
			"created_at":     e.CreatedAt,
			"course": fiber.Map{
				"kode_mk": course.KodeMK,
				"nama_mk": course.NamaMK,
				"sks":     course.SKS,
			},
			"total_sks": int(total) + course.SKS,
			"batas_sks": batas,
		}
		return nil
	})
	if err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusCreated, "Mata kuliah berhasil ditambahkan ke KRS", result)
}

// DELETE /enrollments/:id
func (h *EnrollmentHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return helper.NotFound("Data KRS tidak ditemukan")
	}

	var e model.Enrollment
	if err := h.DB.First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NotFound("Data KRS tidak ditemukan")
		}
		return err
	}

	userID, _ := c.Locals("user_id").(uint)
	var st model.Student
	if err := h.DB.Where("user_id = ?", userID).First(&st).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.Forbidden("Data mahasiswa tidak ditemukan")
		}
		return err
	}
	if e.StudentID != st.ID {
		return helper.Forbidden("Anda hanya dapat membatalkan KRS milik sendiri")
	}

	if err := h.DB.Delete(&e).Error; err != nil {
		return err
	}
	return helper.NoContent(c)
}