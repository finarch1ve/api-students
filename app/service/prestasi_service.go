package service

import (
	"context"

	"api-students/app/model"
	"api-students/app/repository"
)

type PrestasiService struct {
	studentRepo  repository.StudentRepository
	prestasiRepo repository.PrestasiRepository
}

func NewPrestasiService(studentRepo repository.StudentRepository, prestasiRepo repository.PrestasiRepository) *PrestasiService {
	return &PrestasiService{studentRepo: studentRepo, prestasiRepo: prestasiRepo}
}

func (s *PrestasiService) GetByNIM(ctx context.Context, nim string) (model.Student, []model.Prestasi, error) {
	student, err := s.studentRepo.FindByNIM(ctx, nim)
	if err != nil {
		return model.Student{}, nil, err
	}

	prestasi, err := s.prestasiRepo.FindByStudentID(ctx, student.ID)
	if err != nil {
		return model.Student{}, nil, err
	}

	return student, prestasi, nil
}