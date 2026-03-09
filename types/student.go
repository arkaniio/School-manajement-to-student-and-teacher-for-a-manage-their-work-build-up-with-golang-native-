package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StudentStore interface {
	CreateNewStudent(ctx context.Context, students *Student) error
}

type Student struct {
	Id             uuid.UUID `db:"id"`
	Full_Name      string    `db:"full_name"`
	Kelas          string    `db:"kelas"`
	Jurusan        string    `db:"jurusan"`
	Absen          int32     `db:"absen"`
	StudentProfile string    `db:"student_profile"`
	Wali_Kelas     string    `db:"wali_kelas"`
	Created_at     time.Time `db:"created_at"`
	Updated_at     time.Time `db:"updated_at"`
}

type RegisterAsStudent struct {
	Id             uuid.UUID `json:"id"`
	Full_name      string    `json:"full_name"`
	Kelas          string    `json:"kelas"`
	Jurusan        string    `json:"jurusan"`
	Absen          int32     `json:"absen"`
	StudentProfile string    `json:"student_profile"`
	Wali_Kelas     string    `json:"wali_kelas"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

type UpdateAsStudent struct {
	Id             uuid.UUID `json:"id"`
	Full_name      *string   `json:"full_name"`
	Kelas          *string   `json:"kelas"`
	Jurusan        *string   `json:"jurusan"`
	Absen          *int32    `json:"absen"`
	StudentProfile *string   `json:"student_profile"`
	Wali_Kelas     *string   `json:"wali_kelas"`
	Created_at     time.Time `json:"created_at"`
	Updated_at     time.Time `json:"updated_at"`
}

type StudentResponse struct {
	Id             uuid.UUID `json:"id"`
	Full_name      string    `json:"full_name"`
	Kelas          string    `json:"kelas"`
	Jurusan        string    `json:"jurusan"`
	Absen          int32     `json:"absen"`
	StudentProfile string    `json:"student_profile"`
	Wali_Kelas     string    `json:"wali_kelas"`
	Created_at     string    `json:"created_at"`
	Updated_at     string    `json:"updated_at"`
}
