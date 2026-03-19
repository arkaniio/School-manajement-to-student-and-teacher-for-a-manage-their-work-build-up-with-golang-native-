package studentsgrades

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// make the type for repo
type StoreStudentsGrade struct {
	db *sqlx.DB
}

// make the func for repository
func NewHandlerStoreStudentsGrade(db *sqlx.DB) *StoreStudentsGrade {
	return &StoreStudentsGrade{db: db}
}

// add func to create a new handler students grade
func (s *StoreStudentsGrade) CreateNewStudentsGrade(ctx context.Context, payloads *types.StudentsGrade) error {

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

	//execute and settings the query for this method
	query := `
		INSERT INTO students_grades (id, task_id, tanggal, keterangan, grades, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING*;
	`

	if err := tx.QueryRowContext(
		ctx,
		query,
		payloads.Id,
		payloads.Task_Id,
		payloads.Tanggal,
		payloads.Keterangan,
		payloads.Grades,
		payloads.Created_at,
		payloads.Updated_at,
	).Scan(
		&payloads.Id,
		&payloads.Task_Id,
		&payloads.Tanggal,
		&payloads.Keterangan,
		&payloads.Grades,
		&payloads.Created_at,
		&payloads.Updated_at,
	); err != nil {
		return errors.New("Failed to get settings and setup the query for this method!")
	}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction!")
	}

	//return final result
	return nil

}

func (s *StoreStudentsGrade) UpdateStudentsGrade(ctx context.Context, id uuid.UUID, payloads *types.PayloadsStudentGradeUpdate) error {
	option_tx := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := s.db.BeginTxx(ctx, option_tx)
	if err != nil {
		return errors.New("Failed to setup the transaction for update!")
	}
	defer tx.Rollback()

	// Build dynamic update query
	updates := []string{"updated_at = $2"}
	args := []interface{}{payloads.Updated_at}
	argIndex := 3

	if payloads.Task_Id != nil {
		updates = append(updates, fmt.Sprintf("task_id = $%d", argIndex))
		args = append(args, *payloads.Task_Id)
		argIndex++
	}
	if payloads.Tanggal != nil {
		updates = append(updates, fmt.Sprintf("tanggal = $%d", argIndex))
		args = append(args, *payloads.Tanggal)
		argIndex++
	}
	if payloads.Keterangan != nil {
		updates = append(updates, fmt.Sprintf("keterangan = $%d", argIndex))
		args = append(args, *payloads.Keterangan)
		argIndex++
	}
	if payloads.Grades != nil {
		updates = append(updates, fmt.Sprintf("grades = $%d", argIndex))
		args = append(args, *payloads.Grades)
		argIndex++
	}

	args = append([]interface{}{id}, args...)

	query := fmt.Sprintf(`
		UPDATE students_grades 
		SET %s 
		WHERE id = $1
	`, strings.Join(updates, ", "))

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return errors.New("Failed to update students grade!")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("No students grade found with given ID!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit update transaction!")
	}

	return nil
}

func (s *StoreStudentsGrade) DeleteStudentsGrade(ctx context.Context, id uuid.UUID) error {
	option_tx := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	tx, err := s.db.BeginTxx(ctx, option_tx)
	if err != nil {
		return errors.New("Failed to setup delete transaction!")
	}
	defer tx.Rollback()

	query := `DELETE FROM students_grades WHERE id = $1`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return errors.New("Failed to delete students grade!")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("No students grade found with given ID!")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit delete transaction!")
	}

	return nil
}

func (s *StoreStudentsGrade) GetAllStudentsGradesWithTask(ctx context.Context) ([]types.StudentsGradesWithTaskResponse, error) {
	query := `
		SELECT sg.id, sg.task_id, sg.tanggal, sg.keterangan, sg.grades, sg.created_at, sg.updated_at,
		       t.name_task, t.mapel_task
		FROM students_grades sg
		LEFT JOIN tasks t ON sg.task_id = t.id
		ORDER BY sg.created_at DESC
	`

	var grades []types.StudentsGradesWithTaskResponse
	err := s.db.SelectContext(ctx, &grades, query)
	if err != nil {
		return nil, errors.New("Failed to get all students grades with tasks!")
	}

	return grades, nil
}
