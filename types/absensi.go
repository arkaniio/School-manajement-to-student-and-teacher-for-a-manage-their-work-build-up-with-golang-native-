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
	DeleteAbsensisById(id uuid.UUID, ctx context.Context) error
	UpdateAbsensiById(id uuid.UUID, ctx context.Context, payloads PayloadAbsensisUpdate) error

	// New stats methods
	GetWeeklyStats(ctx context.Context) (*AbsensiStats, error)
	GetMonthlyStats(ctx context.Context) (*AbsensiStats, error)
	GetAllAbsensiWithStudents(ctx context.Context) ([]AbsensiWithStudent, error)
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

type PayloadAbsensisUpdate struct {
	Id                   uuid.UUID `json:"id"`
	NameLengkap          *string   `json:"name_lengkap" validate:"required"`
	Kelas                *string   `json:"kelas" validate:"required"`
	Jurusan              *string   `json:"jurusan" validate:"required"`
	Hari                 *string   `json:"hari" validate:"required"`
	Tanggal              *string   `json:"tanggal" validate:"required"`
	Status               *string   `json:"status"`
	Keterangan           *string   `json:"keterangan" validate:"required"`
	Created_at           time.Time `json:"created_at"`
	Updated_at           time.Time `json:"updated_at"`
	KeteranganTidakHadir *string   `json:"keterangan_tidak_hadir"`
	KeteranganDispen     *string   `json:"keterangan_dispen"`
	FileDispen           *string   `json:"file_dispen"`
}

type AbsensiStats struct {
	Hadir      int `json:"hadir"`
	TidakHadir int `json:"tidak_hadir"`
	Izin       int `json:"izin"`
}

type AbsensiWithStudent struct {
	Id                   uuid.UUID `db:"id"`
	NameLengkap          string    `db:"name_lengkap"`
	KelasAbsensi         string    `db:"kelas"`
	JurusanAbsensi       string    `db:"jurusan"`
	Hari                 string    `db:"hari"`
	Tanggal              string    `db:"tanggal"`
	Status               string    `db:"status"`
	Keterangan           string    `db:"keterangan"`
	Created_at           time.Time `db:"created_at"`
	Updated_at           time.Time `db:"updated_at"`
	KeteranganTidakHadir string    `db:"keterangan_tidak_hadir"`
	KeteranganDispen     string    `db:"keterangan_dispen"`
	FileDispen           string    `db:"file_dispen"`
	StudentFullName      string    `db:"full_name"`
	StudentKelas         string    `db:"s_kelas"`
	StudentJurusan       string    `db:"s_jurusan"`
	StudentAbsen         int       `db:"s_absen"`
	StudentProfile       string    `db:"student_profile"`
	StudentWaliKelas     string    `db:"s_wali_kelas"`
	StudentMapel         string    `db:"s_mapel_students"`
}

type AbsensiWithStudentResponse struct {
	Id                   uuid.UUID `json:"id"`
	NameLengkap          string    `json:"name_lengkap"`
	KelasAbsensi         string    `json:"kelas_absensi"`
	JurusanAbsensi       string    `json:"jurusan_absensi"`
	Hari                 string    `json:"hari"`
	Tanggal              string    `json:"tanggal"`
	Status               string    `json:"status"`
	Keterangan           string    `json:"keterangan"`
	Created_at           string    `json:"created_at"`
	Updated_at           string    `json:"updated_at"`
	KeteranganTidakHadir string    `json:"keterangan_tidak_hadir"`
	KeteranganDispen     string    `json:"keterangan_dispen"`
	FileDispen           string    `json:"file_dispen"`
	StudentFullName      string    `json:"student_full_name"`
	StudentKelas         string    `json:"student_kelas"`
	StudentJurusan       string    `json:"student_jurusan"`
	StudentAbsen         int       `json:"student_absen"`
	StudentProfile       string    `json:"student_profile"`
	StudentWaliKelas     string    `json:"student_wali_kelas"`
	StudentMapel         string    `json:"student_mapel"`
}
