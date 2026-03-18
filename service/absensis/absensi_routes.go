package absensis

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

// add func to handle the update absensi repository
func (h *HanlderAbsensi) UpdateAbsensi_Bp(w http.ResponseWriter, r *http.Request) {

	//get request id from middleware
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

	//check the role students from middleware token\
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role students from middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role students from middleware!", err.Error())
		return
	}
	if role_students != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!, invalid role!", false)
		return
	}

	//get the absensi id!
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the absensi vars!", false)
		return
	}
	absensi_id := vars_id["absensi_id"]
	if absensi_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the absensi id for a parameters!", false)
		return
	}

	//convert into an uuid type
	absensi_id_fix, err := uuid.Parse(absensi_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing from string into an uuid type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing from string into and uuid type!", err.Error())
		return
	}
	if absensi_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the uuid type for absensi id!", false)
		return
	}

	//settings for a max byte reader
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	//parsing the multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing into an multipart form data!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing into an multipart form data!", err.Error())
		return
	}

	//settings the form value for this method\
	var payloads types.PayloadAbsensisUpdate
	name_lengkap := r.FormValue("name_lengkap")
	kelas := r.FormValue("kelas")
	jurusan := r.FormValue("jurusan")
	hari := r.FormValue("hari")
	tanggal := r.FormValue("tanggal")
	status := r.FormValue("status")
	keterangan := r.FormValue("keterangan")
	keterangan_tidak_hadir := r.FormValue("keterangan_tidak_hadir")
	keterangan_dispen := r.FormValue("keterangan_dispen")

	//settings the form file for file_dispen
	file_dispen, header, err := r.FormFile("file_dispen")
	if err != nil {
		if err != http.ErrMissingFile {
			//logger the response error for this method
			logger.Log.Error("Failed to get the err missing file for this method!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the file error!", err.Error())
		}
	}
	if err == nil {

		//read the file using buff
		buff := make([]byte, 512)
		read_file, err := file_dispen.Read(buff)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to read the file dispen!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read the file dispen using buff!", err.Error())
			return
		}
		if read_file == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the length of file dispen!", false)
			return
		}

		//detect the content type for file_dispen
		content_type := http.DetectContentType(buff)
		if content_type != "image/jpg" && content_type != "image/png" && content_type != "image/jpeg" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to choose the content type, invalid content type for file fispen!", false)
			return
		}

		//make the filename
		file_name := uuid.New().String() + header.Filename
		path_folder := "uploadsAbsensi"
		if err := os.MkdirAll(path_folder, os.ModePerm); err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to settings the directory for a path folder!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the file name!", err.Error())
			return
		}
		file_name_final := filepath.Join(path_folder, file_name)
		if file_name_final == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the file name for this method!", false)
			return
		}

		//create the folder using os
		dst, err := os.Create(file_name_final)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to create the new path for this method!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the os to create the new file!", err.Error())
			return
		}
		defer dst.Close()

		//copy into an io reader
		io_copy, err := io.Copy(dst, file_dispen)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to copy the io reader for this method!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the io reader!", err.Error())
			return
		}
		if io_copy == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the length of io copy!", false)
			return
		}

		//check if the file dispen is exist
		ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
		defer cancle()
		absensis, err := h.db.GetAbsensiById(absensi_id_fix, ctx)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to get the absensis data by id!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the absensi data by id!", err.Error())
			return
		}
		if absensis.FileDispen != "" {
			file_dispen_old := absensis.FileDispen
			if _, err := os.Stat(file_dispen_old); os.IsNotExist(err) {
				//logger the response error for this method
				logger.Log.Error("Failed to get the absensi file dispen!",
					zap.String("request_id", request_id),
					zap.String("client_ip", r.RemoteAddr),
				)
				utils.ResponseError(w, http.StatusBadRequest, "Failed to get the absensi file dispen from db!", err.Error())
				return
			}
			if err := os.Remove(file_dispen_old); err != nil {
				//logger the response error for this method
				logger.Log.Error("Failed to remove the file old in db to a new file!",
					zap.String("request_id", request_id),
					zap.String("client_ip", r.RemoteAddr),
				)
				utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the data old path from db and change it to a new path!", err.Error())
				return
			}

			utils.SetIsNotEmpty(&payloads.FileDispen, file_name_final)

		}

		//check the conditions
		utils.SetIsNotEmpty(&payloads.NameLengkap, name_lengkap)
		utils.SetIsNotEmpty(&payloads.Kelas, kelas)
		utils.SetIsNotEmpty(&payloads.Jurusan, jurusan)
		utils.SetIsNotEmpty(&payloads.Hari, hari)
		utils.SetIsNotEmpty(&payloads.Tanggal, tanggal)
		utils.SetIsNotEmpty(&payloads.Status, status)
		utils.SetIsNotEmpty(&payloads.Keterangan, keterangan)
		utils.SetIsNotEmpty(&payloads.KeteranganTidakHadir, keterangan_tidak_hadir)
		utils.SetIsNotEmpty(&payloads.KeteranganDispen, keterangan_dispen)

		//execute the query from repository for this method
		if err := h.db.UpdateAbsensiById(absensi_id_fix, ctx, payloads); err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to update the absensi by id!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to update the absensi data by id!", err.Error())
			return
		}

		//return final result
		utils.ResponseSuccess(w, http.StatusOK, "Update the absensi data by id has been successfully!", true)

	}

}
