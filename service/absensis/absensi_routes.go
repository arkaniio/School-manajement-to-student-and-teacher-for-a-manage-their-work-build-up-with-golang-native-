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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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
		utils.ResponseError(w, http.StatusBadRequest, "Failed to validate the payloads!", err.Error())
		return
	}

	//parsing into struct absensi
	absensi := &types.Absensi{
		Id:                   uuid.New(),
		NameLengkap:          payloads.NameLengkap,
		Kelas:                payloads.Kelas,
		Jurusan:              payloads.Jurusan,
		Hari:                 payloads.Hari,
		Tanggal:              payloads.Tanggal,
		Status:               payloads.Status,
		Keterangan:           payloads.Keterangan,
		Created_at:           payloads.Created_at,
		Updated_at:           payloads.Updated_at,
		KeteranganTidakHadir: payloads.KeteranganTidakHadir,
		KeteranganDispen:     payloads.KeteranganDispen,
		FileDispen:           payloads.FileDispen,
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.CreateNewAbsensi(ctx, absensi); err != nil {
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
		Id:                   absensi.Id,
		NameLengkap:          absensi.NameLengkap,
		Kelas:                absensi.Kelas,
		Jurusan:              absensi.Jurusan,
		Hari:                 absensi.Hari,
		Tanggal:              absensi.Tanggal,
		Status:               absensi.Status,
		Keterangan:           absensi.Keterangan,
		Created_at:           time_created,
		Updated_at:           time_updated,
		KeteranganTidakHadir: absensi.KeteranganTidakHadir,
		KeteranganDispen:     absensi.KeteranganDispen,
		FileDispen:           absensi.FileDispen,
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusCreated, "Create the new data absensi has been successfully!", response_absensi)

}

// add the handler routes for delete the absensis
func (h *HanlderAbsensi) DeleteAbsensis_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from middleware
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

	//get the role from middleware token
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the middleware role for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role from middleware!", err.Error())
		return
	}
	if role_students != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method, invalid role!", false)
		return
	}

	//get the params id
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the params for this method!", false)
		return
	}
	absensi_id := vars_id["absensi_id"]
	if absensi_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid parameters for this method!", false)
		return
	}

	//convert the id from parameters rto type uuid
	absensi_id_fix, err := uuid.Parse(absensi_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing the absensi id to type uuid!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the uuid type!", err.Error())
		return
	}
	if absensi_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid value of type uuid!", false)
		return
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.DeleteAbsensisById(absensi_id_fix, ctx); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to delete the absensi data by id!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the data of absensi!", err.Error())
		return
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Delete the absensi has been successfully!", true)

}
