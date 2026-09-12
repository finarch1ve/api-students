package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/app/service"
	"api-students/helper"
)

type StudentHandler struct {
	svc *service.StudentService
}

func NewStudentHandler(svc *service.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func terjemahkanError(c *fiber.Ctx, err error, pesanUmum string) error {
	var valErr *service.ValidationError
	switch {
	case errors.As(err, &valErr):
		return failValidation(c, valErr.Errors)
	case errors.Is(err, service.ErrNoFieldsToUpdate):
		return fail(c, fiber.StatusBadRequest, err.Error())
	case errors.Is(err, repository.ErrNotFound):
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "NIM sudah terdaftar")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}

func (h *StudentHandler) List(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	q := parseListQuery(c)

	hasil, meta, err := h.svc.List(ctx, q)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "gagal mengambil data mahasiswa")
	}
	return okList(c, "daftar mahasiswa berhasil diambil", hasil, &meta)
}

func (h *StudentHandler) Get(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	student, err := h.svc.Get(ctx, id)
	if err != nil {
		return terjemahkanError(c, err, "gagal mengambil data mahasiswa")
	}
	return ok(c, "mahasiswa ditemukan", student)
}

func (h *StudentHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	baru, err := h.svc.Create(ctx, req)
	if err != nil {
		return terjemahkanError(c, err, "gagal menyimpan mahasiswa")
	}

	return created(c, "mahasiswa berhasil dibuat", baru, "/api/v1/students/"+strconv.Itoa(baru.ID))
}

func (h *StudentHandler) Replace(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	hasil, err := h.svc.Replace(ctx, id, req)
	if err != nil {
		return terjemahkanError(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diganti seluruhnya", hasil)
}

func (h *StudentHandler) Patch(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	var req model.PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	hasil, err := h.svc.Patch(ctx, id, req)
	if err != nil {
		return terjemahkanError(c, err, "gagal memperbarui mahasiswa")
	}

	return ok(c, "mahasiswa berhasil diperbarui sebagian", hasil)
}

func (h *StudentHandler) Delete(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	if err := h.svc.Delete(ctx, id); err != nil {
		return terjemahkanError(c, err, "gagal menghapus mahasiswa")
	}
	return noContent(c)
}

// ==================== PRESTASI ====================

type PrestasiHandler struct {
	svc *service.PrestasiService
}

func NewPrestasiHandler(svc *service.PrestasiService) *PrestasiHandler {
	return &PrestasiHandler{svc: svc}
}

func (h *PrestasiHandler) GetByNIM(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	nim := c.Params("nim")
	if strings.TrimSpace(nim) == "" {
		return fail(c, fiber.StatusBadRequest, "nim wajib diisi")
	}

	student, prestasi, err := h.svc.GetByNIM(ctx, nim)
	if err != nil {
		return terjemahkanError(c, err, "gagal mengambil data prestasi")
	}

	return ok(c, "data prestasi berhasil diambil", fiber.Map{
		"student":  student,
		"prestasi": prestasi,
	})
}

// ==================== AUTH ====================

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	user, err := h.svc.Register(ctx, req)
	if err != nil {
		return terjemahkanAuthError(c, err, "gagal mendaftarkan user")
	}

	return created(c, "pendaftaran berhasil", user, "/api/v1/auth/me")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	pair, err := h.svc.Login(ctx, req)
	if err != nil {
		return terjemahkanAuthError(c, err, "gagal login")
	}

	return ok(c, "login berhasil", pair)
}

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	pair, err := h.svc.Refresh(ctx, req)
	if err != nil {
		return terjemahkanAuthError(c, err, "gagal memperbarui token")
	}

	return ok(c, "token berhasil diperbarui", pair)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	var req model.RefreshRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	_ = h.svc.Logout(ctx, req)
	return ok(c, "logout berhasil", nil)
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	ctx, cancel := reqCtx(c)
	defer cancel()

	authUser, ok2 := helper.CurrentUser(c)
	if !ok2 {
		return fail(c, fiber.StatusUnauthorized, "belum terautentikasi")
	}

	user, err := h.svc.Me(ctx, authUser.UserID)
	if err != nil {
		return fail(c, fiber.StatusUnauthorized, "user tidak ditemukan")
	}

	return ok(c, "profil berhasil diambil", user)
}

func terjemahkanAuthError(c *fiber.Ctx, err error, pesanUmum string) error {
	var valErr *service.ValidationError
	switch {
	case errors.As(err, &valErr):
		return failValidation(c, valErr.Errors)
	case errors.Is(err, repository.ErrDuplicate):
		return fail(c, fiber.StatusConflict, "username sudah dipakai")
	case errors.Is(err, service.ErrInvalidCredentials):
		return fail(c, fiber.StatusUnauthorized, "username atau password salah")
	case errors.Is(err, service.ErrAccountDisabled):
		return fail(c, fiber.StatusForbidden, "akun dinonaktifkan")
	default:
		return fail(c, fiber.StatusInternalServerError, pesanUmum)
	}
}