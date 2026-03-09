package students

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
	"go.uber.org/zap"
)

// type handlerequest that declare the student store for a database logic
type HandleRequest struct {
	db types.StudentStore
}

// func that declare the handler for student
func NewHandlerStudent(db types.StudentStore) *HandleRequest {
	return &HandleRequest{db: db}
}

// func to create a new student
func (h *HandleRequest) RegisterAsStudent_Bp(w http.ResponseWriter, r *http.Request) {

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

	//get the role for user from middleware
	role_user, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role user from middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role of user from middleware", err.Error())
		return
	}
	if role_user != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method! Invalid role!", false)
		return
	}

	//make decode for the payload of the students payload struct for register and update
	var payloads types.RegisterAsStudent
	if err := utils.DecodeData(r, &payloads); err != nil {
		//logger the data response of error for this method
		logger.Log.Error("Failed to decode the data payloads of student register struct!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the payloads of the register struct", err.Error())
		return
	}

	//make the validator for a payloads of the students register struct
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payloads); err != nil {
		//logger the response data error for this method
		logger.Log.Error("Failed to make the validator for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		var errors []string
		for _, Err := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Error detected: %s and %s", Err.ActualTag(), Err.Field()))
		}
	}

	//parsing the payload into a students struct in types
	students := &types.Student{
		Id:             uuid.New(),
		Full_Name:      payloads.Full_name,
		Kelas:          payloads.Kelas,
		Jurusan:        payloads.Jurusan,
		Absen:          payloads.Absen,
		StudentProfile: payloads.StudentProfile,
		Wali_Kelas:     payloads.Wali_Kelas,
		Created_at:     payloads.Created_at,
		Updated_at:     payloads.Updated_at,
	}

	//execute the query from student store
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.CreateNewStudent(ctx, students); err != nil {
		//logger the response error data for this method
		logger.Log.Error("Failed to create a new students for a user!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create a new students for user", err.Error())
		return
	}

	//make the time parsin for a response students
	time_created := time.Now().UTC().Format("2006-01-02")
	time_updated := time.Now().UTC().Format("2006-01-02")

	//make the response for a new students
	students_response := types.StudentResponse{
		Id:             students.Id,
		Full_name:      students.Full_Name,
		Kelas:          students.Kelas,
		Jurusan:        students.Jurusan,
		Absen:          students.Absen,
		StudentProfile: students.StudentProfile,
		Wali_Kelas:     students.Wali_Kelas,
		Created_at:     time_created,
		Updated_at:     time_updated,
	}

	//return final result for this method
	utils.ResponseSuccess(w, http.StatusCreated, "Create a new student has been successfully!", students_response)

}
