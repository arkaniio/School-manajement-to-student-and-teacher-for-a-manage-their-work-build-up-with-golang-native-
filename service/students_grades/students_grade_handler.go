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
	"go.uber.org/zap"
)

// make the handler type for a student_grades
type HanlderStudentsGrade struct {
	db types.StudentsGradeStore
}

// make the func to get handler absensis
func NewHandlerStudentsGrade(db types.StudentsGradeStore) *HanlderStudentsGrade {
	return &HanlderStudentsGrade{db: db}
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
	if err := h.db.CreateNewStudentsGrade(ctx, students_grades); err != nil {
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
