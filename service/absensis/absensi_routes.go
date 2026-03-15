package absensis

import (
	"net/http"

	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
)

// make the handler type for a absensi
type HanlderAbsensi struct {
	db types.AbsensiStore
}

// make the func to get handler absensis
func NewHandlerAbsensi(db types.AbsensiStore) *HanlderAbsensi {
	return &HanlderAbsensi{db: db}
}

// func to create the absensi
func (h *HanlderAbsensi) CreateNewAbsensi_Bp(w http.ResponseWriter, r *http.Request) {

}
