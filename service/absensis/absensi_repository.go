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
func (s *StoreAbsensi) CreateNewAbsensi(ctx context.Context, payloads *types.PayloadAbsensis) error {

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
		INSERT INTO absensis (id, name_lengkap, kelas, jurusan, hari, tanngal, status, keterangan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	//execute the query
	result, err := tx.ExecContext(ctx, query, payloads)
	if err != nil {
		return errors.New("Failed to execute the query!")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("Failed to detect the rows into an db")
		}
		return errors.New("Failed to get the rows into an db!")
	}
	if rows == 0 {
		return errors.New("Failed to detect the rows in db!")
	}

	return nil

}
