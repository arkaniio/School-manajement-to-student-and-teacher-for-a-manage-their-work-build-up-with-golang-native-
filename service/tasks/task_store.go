package tasks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

// func to delete the data tasks (only students role or role siswa can access this method)
func (s *StoreTask) DeleteTask(id uuid.UUID, ctx context.Context) error {

	//setup the options transaction for this method
	tx_options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	//setup the transaction with sqlx method
	tx, err := s.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions for this method!")
	}
	defer tx.Rollback()

	//base query for this method
	query := `
		DELETE FROM tasks WHERE id = $1;
	`

	//execute the query
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return errors.New("Failed to execute the query for this method!")
	}

	//scan rows
	rows, err := result.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Failed to detect the rows in db!")
		}
		return errors.New("Failed to check the rows in db")
	}
	if rows == 0 {
		return errors.New("Failed to get the rows in db!")
	}

	//commit the transactions
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	//return final result
	return nil

}

// func to handle a update task for this method
func (s *StoreTask) UpdateTask(id uuid.UUID, ctx context.Context, payloads types.PayloadUpdate) error {

	//setup the options transaction for this method
	tx_options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	//setup the transaction with sqlx method
	tx, err := s.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transactions for this method!")
	}
	defer tx.Rollback()

	//make the variable to put the every query on it
	var settings []string
	argsID := 1
	var args []interface{}

	//if students wants to update their name task
	if payloads.Name_Task != nil {
		settings = append(settings, fmt.Sprintf("name_task=$%d", argsID))
		argsID++
		args = append(args, *payloads.Name_Task)
	}

	//if students wants to update their file_task
	if payloads.File_Task != nil {
		settings = append(settings, fmt.Sprintf("file_task=$%d", argsID))
		argsID++
		args = append(args, *payloads.File_Task)
	}

	//update the updated at
	settings = append(settings, fmt.Sprintf("updated_at=$%d", argsID))
	argsID++
	args = append(args, time.Now().UTC())

	//combine the query
	full_query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = %d", strings.Join(settings, ","), argsID)
	args = append(args, argsID)

	//execute query for this method
	result, err := tx.ExecContext(ctx, full_query, args...)
	if err != nil {
		return errors.New("Failed to execute the query for this method!")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("No rows detected!")
		}
		return errors.New("Failed to get the rows!")
	}
	if rows == 0 {
		return errors.New("Failed to get the rows from db!")
	}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions")
	}

	//return final result
	return nil

}

// func to handle a get task by id!
func (s *StoreTask) GetTaskByIdIncludeStudents(id uuid.UUID, ctx context.Context) (*types.TaskIncludeStudents, error) {

	//base query for get the task by id
	query := `
		SELECT * FROM tasks t
		INNER JOIN students s ON t.student_id = s.id
		WHERE t.id = $1;
	`

	//execute the query
	var tasks types.TaskIncludeStudents
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, errors.New("Failed to get the tasks data!")
	}
	for rows.Next() {
		var students types.Student
		if err := rows.Scan(tasks); err != nil {
			return nil, errors.New("Failed to scan the data in db!")
		}
		tasks.Students = append(tasks.Students, students)
	}

	//return final result
	return &tasks, nil

}
