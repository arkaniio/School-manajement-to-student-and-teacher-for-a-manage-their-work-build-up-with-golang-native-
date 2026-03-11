package students

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
		INSERT INTO students (id, full_name, kelas, jurusan, absen, student_profile, wali_kelas, created_at, updated_at)
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

// func to get the data of students by their id to detect their data in table
func (s *StudentStore) GetStudentById(id uuid.UUID, ctx context.Context) (*types.Student, error) {

	//base query for this method
	query := `
		SELECT full_name, kelas, jurusan, absen, student_profile, wali_kelas, created_at, updated_at FROM students
		WHERE id = $1;
	`

	//execute the query
	var students types.Student
	if err := s.db.GetContext(ctx, &students, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get the students by id! result is nil" + err.Error())
		}
		return nil, errors.New("Failed to get the students by id!" + err.Error())
	}

	//return final result
	return &students, nil

}

// func and method to update the students data
func (s *StudentStore) UpdateStudentsData(id uuid.UUID, payload types.UpdateAsStudent, ctx context.Context) error {

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

	//setup the base query and args
	var settings []string
	argsID := 1
	var args []interface{}

	//if students wants to update their full_name
	if payload.Full_name != nil {
		settings = append(settings, fmt.Sprintf("full_name=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//if students wants to update their kelas
	if payload.Kelas != nil {
		settings = append(settings, fmt.Sprintf("kelas=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//if students wants to update their absen
	if payload.Absen != nil {
		settings = append(settings, fmt.Sprintf("absen=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//if students wants to update their jurusan
	if payload.Jurusan != nil {
		settings = append(settings, fmt.Sprintf("jurusan=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//if students wants to update their profile
	if payload.StudentProfile != nil {
		settings = append(settings, fmt.Sprintf("student_profile=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//if students wants to update their wali_kelas
	if payload.Wali_Kelas != nil {
		settings = append(settings, fmt.Sprintf("wali_kelas=$%d", argsID))
		argsID++
		args = append(args, settings)
	}

	//update the updated at
	settings = append(settings, fmt.Sprintf("updated_at=$%d", argsID))
	argsID++
	args = append(args, settings)

	//combine all of interface to one interface
	base_query := fmt.Sprintf("UPDATE students SET %s WHERE id = $%d", strings.Join(settings, ","), argsID)
	settings = append(settings, base_query)
	args = append(args, settings)

	//execute the query context
	result, err := tx.ExecContext(ctx, base_query, args...)
	if err != nil {
		return errors.New("Failed to update the data students")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to detect the rows affected based on db")
	}
	if rows == 0 {
		return errors.New("Invalid rows!")
	}

	//commit the transactions
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction")
	}

	//return final result based on returning in this method or func
	return nil

}
