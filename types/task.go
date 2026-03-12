package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TaskStore interface {
	CreateNewTasks(ctx context.Context, task *Task) error
	GetTaskById(id uuid.UUID, ctx context.Context) (*Task, error)
}

type Task struct {
	Id         uuid.UUID `db:"id"`
	Name_Task  string    `db:"name_task"`
	File_Task  string    `db:"file_task"`
	Date_Task  time.Time `db:"date_task"`
	Student_Id uuid.UUID `db:"student_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}

type Payload struct {
	Id         uuid.UUID `json:"id"`
	Name_Task  string    `json:"name_task" validate:"required"`
	File_Task  string    `json:"file_task" validate:"required"`
	Date_Task  time.Time `json:"date_task" validate:"required"`
	Student_Id uuid.UUID `json:"student_id" validate:"required"`
	Created_at time.Time `json:"created_at" validate:"required"`
	Updated_at time.Time `json:"updated_at" validate:"required"`
}

type ResponseTask struct {
	Id         uuid.UUID `json:"id"`
	Name_Task  string    `json:"name_task"`
	File_Task  string    `json:"file_task"`
	Date_Task  string    `json:"date_task"`
	Student_Id uuid.UUID `json:"student_id"`
	Created_at string    `json:"created_at"`
	Updated_at string    `json:"updated_at"`
}
