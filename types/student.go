package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StudentStore interface {
	CreateNewStudent(ctx context.Context, students *Student) error
	DeleteStudents(id uuid.UUID, ctx context.Context) error
	UpdateStudentsData(id uuid.UUID, payload UpdateAsStudent, ctx context.Context) error
	GetStudentById(id uuid.UUID, ctx context.Context) (*Student, error)
	GetAllStudents(ctx context.Context) ([]Student, error)
}

type Student struct {
	Id             uuid.UUID `db:"id"`
	Full_Name      string    `db:"full_name"`
	Kelas          string    `db:"kelas"`
	Jurusan        string    `db:"jurusan"`
	Absen          int       `db:"absen"`
	StudentProfile string    `db:"student_profile"`
	Wali_Kelas     string    `db:"wali_kelas"`
	Created_at     time.Time `db:"created_at"`
	Updated_at     time.Time `db:"updated_at"`
	MapelStudents  string    `db:"mapel_students"`
}

type RegisterAsStudent struct {
	Id             uuid.UUID `json:"id"`
	Full_name      string    `json:"full_name" validate:"required"`
	Kelas          string    `json:"kelas" validate:"required"`
	Jurusan        string    `json:"jurusan" validate:"required"`
	Absen          int       `json:"absen" validate:"required"`
	StudentProfile string    `json:"student_profile"`
	Wali_Kelas     string    `json:"wali_kelas" validate:"required"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	MapelStudents  string    `json:"mapel_students" validate:"required"`
}

type UpdateAsStudent struct {
	Id             uuid.UUID `json:"id"`
	Full_name      *string   `json:"full_name"`
	Kelas          *string   `json:"kelas"`
	Jurusan        *string   `json:"jurusan"`
	Absen          *int      `json:"absen"`
	StudentProfile *string   `json:"student_profile"`
	Wali_Kelas     *string   `json:"wali_kelas"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
	MapelStudents  *string   `json:"mapel_students"`
}

type StudentResponse struct {
	Id             uuid.UUID `json:"id"`
	Full_name      string    `json:"full_name"`
	Kelas          string    `json:"kelas"`
	Jurusan        string    `json:"jurusan"`
	Absen          int       `json:"absen"`
	StudentProfile string    `json:"student_profile"`
	Wali_Kelas     string    `json:"wali_kelas"`
	Created_at     string    `json:"created_at"`
	Updated_at     string    `json:"updated_at"`
	MapelStudents  string    `json:"mapel_students"`
}
