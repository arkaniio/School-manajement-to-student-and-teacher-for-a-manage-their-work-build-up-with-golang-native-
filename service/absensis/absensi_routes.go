package absensis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware"
	"github.com/ArkaniLoveCoding/Shcool-manajement/middleware/logger"
	"github.com/ArkaniLoveCoding/Shcool-manajement/types"
	"github.com/ArkaniLoveCoding/Shcool-manajement/utils"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
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

	//get request from middleware
	request_id := middleware.GetRequestID(r)
	if request_id == "" {
		//make the logger data response for info
		logger.Log.Info("Failed to get the request id from this func!",
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get request id for this method!", false)
		return
	}

	//get the role from middleware and validate the role
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role students from middleware token!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role students from middleware!", err.Error())
		return
	}
	if role_students != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method, invalid role!", false)
		return
	}

	//decode the payloads from type tasks
	var payloads types.PayloadAbsensis
	if err := utils.DecodeData(r, &payloads); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to decod the payloads of the data",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the payloads!", err.Error())
		return
	}

	//make the validator for this method
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payloads); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to check the validator for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		var errors []string
		for _, Err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Error Detected: %s, %s", Err.ActualTag(), Err.Field()))
		}
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.CreateNewAbsensi(ctx, &payloads); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to create the new absensi for students!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create the new absensi for the students!", err.Error())
		return
	}

	//parsing time for created and updated
	time_created := time.Now().UTC().Format("2006-01-02")
	time_updated := time.Now().UTC().Format("2006-01-02")

	//make the response for this method
	response_absensi := types.AbsensiResponse{
		Id:          payloads.Id,
		NameLengkap: payloads.NameLengkap,
		Kelas:       payloads.Kelas,
		Jurusan:     payloads.Jurusan,
		Hari:        payloads.Hari,
		Tanggal:     payloads.Tanggal,
		Status:      payloads.Status,
		Keterangan:  payloads.Keterangan,
		Created_at:  time_created,
		Updated_at:  time_updated,
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusCreated, "Create the new data absensi has been successfully!", response_absensi)

}
