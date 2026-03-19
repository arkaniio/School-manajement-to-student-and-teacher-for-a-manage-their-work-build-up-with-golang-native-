package studentsgrades

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
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// make the handler type for a student_grades
type HanlderStudentsGrade struct {
	repo    types.StudentsGradeStore
	service *ServiceStudentGrade
}

// make the func to get handler absensis
func NewHandlerStudentsGrade(db *sqlx.DB) *HanlderStudentsGrade {
	repo := NewHandlerStoreStudentsGrade(db)
	service := NewServiceStudentsGrade(repo)
	return &HanlderStudentsGrade{
		repo:    repo,
		service: service,
	}
}

// add func to create a new students greade for every students
func (h *HanlderStudentsGrade) CreateNewStudentsGrade_Bp(w http.ResponseWriter, r *http.Request) {

	//get request if from middleware
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
	role_grades, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role grades from middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role grades from middleware!", err.Error())
		return
	}
	if role_grades != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Cannot access this method, invalid role!", false)
		return
	}

	//decode the payloads
	var payloads types.PayloadsStudentGrade
	if err := utils.DecodeData(r, &payloads); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to decode the payload!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode the data payloads!", err.Error())
		return
	}

	//validate the payloads
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payloads); err != nil {
		var errors []string
		for _, ErrMsg := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Error detected for required! %s, %s", ErrMsg.ActualTag(), ErrMsg.Field()))
		}
		//logger the response error for this method
		logger.Log.Error("Failed to settings the validator for this method",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the validator!", err.Error())
		return
	}

	//parsing into an another type
	students_grades := &types.StudentsGrade{
		Id:         payloads.Id,
		Task_Id:    payloads.Task_Id,
		Tanggal:    payloads.Tanggal,
		Keterangan: payloads.Keterangan,
		Grades:     payloads.Grades,
		Created_at: payloads.Created_at,
		Updated_at: payloads.Updated_at,
	}

	//execute the method from repository
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.repo.CreateNewStudentsGrade(ctx, students_grades); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to create new task grade!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create a new task grade!", err.Error())
		return
	}

	//make the students_grade response
	students_grade_response := types.StudentsGradeResponse{
		Id:         students_grades.Id,
		Task_Id:    students_grades.Task_Id,
		Tanggal:    students_grades.Tanggal,
		Keterangan: students_grades.Keterangan,
		Grades:     students_grades.Grades,
		Created_at: utils.SetResponseTime("2006-01-02"),
		Updated_at: utils.SetResponseTime("2006-01-02"),
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusCreated, "Create a new grades has been successfully!", students_grade_response)

}

func (h *HanlderStudentsGrade) UpdateStudentsGrade_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from middleware
	request_id := middleware.GetRequestID(r)
	if request_id == "" {
		logger.Log.Info("Failed to get the request id from update func!",
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get request id!", false)
		return
	}

	//get the role students from middlewae
	role_grades, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		logger.Log.Error("Failed to get role for update!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role!", err.Error())
		return
	}
	if role_grades != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Guru role required!", false)
		return
	}

	//setup the params id
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("Invalid grade ID!",
			zap.String("request_id", request_id),
			zap.String("id", idStr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid grade ID!", err.Error())
		return
	}

	//decode the data payload
	var payloads types.PayloadsStudentGradeUpdate
	if err := utils.DecodeData(r, &payloads); err != nil {
		logger.Log.Error("Failed to decode update payload!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to decode payload!", err.Error())
		return
	}

	//make the validator for a payloads update
	validate := validator.New()
	if err := validate.Struct(&payloads); err != nil {
		logger.Log.Error("Validation failed for update!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Validation failed!", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()
	if err := h.repo.UpdateStudentsGrade(ctx, id, &payloads); err != nil {
		logger.Log.Error("Failed to update students grade!",
			zap.String("request_id", request_id),
			zap.String("id", id.String()),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update grade!", err.Error())
		return
	}

	students_grade_response := types.StudentsGradeResponse{
		Id: id,
	}

	utils.ResponseSuccess(w, http.StatusOK, "Grade updated successfully!", students_grade_response)
}

func (h *HanlderStudentsGrade) DeleteStudentsGrade_Bp(w http.ResponseWriter, r *http.Request) {

	//get request id from middleware
	request_id := middleware.GetRequestID(r)
	if request_id == "" {
		logger.Log.Info("Failed to get request id for delete!",
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get request id!", false)
		return
	}

	//get the role students from middleware
	role, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		logger.Log.Error("Failed to get role for delete!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role!", err.Error())
		return
	}
	if role != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Guru role required!", false)
		return
	}

	//setup the params id
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.Error("Invalid ID for delete!",
			zap.String("request_id", request_id),
			zap.String("id", idStr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Invalid ID!", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	if err := h.repo.DeleteStudentsGrade(ctx, id); err != nil {
		logger.Log.Error("Failed to delete grade!",
			zap.String("request_id", request_id),
			zap.String("id", idStr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete grade!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Grade deleted successfully!", map[string]interface{}{"message": "Deleted", "id": id.String()})
}

func (h *HanlderStudentsGrade) GetAllStudentsGradesWithTask_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from middleware
	request_id := middleware.GetRequestID(r)
	if request_id == "" {
		logger.Log.Info("Failed to get request id for get all!",
			zap.String("client_ip", r.RemoteAddr),
			zap.String("path", r.URL.Path),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get request id!", false)
		return
	}

	//get the role students from middleware
	role, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		logger.Log.Error("Failed to get role for get all!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role!", err.Error())
		return
	}
	if role != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Guru role required!", false)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
	defer cancel()

	grades, err := h.repo.GetAllStudentsGradesWithTask(ctx)
	if err != nil {
		logger.Log.Error("Failed to get all grades!",
			zap.String("request_id", request_id),
		)
		utils.ResponseError(w, http.StatusInternalServerError, "Failed to get grades!", err.Error())
		return
	}

	utils.ResponseSuccess(w, http.StatusOK, "Grades retrieved successfully!", grades)
}
