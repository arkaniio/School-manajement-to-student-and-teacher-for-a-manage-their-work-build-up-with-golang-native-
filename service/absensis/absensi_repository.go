package absensis

import "github.com/jmoiron/sqlx"

//make the type for repo
type StoreAbsensi struct {
	db *sqlx.DB
}

//make the func for repository
func NewHandlerStoreAbsensi(db *sqlx.DB) *StoreAbsensi {
	return &StoreAbsensi{db: db}
}

//func to make the new absensi
