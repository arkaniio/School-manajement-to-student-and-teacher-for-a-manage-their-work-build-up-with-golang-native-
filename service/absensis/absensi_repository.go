package absensis

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
type StoreAbsensi struct {
	db *sqlx.DB
}

// make the func for repository
func NewHandlerStoreAbsensi(db *sqlx.DB) *StoreAbsensi {
	return &StoreAbsensi{db: db}
}

// func to make the new absensi
func (s *StoreAbsensi) CreateNewAbsensi(ctx context.Context, payloads *types.Absensi) error {

	//make the transactions method
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

	//base query
	query := `
		INSERT INTO absensis (id, name_lengkap, kelas, jurusan, hari, tanggal, status, keterangan, created_at, updated_at, keterangan_tidak_hadir, keterangan_dispen, file_dispen)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING*;
	`

	//validate if payloads is keterangan hadir
	if payloads.Keterangan == "hadir" {

		payloads.Status = "hadir"

	}

	//validate if payloads is keterangan tidak hadir but keterangan_tidak_hadir is nil
	if payloads.Keterangan == "tidak hadir" {

		if payloads.KeteranganTidakHadir == "" {
			return errors.New("Failed to izin tidak hadir")
		}

		payloads.Status = "not accepted"

	}
	//validate if payloads is izin and dispen but keterangan dispen is nill
	if payloads.Keterangan == "izin" || payloads.Keterangan == "dispen" {

		if payloads.KeteranganDispen == "" {
			return errors.New("Failed to izin or dispen!")
		}

		payloads.Status = "permissions"

	}

	//execute the query
	if err := tx.QueryRowxContext(
		ctx,
		query,
		payloads.Id,
		payloads.NameLengkap,
		payloads.Kelas,
		payloads.Jurusan,
		payloads.Hari,
		payloads.Tanggal,
		payloads.Status,
		payloads.Keterangan,
		payloads.Created_at,
		payloads.Updated_at,
		payloads.KeteranganTidakHadir,
		payloads.KeteranganDispen,
		payloads.FileDispen,
	).Scan(
		&payloads.Id,
		&payloads.NameLengkap,
		&payloads.Kelas,
		&payloads.Jurusan,
		&payloads.Hari,
		&payloads.Tanggal,
		&payloads.Status,
		&payloads.Keterangan,
		&payloads.Created_at,
		&payloads.Updated_at,
		&payloads.KeteranganTidakHadir,
		&payloads.KeteranganDispen,
		&payloads.FileDispen,
	); err != nil {
		return errors.New("Failed to execute the query!" + err.Error())
	}

	//commit the transaction
	if err := tx.Commit(); err != nil {
		return errors.New("Failed to commit the transaction")
	}

	return nil

}

// func to update the absensis especially in part of status
func (s *StoreAbsensi) UpdateStatusAbsensi(id uuid.UUID, ctx context.Context, status string) error {

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

	//make the variable
	var args []interface{}
	var settings []string
	argsID := 1

	//base query
	if status != "" {
		settings = append(settings, fmt.Sprintf("status=$%d", argsID))
		argsID++
		args = append(args, status)
	}

	//full query
	full_query := fmt.Sprintf("UPDATE absensis SET %s WHERE id = $%d", strings.Join(settings, ","), argsID)
	args = append(args, id)

	//execute the query
	result, err := tx.ExecContext(ctx, full_query, args...)
	if err != nil {
		return errors.New("Failed to update the status!" + err.Error())
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.New("Failed to detect the rows affected based on db" + err.Error())
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

// func to update the keterangan tidak hadir at absensis table
func (s *StoreAbsensi) UpdateKeteranganTidakHadirAbsensi(id uuid.UUID, ctx context.Context, keterangan_tidak_hadir string) error {

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

	//make the variable
	var args []interface{}
	var settings []string
	argsID := 1

	//base query
	if keterangan_tidak_hadir != "" {
		settings = append(settings, fmt.Sprintf("status=$%d", argsID))
		argsID++
		args = append(args, keterangan_tidak_hadir)
	}

	//full query
	full_query := fmt.Sprintf("UPDATE absensis SET %s WHERE id = $%d", strings.Join(settings, ","), argsID)
	args = append(args, id)

	//execute the query
	result, err := tx.ExecContext(ctx, full_query, keterangan_tidak_hadir)
	if err != nil {
		return errors.New("Failed to update the keterangan_tidak_hadir!")
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

func (s *StoreAbsensi) GetAbsensiById(id uuid.UUID, ctx context.Context) (*types.Absensi, error) {

	//query
	query := `
		SELECT id, name_lengkap, kelas, jurusan, hari, tanggal, status, keterangan,
			   created_at, updated_at, keterangan_tidak_hadir, keterangan_dispen, file_dispen
		FROM absensis WHERE id = $1;
	`

	//execute the query
	var absensis types.Absensi
	if err := s.db.GetContext(ctx, &absensis, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid rows!")
		}
		return nil, errors.New("Failed to get the absensis data by id!" + err.Error())
	}

	//return final result
	return &absensis, nil

}
