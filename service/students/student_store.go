package students

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// type for a store student
type StudentStore struct {
	db *sqlx.DB
}

// func that we use when we want to use the store from this db
func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{db: db}
}

// make the func create a new student data
func (s *StudentStore) CreateNewStudent(ctx context.Context, students *types.Student) error {

	//settings the transaction options
	tx_options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	//settings the transaction method for db
	tx, err := s.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to settings the transaction method for this func create!" + err.Error())
	}
	defer tx.Rollback()

	//query for this func create a new student
	query := `
		INSERT INTO students (id, full_name, kelas, jurusan, absen, wali_kelas, student_profile, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING*;
	`

	//execute the query
	if err := tx.QueryRowContext(
		ctx,
		query,
		students.Id,
		students.Full_Name,
		students.Kelas,
		students.Jurusan,
		students.Absen,
		students.StudentProfile,
		students.Wali_Kelas,
		students.Created_at,
		students.Updated_at,
	).Scan(
		&students.Id,
		&students.Full_Name,
		&students.Kelas,
		&students.Jurusan,
		&students.Absen,
		&students.StudentProfile,
		&students.Wali_Kelas,
		&students.Created_at,
		&students.Updated_at,
	); err != nil {
		return errors.New("Failed to scan the students struct!" + err.Error())
	}

	//commit the transaction to final result
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction" + err.Error())
	}

	return nil

}

// func that containt the delete query at here
func (s *StudentStore) DeleteStudents(id uuid.UUID, ctx context.Context) error {

	//setup the transactions option here
	options := &sql.TxOptions{
		Isolation: sql.LevelSerializable,
		ReadOnly:  false,
	}

	//setup the final transaction
	tx, err := s.db.BeginTxx(ctx, options)
	if err != nil {
		return errors.New("Failed to setup the final transaction for this method")
	}
	defer tx.Rollback()

	//base query for this methid
	query := `
		DELETE FROM students WHERE id = $1;
	`

	//execute the query for this method
	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return errors.New("Failed to get the result of every rows in students table!")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to detected the rows in table students")
	}
	if rows == 0 {
		return errors.New("Invali rows!")
	}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transactions!")
	}

	//return final result
	return nil

}
