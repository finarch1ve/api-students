package service

import (
	"context"
	"errors"

	"api-students/app/model"
	"api-students/app/repository"
	"api-students/helper"
)

var ErrNoFieldsToUpdate = errors.New("tidak ada field yang diubah")
var ErrForbidden = errors.New("tidak berhak mengakses data ini")

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validasi gagal"
}

func newValidationError(errs map[string]string) error {
	if len(errs) == 0 {
		return nil
	}
	return &ValidationError{Errors: errs}
}

type StudentService struct {
	repo  repository.StudentRepository
	perms *helper.PermissionSet
}

func NewStudentService(repo repository.StudentRepository, perms *helper.PermissionSet) *StudentService {
	return &StudentService{repo: repo, perms: perms}
}

func (s *StudentService) List(ctx context.Context, q model.ListQuery) ([]model.Student, model.Meta, error) {
	hasil, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return nil, model.Meta{}, err
	}
	return hasil, model.Meta{
		Page: q.Page, Limit: q.Limit, Total: total,
		TotalPages: CountTotalPages(total, q.Limit),
	}, nil
}

func (s *StudentService) Get(ctx context.Context, current model.AuthUser, id int) (model.Student, error) {
	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Student{}, err
	}
	if !CanAccessStudent(current, student.OwnerID, s.perms, "student:read:any") {
		return model.Student{}, ErrForbidden
	}
	return student, nil
}

func (s *StudentService) Create(ctx context.Context, ownerID int, req model.CreateStudentRequest) (model.Student, error) {
	if errs := ValidateCreate(req); len(errs) > 0 {
		return model.Student{}, newValidationError(errs)
	}

	return s.repo.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
		OwnerID:  ownerID, // selalu dari identitas pemanggil, TIDAK dari body request
	})
}

func (s *StudentService) Replace(ctx context.Context, current model.AuthUser, id int, req model.ReplaceStudentRequest) (model.Student, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Student{}, err
	}
	if !CanAccessStudent(current, existing.OwnerID, s.perms, "student:update:any") {
		return model.Student{}, ErrForbidden
	}

	if errs := ValidateReplace(req); len(errs) > 0 {
		return model.Student{}, newValidationError(errs)
	}

	return s.repo.Update(ctx, model.Student{
		ID:       id,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: req.IsActive,
	})
}

func (s *StudentService) Patch(ctx context.Context, current model.AuthUser, id int, req model.PatchStudentRequest) (model.Student, error) {
	if IsEmptyPatch(req) {
		return model.Student{}, ErrNoFieldsToUpdate
	}

	saatIni, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Student{}, err
	}
	if !CanAccessStudent(current, saatIni.OwnerID, s.perms, "student:update:any") {
		return model.Student{}, ErrForbidden
	}

	updated, errs := ApplyPatch(saatIni, req)
	if len(errs) > 0 {
		return model.Student{}, newValidationError(errs)
	}

	return s.repo.Update(ctx, updated)
}

func (s *StudentService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}