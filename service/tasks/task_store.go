package tasks

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StoreTask struct {
	db *sqlx.DB
}

// make the new handler store for tasks
func NewTaskStore(s *sqlx.DB) *StoreTask {
	return &StoreTask{db: s}
}

// func to make the new tasks for a student
func (s *StoreTask) CreateNewTasks(ctx context.Context, task *types.Task) error {

	//setup the options for a transaction
	option_tx := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	//begin the transaction for this method
	tx, err := s.db.BeginTxx(ctx, option_tx)
	if err != nil {
		return errors.New("Failed to setup the transaction for this method!")
	}
	defer tx.Rollback()

	//query for find a student_id
	// var students types.Student
	// query_student_id := `
	// 	SELECT name, kelas, jurusan, absen, wali_kelas FROM students
	// 	WHERE id = $1;
	// `

	//execute the method query to find the student id
	// if err := s.db.Get(students.Id, query_student_id, id_student); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return errors.New("Failed to find the student id!")
	// 	}
	// 	return errors.New("Failed to get the student_id from student table!")
	// }

	//query for create a new task
	query_task := `
		INSERT INTO tasks (id, name_task, file_task, date_task, student_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING*;
	`

	//execute the method query to create a new data for task table
	if err := tx.QueryRowContext(
		ctx,
		query_task,
		task.Id,
		task.Name_Task,
		task.File_Task,
		task.Date_Task,
		task.Student_Id,
		task.Created_at,
		task.Updated_at,
	).Scan(
		&task.Id,
		&task.Name_Task,
		&task.File_Task,
		&task.Date_Task,
		&task.Student_Id,
		&task.Created_at,
		&task.Updated_at,
	); err != nil {
		return errors.New("Failed to scan the payload to real data in db!" + err.Error())
	}

	//commit the transaction if transaction has been successfully!
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the data create!")
	}

	return nil

}

// func to get task by id
func (s *StoreTask) GetTaskById(id uuid.UUID, ctx context.Context) (*types.Task, error) {

	//base query for this method to handle get task by id
	query := `
		SELECT name_task, file_task, date_task, student_id, created_at, updated_at 
		FROM tasks WHERE id = $1;
	`

	//execute the func for this method
	var tasks types.Task
	if err := s.db.GetContext(ctx, &tasks, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to see the rows in db!")
		}
		return nil, errors.New("Failed to get the task by id!" + err.Error())
	}

	//return final result
	return &tasks, nil

}
