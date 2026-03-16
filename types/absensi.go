package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AbsensiStore interface {
	CreateNewAbsensi(ctx context.Context, payloads *Absensi) error
	UpdateStatusAbsensi(id uuid.UUID, ctx context.Context, status string) error
	UpdateKeteranganTidakHadirAbsensi(id uuid.UUID, ctx context.Context, keterangan_tidak_hadir string) error
	GetAbsensiById(id uuid.UUID, ctx context.Context) (*Absensi, error)
}

type Absensi struct {
	Id                   uuid.UUID `db:"id"`
	NameLengkap          string    `db:"name_lengkap"`
	Kelas                string    `db:"kelas"`
	Jurusan              string    `db:"jurusan"`
	Hari                 string    `db:"hari"`
	Tanggal              string    `db:"tanggal"`
	Status               string    `db:"status"`
	Keterangan           string    `db:"keterangan"`
	Created_at           time.Time `db:"created_at"`
	Updated_at           time.Time `db:"updated_at"`
	KeteranganTidakHadir string    `db:"keterangan_tidak_hadir"`
	KeteranganDispen     string    `db:"keterangan_dispen"`
	FileDispen           string    `db:"file_dispen"`
}

type PayloadAbsensis struct {
	Id                   uuid.UUID `json:"id"`
	NameLengkap          string    `json:"name_lengkap" validate:"required"`
	Kelas                string    `json:"kelas" validate:"required"`
	Jurusan              string    `json:"jurusan" validate:"required"`
	Hari                 string    `json:"hari" validate:"required"`
	Tanggal              string    `json:"tanggal" validate:"required"`
	Status               string    `json:"status"`
	Keterangan           string    `json:"keterangan" validate:"required"`
	Created_at           time.Time `json:"created_at"`
	Updated_at           time.Time `json:"updated_at"`
	KeteranganTidakHadir string    `json:"keterangan_tidak_hadir"`
	KeteranganDispen     string    `json:"keterangan_dispen"`
	FileDispen           string    `json:"file_dispen"`
}

type AbsensiResponse struct {
	Id                   uuid.UUID `json:"id"`
	NameLengkap          string    `json:"name_lengkap"`
	Kelas                string    `json:"kelas"`
	Jurusan              string    `json:"jurusan"`
	Hari                 string    `json:"hari"`
	Tanggal              string    `json:"tanggal"`
	Status               string    `json:"status"`
	Keterangan           string    `json:"keterangan"`
	Created_at           string    `json:"created_at"`
	Updated_at           string    `json:"updated_at"`
	KeteranganTidakHadir string    `json:"keterangan_tidak_hadir"`
	KeteranganDispen     string    `json:"keterangan_dispen"`
	FileDispen           string    `json:"file_dispen"`
}
