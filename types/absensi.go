package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AbsensiStore interface {
	CreateNewAbsensi(ctx context.Context, payloads *Absensi) error
	UpdateStatusAbsensi(ctx context.Context, status string) error
}

type Absensi struct {
	Id          uuid.UUID `db:"id"`
	NameLengkap string    `db:"name_lengkap"`
	Kelas       string    `db:"kelas"`
	Jurusan     string    `db:"jurusan"`
	Hari        string    `db:"hari"`
	Tanggal     string    `db:"tanggal"`
	Status      string    `db:"status"`
	Keterangan  string    `db:"keterangan"`
	Created_at  time.Time `db:"created_at"`
	Updated_at  time.Time `db:"updated_at"`
}

type PayloadAbsensis struct {
	Id          uuid.UUID `json:"id"`
	NameLengkap string    `json:"name_lengkap" validate:"required"`
	Kelas       string    `json:"kelas" validate:"required"`
	Jurusan     string    `json:"jurusan" validate:"required"`
	Hari        string    `json:"hari" validate:"required"`
	Tanggal     string    `json:"tanggal" validate:"required"`
	Status      string    `json:"status"`
	Keterangan  string    `json:"keterangan" validate:"required"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

type AbsensiResponse struct {
	Id          uuid.UUID `json:"id"`
	NameLengkap string    `json:"name_lengkap"`
	Kelas       string    `json:"kelas"`
	Jurusan     string    `json:"jurusan"`
	Hari        string    `json:"hari"`
	Tanggal     string    `json:"tanggal"`
	Status      string    `json:"status"`
	Keterangan  string    `json:"keterangan"`
	Created_at  string    `json:"created_at"`
	Updated_at  string    `json:"updated_at"`
}
