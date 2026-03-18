package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type StudentsGradeStore interface {
	CreateNewStudentsGrade(ctx context.Context, payloads *StudentsGrade) error
}

type StudentsGrade struct {
	Id         uuid.UUID `db:"id"`
	Task_Id    uuid.UUID `db:"task_id"`
	Tanggal    string    `db:"tanggal"`
	Keterangan string    `db:"keterangan"`
	Grades     string    `db:"grades"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}

type PayloadsStudentGrade struct {
	Id         uuid.UUID `json:"id"`
	Task_Id    uuid.UUID `json:"task_id" validate:"required"`
	Tanggal    string    `json:"tanggal" validate:"required"`
	Keterangan string    `json:"keterangan" validate:"required"`
	Grades     string    `json:"grades" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
	Updated_at time.Time `json:"updated_at" validate:"required"`
}

type PayloadsStudentGradeUpdate struct {
	Id         uuid.UUID  `json:"id"`
	Task_Id    *uuid.UUID `json:"task_id"`
	Tanggal    *string    `json:"tanggal"`
	Keterangan *string    `json:"keterangan"`
	Grades     *string    `jso:"grades"`
	Created_at time.Time  `json:"created_at"`
	Updated_at time.Time  `json:"updated_at"`
}

type StudentsGradeResponse struct {
	Id         uuid.UUID `json:"id"`
	Task_Id    uuid.UUID `json:"task_id" validate:"required"`
	Tanggal    string    `json:"tanggal" validate:"required"`
	Keterangan string    `json:"keterangan" validate:"required"`
	Grades     string    `json:"grades" validate:"required"`
	Created_at string    `json:"created_at" validate:"required"`
	Updated_at string    `json:"updated_at" validate:"required"`
}
