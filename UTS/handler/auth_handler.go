package handler

import (
	"errors"
	"strings"

	"siakad-mini/app/model"
	"siakad-mini/helper"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.NewAppError(fiber.StatusBadRequest, "Body request tidak valid")
	}
	if errs := helper.ValidateStruct(req); errs != nil {
		return helper.ValidationFailed(errs)
	}

	var user model.User
	err := h.DB.Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(req.Email))).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.Unauthorized("Email atau password salah")
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return helper.Unauthorized("Email atau password salah")
	}

	token, err := helper.GenerateToken(user.ID, user.Role)
	if err != nil {
		return err
	}

	return helper.Success(c, fiber.StatusOK, "Login berhasil", fiber.Map{
		"token":      token,
		"token_type": "Bearer",
		"expires_in": int(helper.TokenTTL().Seconds()),
		"user":       user,
	})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return helper.Unauthorized("Sesi tidak valid")
	}

	var user model.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.Unauthorized("Pengguna tidak ditemukan")
		}
		return err
	}

	data := fiber.Map{"user": user}
	if user.Role == "mahasiswa" {
		var st model.Student
		if err := h.DB.Where("user_id = ?", user.ID).First(&st).Error; err == nil {
			data["student"] = st
		}
	}
	return helper.Success(c, fiber.StatusOK, "Data pengguna", data)
}