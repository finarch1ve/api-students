package service

import (
	"context"
	"errors"

	"api-students/app/model"
	"api-students/app/repository"
)

// ErrNoFieldsToUpdate dipakai saat request PATCH tidak mengubah field apa pun.
var ErrNoFieldsToUpdate = errors.New("tidak ada field yang diubah")

// ValidationError merepresentasikan kumpulan error validasi per field.
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
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) *StudentService {
	return &StudentService{repo: repo}
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

func (s *StudentService) Get(ctx context.Context, id int) (model.Student, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *StudentService) Create(ctx context.Context, req model.CreateStudentRequest) (model.Student, error) {
	if errs := ValidateCreate(req); len(errs) > 0 {
		return model.Student{}, newValidationError(errs)
	}

	return s.repo.Create(ctx, model.Student{
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
	})
}

func (s *StudentService) Replace(ctx context.Context, id int, req model.ReplaceStudentRequest) (model.Student, error) {
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

func (s *StudentService) Patch(ctx context.Context, id int, req model.PatchStudentRequest) (model.Student, error) {
	if IsEmptyPatch(req) {
		return model.Student{}, ErrNoFieldsToUpdate
	}

	saatIni, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return model.Student{}, err
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