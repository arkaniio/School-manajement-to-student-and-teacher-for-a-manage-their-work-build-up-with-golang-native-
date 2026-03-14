package students

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// type handlerequest that declare the student store for a database logic
type HandleStudentsRequest struct {
	db types.StudentStore
}

// func that declare the handler for student
func NewHandlerStudent(db types.StudentStore) *HandleStudentsRequest {
	return &HandleStudentsRequest{db: db}
}

// func to create a new student
func (h *HandleStudentsRequest) RegisterAsStudent_Bp(w http.ResponseWriter, r *http.Request) {

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
		MapelStudents:  payloads.MapelStudents,
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
		MapelStudents:  students.MapelStudents,
	}

	//return final result for this method
	utils.ResponseSuccess(w, http.StatusCreated, "Create a new student has been successfully!", students_response)

}

// func to delete the students (only admin and teacher can do and access this method)
func (h *HandleStudentsRequest) DeleteStudent_Bp(w http.ResponseWriter, r *http.Request) {

	//make the request id from logger middleware
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

	//get the role in middleware token to check the role of the users students or admin and or teacher
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the data response error for this method
		logger.Log.Error("Failed to get the role middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role from middleware token role", err.Error())
		return
	}
	if role_students != "admin" && role_students != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method, invalid role!", false)
		return
	}

	//get the params id from params url
	vars_params := mux.Vars(r)
	if vars_params == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params id from url!", false)
		return
	}
	user_id := vars_params["id"]
	if user_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id from url params", false)
		return
	}

	//convert into an uuid type for user_id
	uuid_user, err := uuid.Parse(user_id)
	if err != nil {
		//logger the response data error from this method
		logger.Log.Error("Failed to convert into an uuid type for user_id!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert type string into an uuid!", err.Error())
		return
	}
	if uuid_user == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Nil!", false)
		return
	}

	//execute the query from store method!
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.DeleteStudents(uuid_user, ctx); err != nil {
		//logger the data response error for this method
		logger.Log.Error("Failed to delete the students data!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the students data!", err.Error())
		return
	}

	//return final response
	utils.ResponseError(w, http.StatusOK, "Delete data has been successfully!", true)

}

// func to handle the routes as a part of routes in this method
func (h *HandleStudentsRequest) UpdateStudents_Bp(w http.ResponseWriter, r *http.Request) {

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

	//get the role students user from middleware token id
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
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!, Invalid role students!", false)
		return
	}

	//get the params for a id in url endpoint
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params id at url endpoint because nill!", false)
		return
	}
	user_id := vars_id["id"]
	if user_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the string of id parameter in url endpoint!", false)
		return
	}

	//parsing into an uuid type
	uuid_user, err := uuid.Parse(user_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to convert from string type into an uuid type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convet from string type into an uuid type!", err.Error())
		return
	}
	if uuid_user == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid result of uuid convert!", false)
		return
	}

	//parsing into a max byte of file type for a request from every students
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	//parse every requests to a multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing every requests from request students!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the multipart form data", err.Error())
		return
	}

	//declare the variable form file
	var payloads types.UpdateAsStudent
	full_name := r.FormValue("full_name")
	kelas := r.FormValue("kelas")
	jurusan := r.FormValue("jurusan")
	wali_kelas := r.FormValue("wali_kelas")
	absen := r.FormValue("absen")
	mapel_students := r.FormValue("mapel_students")

	//settings the form file to a student profile field in db
	student_profile_name, header, err := r.FormFile("student_profile")
	if err != nil {
		//check if th error is type of missing file error
		if err == http.ErrMissingFile {
			//logger the response error for this method
			logger.Log.Error("Failed to check the error for type missing file error!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to cheking the error type!", err.Error())
			return
		}
	}
	//because this is the update method, we have to detect if the error is nill, it will be read the file and
	//put the logic on that condition to make sure if students is right if they not update their students profile
	if err == nil {

		//reade the file and detect the content type for this file
		buff := make([]byte, 512)
		read, err := student_profile_name.Read(buff)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to read the file buff!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read the file buff !", err.Error())
			return
		}
		if read == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Invalid length of the read file!", false)
			return
		}
		content_type_file := http.DetectContentType(buff)
		if content_type_file != "image/jpg" && content_type_file != "image/jpeg" && content_type_file != "image/png" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to access this file! Invalid content type", false)
			return
		}

		//create the filename and os name
		file_name_student := uuid.New().String() + header.Filename
		if file_name_student == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the file name student profile!", false)
			return
		}
		path_os := "uploadsStudent"
		if err := os.MkdirAll(path_os, os.ModePerm); err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to settings the os for this file!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the os for this file!", err.Error())
			return
		}
		path_file := filepath.Join(path_os, file_name_student)
		if path_file == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get settings the file name!", false)
			return
		}

		//create os for a folder to this students file
		dst, err := os.Create(path_file)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to create the os file!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to create the os file based on file name on this method!", err.Error())
			return
		}
		defer dst.Close()

		//copy the file into an io reader to communicate with a request from students
		copy, err := io.Copy(dst, student_profile_name)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to copy the file!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the file name!", err.Error())
			return
		}
		if copy == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Invalid io reader!", false)
			return
		}

		//check if the students profile is exist in folder, if the students is updating again their students profile
		//its gonna be replace from old path to a new path in folder students profile
		ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
		defer cancle()
		students, err := h.db.GetStudentById(uuid_user, ctx)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to get the students id from db!",
				zap.String("request_id", request_id),
				zap.String("cllient_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get students data from db by id!", err.Error())
			return
		}
		if students.StudentProfile != "" {
			path_old := students.StudentProfile
			if _, err := os.Stat(path_old); os.IsNotExist(err) {
				if err := os.Remove(path_old); err != nil {
					//logger the response error for this method
					logger.Log.Error("Failed to remove the data old from db!",
						zap.String("request_id", request_id),
						zap.String("client_ip", r.RemoteAddr),
					)
					utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the old path!", err.Error())
					return
				}
				//logger the response error for this method
				logger.Log.Error("Failed to check the file old is exist or not!",
					zap.String("request_id", request_id),
					zap.String("client_ip", r.RemoteAddr),
				)
				utils.ResponseError(w, http.StatusBadRequest, "Failed to check the data old file!", err.Error())
				return
			}
		}

		//if the payload in students profile is not nill
		payloads.StudentProfile = &path_file

	}

	//checking and validate again
	if full_name != "" {
		payloads.Full_name = &full_name
	}
	if kelas != "" {
		payloads.Kelas = &kelas
	}
	if jurusan != "" {
		payloads.Jurusan = &jurusan
	}
	if absen != "" {
		//convert the value of absen into an integer
		absen_fix, err := strconv.Atoi(absen)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to convert from string into an integer for a absen!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to convert type strings into an integer for a absen!", err.Error())
			return
		}
		payloads.Absen = &absen_fix

	}
	if wali_kelas != "" {
		payloads.Wali_Kelas = &wali_kelas
	}
	if mapel_students != "" {
		payloads.MapelStudents = &mapel_students
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.UpdateStudentsData(uuid_user, payloads, ctx); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to update and execute the query!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update and execute the query!", err.Error())
		return
	}

	//return final result and message
	utils.ResponseSuccess(w, http.StatusOK, "Update data students has been successfully!", true)

}

// func to get the filename and we can opent it into a postman or url
func (h *HandleStudentsRequest) ReadFilename_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from request user
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

	//make the params for the path in url postman
	vars_fileName := mux.Vars(r)
	if vars_fileName == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the params into an path in postman", false)
		return
	}
	file_nameParams := vars_fileName["file_name"]
	if file_nameParams == "" {
		utils.ResponseError(w, http.StatusBadRequest, "No one params detected in this route!", false)
		return
	}

	//combine path
	path_folder := "uploadsStudent"
	file_name := filepath.Join(path_folder, file_nameParams)

	//check if the file is exist or not
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		//logger the response error for this method
		logger.Log.Error("Failed to check the data is exist or not!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to check the file name is exist or not!", err.Error())
		return
	}

	//serve the file_name in http
	http.ServeFile(w, r, file_name)

}

// func to handle the get all students store
func (h *HandleStudentsRequest) GetAllStudents_Bp(w http.ResponseWriter, r *http.Request) {

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

	//get the role and validate the role in this method
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role students!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role students from middleware token!", err.Error())
		return
	}
	if role_students != "guru" && role_students != "admin" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!", false)
		return
	}

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	students_data, err := h.db.GetAllStudents(ctx)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get all the data of students!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the all of students data!", err.Error())
		return
	}
	//debug
	fmt.Print(students_data)

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Get all students data has been successfully!", students_data)

}
