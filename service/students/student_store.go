package students

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
)

// type for a store student
type StudentStore struct {
	db *sqlx.DB
}

// func that we use when we want to use the store from this db
func NewStudentStore(db *sqlx.DB) *StudentStore {
	return &StudentStore{db: db}
}

//helper for a pagination student

type SortedConfig struct {
	Column   string
	Operator string
	Order    string
}

type Cursor struct {
	Value any
	Id    string
}

func getSorted(sorted string, order string) SortedConfig {

	//make the interface
	sorted_descasc := map[string]string{
		"created_at": "created_at",
		"name":       "name",
	}

	col, ok := sorted_descasc[sorted]
	if !ok {
		col = "created_at"
	}

	operator := "<"
	description := "DESC"

	if order == strings.ToLower("asc") {
		operator = ">"
		description = "ASC"
	}

	return SortedConfig{
		Column:   col,
		Operator: operator,
		Order:    description,
	}

}

// func that create a new student
func (s *StudentStore) CreateNewStudent(ctx context.Context, student *types.Student) error {

	//make the options of transaction of create
	options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	//make the ctx for context in transaction
	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
	defer cancle()

	//make the new begintxx for transaction
	tx, err := s.db.BeginTxx(ctx, options)
	if err != nil {
		return errors.New("Failed to settings the db transactions")
	}
	defer tx.Rollback()

	//make the base query for create a new student
	query := `
		INSERT INTO students 
		(id, name, class, address, major, student_profile, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING *;
	`

	//make the method query
	if err := tx.QueryRowContext(
		ctx,
		query,
		student.Id,
		student.Name,
		student.Class,
		student.Address,
		student.Major,
		student.StudentProfile,
		student.Created_at,
		student.Updated_at,
	).Scan(
		&student.Id,
		&student.Name,
		&student.Class,
		&student.Address,
		&student.Major,
		&student.StudentProfile,
		&student.Created_at,
		&student.Updated_at,
	); err != nil {
		return errors.New("Failed to create a new user!" + err.Error())
	}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the query of transaction!" + err.Error())
	}

	//final return
	return nil

}

// func to get the name in student db
func (s *StudentStore) GetStudentByName(name string) (*types.Student, error) {

	//make base query
	query := `
		SELECT id, name, class, major, student_profile, created_at, updated_at 
		FROM students WHERE name = $1;
	`

	//make the method of query
	var students types.Student
	if err := s.db.Get(&students, query, name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &students, nil

}

// get all student from db
func (s *StudentStore) GetAllStudents(
	ctx context.Context,
	limit int,
	sort string,
	order string,
	cursorValue any,
	cursorID string,
) ([]types.Student, error) {

	sort_config := getSorted(sort, order)

	//base query
	query := fmt.Sprintf(
		`
			SELECT id, name, class, address, major, student_profile, created_at, updated_at 
			FROM users WHERE ($1 IS NULL OR (%s, id) %s ($1, $2))
			ORDER BY %s, %s, id %s LIMIT $3;
		`, sort_config.Column, sort_config.Operator, sort_config.Column, sort_config.Order, sort_config.Order,
	)

	//execute the query
	var students []types.Student
	if err := s.db.SelectContext(ctx, &students, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed to get the students data")
		}
		return nil, nil
	}

	//return final result
	return students, nil

}

// func to update their student
func (s *StudentStore) UpdateStudent(
	id uuid.UUID,
	ctx context.Context,
	payload types.UpdateAsStudent,
) error {

	//make the options for a options of transaction
	options := &sql.TxOptions{
		ReadOnly:  false,
		Isolation: sql.LevelSerializable,
	}

	//setup the transaction
	tx, err := s.db.BeginTxx(ctx, options)
	if err != nil {
		return errors.New("Failed to get the transactrion!" + err.Error())
	}

	//setup to join the query for a settings and dll
	var settings_query []string
	var args []interface{}
	var argID = 1

	//if the students wants to update their name
	if payload.Name != nil {
		settings_query = append(settings_query, fmt.Sprintf("name=$%d", argID))
		args = append(args, payload.Name)
		argID++
	}

	//if the students wants to update their class
	if payload.Class != nil {
		settings_query = append(settings_query, fmt.Sprintf("class=$%d", argID))
		args = append(args, payload.Class)
		argID++
	}

	//if the students wants to update their address
	if payload.Address != nil {
		settings_query = append(settings_query, fmt.Sprintf("address=$%d", argID))
		args = append(args, payload.Address)
		argID++
	}

	//if the students wants to update their student_profile
	if payload.StudentProfile != nil {
		settings_query = append(settings_query, fmt.Sprintf("student_profile=$%d", argID))
		args = append(args, payload.StudentProfile)
		argID++
	}

	//if students wants to update their major
	if payload.Major != nil {
		settings_query = append(settings_query, fmt.Sprintf("major=$%d", argID))
		args = append(args, payload.Major)
		argID++
	}

	//update the updated at
	settings_query = append(settings_query, fmt.Sprintf("updated_at=$%d", argID))
	args = append(args, payload.Updated_at)
	argID++

	//set the full query
	full_query := fmt.Sprintf("UPDATE students SET %s WHERE id = $%d", strings.Join(settings_query, ", "), argID)
	args = append(args, full_query)

	//execute the query
	if result, err := tx.ExecContext(ctx, full_query, args...); err != nil {
		rows, err := result.RowsAffected()
		if err != nil {
			return errors.New("Failed to detect the rows in the db!" + err.Error())
		}
		if rows == 0 {
			return errors.New("No one rows detected !")
		}
		return errors.New("Failed to execute the query")
	}

	//return the error
	return nil

}

// func to get the student by id
func (s *StudentStore) GetStudentById(id uuid.UUID, ctx context.Context) (*types.Student, error) {

	//base query
	query := `
		SELECT id, name, class, address, major, student_profile, created_at, updated_at
		WHERE id = $1
	`

	//execute the query
	var students types.Student
	if err := s.db.GetContext(ctx, students, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Failed error no rows detected")
		}
		return nil, nil
	}

	//return final result
	return &students, nil

}

// func to delete the student data
func (s *StudentStore) DeleteStudentsData(id uuid.UUID, ctx context.Context) error {

	//setup the transaction
	tx_options := &sql.TxOptions{
		Isolation: sql.LevelLinearizable,
		ReadOnly:  false,
	}

	//begin tx / transaction
	tx, err := s.db.BeginTxx(ctx, tx_options)
	if err != nil {
		return errors.New("Failed to setup the transaction!")
	}
	defer tx.Rollback()

	//setup the query
	query := `
	DELETE FROM students WHERE id = $1;
	`

	//execute the query for deleting students data
	if result, err := tx.ExecContext(ctx, query, id); err != nil {
		final_result, err := result.RowsAffected()
		if err != nil {
			return errors.New("Failed to see the data of result in db!")
		}
		if final_result == 0 {
			return errors.New("No macthed with the result in db!")
		}
		return errors.New("Failed to execute the query with a context!")
	}

	return nil

}
