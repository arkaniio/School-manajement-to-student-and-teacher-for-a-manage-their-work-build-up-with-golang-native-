package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStore interface {
	CreateNewTasks(ctx context.Context, task *Task) error
	GetTaskById(id uuid.UUID, ctx context.Context) (*Task, error)
	DeleteTask(id uuid.UUID, ctx context.Context) error
	UpdateTask(id uuid.UUID, ctx context.Context, payloads PayloadUpdate) error
	GetTaskByIdIncludeStudents(id uuid.UUID, ctx context.Context) ([]TaskWithStudents, error)
}

type Task struct {
	Id         uuid.UUID `db:"id"`
	Name_Task  string    `db:"name_task"`
	File_Task  string    `db:"file_task"`
	Date_Task  time.Time `db:"date_task"`
	Student_Id uuid.UUID `db:"student_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Mapel_Task string    `db:"mapel_task"`
}

type TaskIncludeStudents struct {
	Id             uuid.UUID `db:"id"`
	NameTask       string    `db:"name_task"`
	File_Task      string    `db:"file_task"`
	Date_Task      time.Time `db:"date_task"`
	Student_Id     uuid.UUID `db:"student_id"`
	MapelTask      string    `db:"mapel_task"`
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

type TaskWithStudents struct {
	Id         uuid.UUID `db:"id"`
	Name_Task  string    `db:"name_task"`
	File_Task  string    `db:"file_task"`
	Date_Task  time.Time `db:"date_task"`
	Student_Id uuid.UUID `db:"student_id"`
	Students   Student   `db:"students"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	MapelTask  string    `db:"mapel_task"`
}

type Payload struct {
	Id         uuid.UUID `json:"id"`
	Name_Task  string    `json:"name_task" validate:"required"`
	File_Task  string    `json:"file_task" validate:"required"`
	Date_Task  time.Time `json:"date_task" validate:"required"`
	Student_Id uuid.UUID `json:"student_id" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
	Updated_at time.Time `json:"updated_at" validate:"required"`
	MapelTask  string    `json:"mapel_task" validate:"required"`
}

type PayloadUpdate struct {
	Id         uuid.UUID  `json:"id"`
	Name_Task  *string    `json:"name_task"`
	File_Task  *string    `json:"file_task"`
	Date_Task  time.Time  `json:"date_task"`
	Student_Id *uuid.UUID `json:"student_id"`
	Created_at time.Time  `json:"created_at"`
	Updated_at time.Time  `json:"updated_at"`
	MapelTask  *string    `json:"mapel_task"`
}

type ResponseTask struct {
	Id         uuid.UUID `json:"id"`
	Name_Task  string    `json:"name_task"`
	File_Task  string    `json:"file_task"`
	Date_Task  string    `json:"date_task"`
	Student_Id uuid.UUID `json:"student_id"`
	Created_at string    `json:"created_at"`
	Updated_at string    `json:"updated_at"`
	MapelTask  string    `json:"mapel_task"`
}
