package absensis

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
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
		INSERT INTO absensis (id, name_lengkap, kelas, jurusan, hari, tanggal, status, keterangan, created_at, updated_at
		keterangan_tidak_hadir, keterangan_dispen, file_dispen)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING*;
	`

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
func (s *StoreAbsensi) UpdateStatusAbsensi(ctx context.Context, status string) error {

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
		UPDATE absensis SET status = $1;
	`

	result, err := tx.ExecContext(ctx, query, status)
	if err != nil {
		return errors.New("Failed to update the status!")
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

// func to update the keterangan tidak hadir at absensis table
func (s *StoreAbsensi) UpdateKeteranganTidakHadirAbsensi(ctx context.Context, keterangan_tidak_hadir string) error {

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
		UPDATE absensis SET keterangan_tida_hadir = $1;
	`

	result, err := tx.ExecContext(ctx, query, keterangan_tidak_hadir)
	if err != nil {
		return errors.New("Failed to update the status!")
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
