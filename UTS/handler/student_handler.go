package handler

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"siakad-mini/app/model"
	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentHandler struct {
	DB *gorm.DB
}

type CreateStudentRequest struct {
	NIM         string   `json:"nim" validate:"required,len=12,numeric"`
	Nama        string   `json:"nama" validate:"required"`
	Email       string   `json:"email" validate:"required,email"`
	Prodi       string   `json:"prodi" validate:"required"`
	Angkatan    int      `json:"angkatan" validate:"required"`
	IPKTerakhir *float64 `json:"ipk_terakhir" validate:"omitempty,gte=0,lte=4"`
}

type UpdateStudentRequest struct {
	Nama        string   `json:"nama" validate:"required"`
	Prodi       string   `json:"prodi" validate:"required"`
	Angkatan    int      `json:"angkatan" validate:"required"`
	IPKTerakhir *float64 `json:"ipk_terakhir" validate:"omitempty,gte=0,lte=4"`
}

func parseStudentID(c *fiber.Ctx) (uint, error) {
	n, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || n == 0 {
		return 0, helper.NotFound("Mahasiswa tidak ditemukan")
	}
	return uint(n), nil
}

func checkAngkatan(errs map[string][]string, angkatan int) {
	if _, bad := errs["angkatan"]; bad {
		return
	}
	if angkatan < 1000 || angkatan > time.Now().Year() {
		errs["angkatan"] = []string{"harus 4 digit dan tidak lebih dari tahun berjalan"}
	}
}

// GET /students
func (h *StudentHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	perPage := c.QueryInt("per_page", 10)
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 50 {
		perPage = 50
	}

	q := h.DB.Model(&model.Student{})

	if v := strings.TrimSpace(c.Query("prodi")); v != "" {
		q = q.Where("LOWER(prodi) = ?", strings.ToLower(v))
	}
	if v := strings.TrimSpace(c.Query("angkatan")); v != "" {
		a, err := strconv.Atoi(v)
		if err != nil {
			return helper.ValidationFailed(map[string][]string{"angkatan": {"harus berupa angka"}})
		}
		q = q.Where("angkatan = ?", a)
	}
	if v := strings.TrimSpace(c.Query("search")); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		q = q.Where("(LOWER(nim) LIKE ? OR LOWER(nama) LIKE ?)", like, like)
	}

	order := "id ASC"
	switch c.Query("sort") {
	case "":
	case "nama":
		order = "nama ASC, id ASC"
	case "-ipk_terakhir":
		order = "ipk_terakhir DESC, id ASC"
	default:
		return helper.ValidationFailed(map[string][]string{"sort": {"hanya boleh nama atau -ipk_terakhir"}})
	}

	base := q.Session(&gorm.Session{})

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return err
	}

	var items []model.Student
	if err := base.Order(order).Offset((page - 1) * perPage).Limit(perPage).Find(&items).Error; err != nil {
		return err
	}

	return helper.SuccessWithMeta(c, "Data mahasiswa berhasil diambil", items, helper.NewMeta(page, perPage, total))
}

// POST /students
func (h *StudentHandler) Create(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Unprocessable("Body request tidak valid")
	}

	errs := helper.ValidateStruct(req)
	if errs == nil {
		errs = map[string][]string{}
	}
	checkAngkatan(errs, req.Angkatan)

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var n int64
	if err := h.DB.Unscoped().Model(&model.Student{}).Where("nim = ?", req.NIM).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		errs["nim"] = append(errs["nim"], "NIM sudah terdaftar")
	}
	if err := h.DB.Model(&model.User{}).Where("LOWER(email) = ?", email).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		errs["email"] = append(errs["email"], "Email sudah terdaftar")
	}
	if len(errs) > 0 {
		return helper.ValidationFailed(errs)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NIM), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	ipk := 0.0
	if req.IPKTerakhir != nil {
		ipk = *req.IPKTerakhir
	}

	var student model.Student
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		user := model.User{Email: email, Password: string(hash), Role: "mahasiswa"}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		student = model.Student{
			UserID:      user.ID,
			NIM:         req.NIM,
			Nama:        strings.TrimSpace(req.Nama),
			Prodi:       strings.TrimSpace(req.Prodi),
			Angkatan:    req.Angkatan,
			IPKTerakhir: ipk,
		}
		return tx.Create(&student).Error
	})
	if err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusCreated, "Mahasiswa berhasil ditambahkan", fiber.Map{
		"id":           student.ID,
		"user_id":      student.UserID,
		"email":        email,
		"nim":          student.NIM,
		"nama":         student.Nama,
		"prodi":        student.Prodi,
		"angkatan":     student.Angkatan,
		"ipk_terakhir": student.IPKTerakhir,
	})
}

// GET /students/:id
func (h *StudentHandler) Show(c *fiber.Ctx) error {
	id, err := parseStudentID(c)
	if err != nil {
		return err
	}

	role, _ := c.Locals("role").(string)
	if role == "mahasiswa" {
		userID, _ := c.Locals("user_id").(uint)
		var own model.Student
		if err := h.DB.Where("user_id = ?", userID).First(&own).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return helper.NotFound("Mahasiswa tidak ditemukan")
			}
			return err
		}
		if own.ID != id {
			return helper.Forbidden("Anda hanya dapat mengakses data milik sendiri")
		}
	} else if role != "admin" {
		return helper.Forbidden("Anda tidak memiliki akses ke resource ini")
	}

	ta := strings.TrimSpace(c.Query("tahun_akademik"))

	var st model.Student
	err = h.DB.
		Preload("Enrollments", func(db *gorm.DB) *gorm.DB {
			if ta != "" {
				db = db.Where("tahun_akademik = ?", ta)
			}
			return db.Order("id ASC")
		}).
		Preload("Enrollments.Course").
		First(&st, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NotFound("Mahasiswa tidak ditemukan")
		}
		return err
	}

	totalSKS := 0
	for _, e := range st.Enrollments {
		totalSKS += e.Course.SKS
	}
	enrollments := st.Enrollments
	if enrollments == nil {
		enrollments = []model.Enrollment{}
	}

	return helper.Success(c, fiber.StatusOK, "Detail mahasiswa berhasil diambil", fiber.Map{
		"id":           st.ID,
		"user_id":      st.UserID,
		"nim":          st.NIM,
		"nama":         st.Nama,
		"prodi":        st.Prodi,
		"angkatan":     st.Angkatan,
		"ipk_terakhir": st.IPKTerakhir,
		"enrollments":  enrollments,
		"total_sks":    totalSKS,
		"batas_sks":    helper.MaxSKS(st.IPKTerakhir),
	})
}

// PUT /students/:id
func (h *StudentHandler) Update(c *fiber.Ctx) error {
	id, err := parseStudentID(c)
	if err != nil {
		return err
	}

	var st model.Student
	if err := h.DB.First(&st, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NotFound("Mahasiswa tidak ditemukan")
		}
		return err
	}

	var req UpdateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.Unprocessable("Body request tidak valid")
	}
	errs := helper.ValidateStruct(req)
	if errs == nil {
		errs = map[string][]string{}
	}
	checkAngkatan(errs, req.Angkatan)
	if len(errs) > 0 {
		return helper.ValidationFailed(errs)
	}

	ipk := st.IPKTerakhir
	if req.IPKTerakhir != nil {
		ipk = *req.IPKTerakhir
	}

	err = h.DB.Model(&st).Updates(map[string]interface{}{
		"nama":         strings.TrimSpace(req.Nama),
		"prodi":        strings.TrimSpace(req.Prodi),
		"angkatan":     req.Angkatan,
		"ipk_terakhir": ipk,
	}).Error
	if err != nil {
		return err
	}
	if err := h.DB.First(&st, st.ID).Error; err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusOK, "Data mahasiswa berhasil diperbarui", st)
}

// DELETE /students/:id (soft delete)
func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	id, err := parseStudentID(c)
	if err != nil {
		return err
	}

	var st model.Student
	if err := h.DB.First(&st, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NotFound("Mahasiswa tidak ditemukan")
		}
		return err
	}
	if err := h.DB.Delete(&st).Error; err != nil {
		return err
	}
	return helper.NoContent(c)
}