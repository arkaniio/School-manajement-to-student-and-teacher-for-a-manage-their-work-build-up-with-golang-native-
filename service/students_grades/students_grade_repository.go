package studentsgrades

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
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
		INSERT INTO students_grade (id, task_id, tanggal, keterangan, grades, created_at, updated_at)
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
