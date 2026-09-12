package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"api-students/app/model"
)

type PrestasiRepository interface {
	FindByStudentID(ctx context.Context, studentID int) ([]model.Prestasi, error)
}

type prestasiRepository struct {
	pool *pgxpool.Pool
}

func NewPrestasiRepository(pool *pgxpool.Pool) PrestasiRepository {
	return &prestasiRepository{pool: pool}
}

func (r *prestasiRepository) FindByStudentID(ctx context.Context, studentID int) ([]model.Prestasi, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, nama_prestasi, id_student, juara FROM prestasi WHERE id_student = $1`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hasil := []model.Prestasi{}
	for rows.Next() {
		var p model.Prestasi
		if err := rows.Scan(&p.ID, &p.NamaPrestasi, &p.StudentID, &p.Juara); err != nil {
			return nil, err
		}
		hasil = append(hasil, p)
	}
	return hasil, rows.Err()
}