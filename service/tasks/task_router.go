package tasks

import (
	"context"
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

type HandleTaskRequest struct {
	db types.TaskStore
}

// make func to handler the request for every single method in this task table
func NewHandlerTask(db types.TaskStore) *HandleTaskRequest {
	return &HandleTaskRequest{db: db}
}

// func to create a new task for a student
func (h *HandleTaskRequest) Create_TaskBp(w http.ResponseWriter, r *http.Request) {

	//get the middleware for a request id
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

	//get the role and validate role user
	role_user, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role user from middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role user from middleware", err.Error())
		return
	}
	if role_user != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!", false)
		return
	}

	//settings the maxbyte of the file form to input task
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	//settings the multipart form data for this method
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parse the multipart form data for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the multipart form data!", err.Error())
		return
	}

	//make the formfile for this method and validate if the input user is nill
	name_task := r.FormValue("name_task")
	student_id := r.FormValue("student_id")
	mapel_task := r.FormValue("mapel_task")
	if name_task == "" && student_id == "" && mapel_task != "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to detect the form file in request", false)
		return
	}
	//convert the student id into an uuid type
	student_id_fix, err := uuid.Parse(student_id)
	if err != nil {
		//logger the response error for this method\
		logger.Log.Error("Failed to convert into an type uuid!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the student id into an uuid type!", err.Error())
		return
	}
	if student_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get settings the student uuid!", false)
		return
	}

	//make the form file to input the file task file
	file_task, header, err := r.FormFile("file_task")
	if err != nil {
		//logger the data response for this method to check there is an error or not
		logger.Log.Error("Failed to settings the formfile to task_file",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the form file for this method!", err.Error())
		return
	}
	buff := make([]byte, 255)
	read_file, err := file_task.Read(buff)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to read the file data task!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to read the form file data!", err.Error())
		return
	}
	if read_file == 0 {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to read the length of the data!", false)
		return
	}

	//check the content type for the file task
	check_content_type := http.DetectContentType(buff)
	if check_content_type != "image/jpg" && check_content_type != "image/png" && check_content_type != "image/jpeg" {
		utils.ResponseError(w, http.StatusBadRequest, "Invalid content type for the file task!", false)
		return
	}

	//make the file name for the file task and folder that place to save the file task url
	file_name := uuid.New().String() + header.Filename
	if file_name == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the name for the task file", false)
		return
	}
	path_folder := `uploadsTask`
	if err := os.MkdirAll(path_folder, os.ModePerm); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to settings the os for this path!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the os for this path", err.Error())
		return
	}
	file_path := filepath.Join(path_folder, file_name)
	if file_path == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to combine the filename with the file path!", false)
		return
	}

	//create the folder
	folder_create, err := os.Create(file_path)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to create the path folder for a task file!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create the folder name to this method!", err.Error())
		return
	}
	defer folder_create.Close()

	//copy into a io reader to validate the file
	dst, err := io.Copy(folder_create, file_task)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to copy the filepath with a original data file!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the file name!", err.Error())
		return
	}
	if dst == 0 {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to written the data for a file task!", false)
		return
	}

	//declare the time for updated_at and created_at
	time_updated := time.Now().UTC()
	time_created := time.Now().UTC()
	time_date_task := time.Now().UTC()

	//declare the time format for response task
	time_updated_format := time.Now().UTC().Format("2006-01-02")
	time_created_format := time.Now().UTC().Format("2006-01-02")
	time_date_task_format := time.Now().UTC().Format("2006-01-02")

	//parsing into a payload
	payload := types.Payload{
		Id:         uuid.New(),
		Name_Task:  name_task,
		File_Task:  file_path,
		Date_Task:  time_date_task,
		Student_Id: student_id_fix,
		Created_at: time_created,
		Updated_at: time_updated,
	}

	//validate the payload
	var validate *validator.Validate
	validate = validator.New()
	if err := validate.Struct(&payload); err != nil {
		//logger the error response data for this method
		logger.Log.Error("Failed to check the validator for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		//use the invalid validation errors for this method
		if _, ok := err.(*validator.InvalidValidationError); !ok {
			utils.ResponseError(w, http.StatusBadRequest, "Error detected!", err)
		}
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the validation error!", err.Error())
		return
	}

	//parsing payload into a task main struct
	tasks := &types.Task{
		Id:         payload.Id,
		Name_Task:  payload.Name_Task,
		File_Task:  payload.File_Task,
		Date_Task:  payload.Date_Task,
		Student_Id: payload.Student_Id,
		Created_at: payload.Created_at,
		Updated_at: payload.Updated_at,
		Mapel_Task: payload.MapelTask,
	}

	//execute the query from task store
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.CreateNewTasks(ctx, tasks); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to create a new task!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to create a new task for a student!", err.Error())
		return
	}

	//validate the task file
	ctx_id, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	task_data, err := h.db.GetTaskById(tasks.Id, ctx_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get tasks data by id!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the task data by id!", err.Error())
		return
	}
	if task_data.File_Task != "" {
		file_old := task_data.File_Task
		if _, err := os.Stat(file_old); os.IsNotExist(err) {
			if err := os.Remove(file_old); err != nil {
				//logger the response error for this method
				logger.Log.Error("Failed to remove the old file!",
					zap.String("request_id", request_id),
					zap.String("client_ip", r.RemoteAddr),
				)
				utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the data!", err.Error())
				return
			}
			//logger the response error for this method
			logger.Log.Error("Failed to check the data is exist or not!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to checking the file is exist or not!", err.Error())
			return
		}
	}

	//parsing into a task response
	task_response := types.ResponseTask{
		Id:         tasks.Id,
		Name_Task:  tasks.Name_Task,
		File_Task:  tasks.File_Task,
		Date_Task:  time_date_task_format,
		Student_Id: tasks.Student_Id,
		Created_at: time_created_format,
		Updated_at: time_updated_format,
		MapelTask:  tasks.Mapel_Task,
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Create a new task has been successfully!", task_response)

}

// func to see the file task from directory
func (h *HandleTaskRequest) ReadFile_Bp(w http.ResponseWriter, r *http.Request) {

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

	//get the params for a file name to see the file is works to see
	vars_filename := mux.Vars(r)
	if vars_filename == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the filename params in url!", false)
	}
	filename := vars_filename["filename"]
	if filename == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params for a url!", false)
		return
	}

	//join the file name with a path folder
	path_folder := "/uploadsTask"
	file_path := filepath.Join(path_folder, filename)
	if file_path == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to join the file path for a task file!", false)
		return
	}

	//serve the file to http response
	http.ServeFile(w, r, file_path)

}

// func to delete the task in this method
func (h *HandleTaskRequest) Delete_Bp(w http.ResponseWriter, r *http.Request) {

	//get the request id from middlware for this method
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

	//get the role students from middleware and validate wich role that have to access this method
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role students from middleware",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role students from middleware!", err.Error())
		return
	}
	if role_students != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method, invalid role students!", false)
		return
	}

	//get the params for id tasks
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the params url for a request tasks!", false)
		return
	}
	task_id := vars_id["task_id"]
	if task_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to put the params on url id!", false)
		return
	}

	//convert the params for task id into an uuid
	task_id_fix, err := uuid.Parse(task_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to convert the type string into an uuid!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert the type from string into an uuid!", err.Error())
	}
	if task_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid parsing data!", false)
		return
	}

	//execute the query for this method
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.DeleteTask(task_id_fix, ctx); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to delete the task by id!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to delete the task by id!", err.Error())
		return
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Delete task has been successfully!", true)

}

// func to handle routes for update tasks students
func (h *HandleTaskRequest) UpdateTask_Bp(w http.ResponseWriter, r *http.Request) {

	//get request id from middleware token
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

	//get role from middleware token and validate wich role that can access this method
	role_students, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role students from middleware",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get role students from middleware token!", err.Error())
		return
	}
	if role_students != "siswa" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!, invalid role students!", false)
		return
	}

	//settings the params for id
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the id params!", false)
		return
	}
	task_id := vars_id["task_id"]
	if task_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the user id string for this method!", false)
		return
	}

	//parsing into an uuid
	task_id_fix, err := uuid.Parse(task_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to convert the data into an uuid type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to convert the type into an uuid type!", err.Error())
		return
	}
	if task_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid type for this method!", false)
		return
	}

	//settings the http max byte for a request multipart form
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	//parsing a request for multipart form data
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing the request to a multipart form data type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing a multipart form data for a client request!", err.Error())
		return
	}

	//settings the form value for a multipart form request
	var payloads types.PayloadUpdate
	name_task := r.FormValue("name_task")
	student_id := r.FormValue("student_id")
	mapel_task := r.FormValue("mapel_task")

	//convert the student id into an uuid type
	student_id_fix, err := uuid.Parse(student_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing the uuid type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the student id value!", err.Error())
		return
	}
	if student_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid student!", false)
		return
	}

	//settings for a file task to a file type for a request
	file_task, header, err := r.FormFile("file_task")
	if err != nil {
		if err == http.ErrMissingFile {
			//logger the response error for this method
			logger.Log.Error("Failed to check the file is missing or not!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the file task! file is missing", err.Error())
			return
		}
	}
	if err == nil {

		//make the some places to put the data of the file into it
		buff := make([]byte, 512)
		readFile, err := file_task.Read(buff)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to read the file!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read the length of file!", err.Error())
			return
		}
		if readFile == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to read the length of file!", false)
			return
		}

		//detect the content type
		content_type := http.DetectContentType(buff)
		if content_type != "image/jpg" && content_type != "image/jpeg" && content_type != "image/png" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to choose the file tyoe, invalid file type!", false)
			return
		}

		//make the filename for this method
		file_name := uuid.New().String() + header.Filename
		if file_name == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the file_name!", false)
			return
		}
		path_folder := "uploadsTaskUpdate"
		path_file := filepath.Join(path_folder, file_name)
		if path_file == "" {
			utils.ResponseError(w, http.StatusBadRequest, "Failed to join the path file for this method!", false)
			return
		}

		//create the folder into an os
		dst, err := os.Create(path_file)
		if err != nil {
			//logger response error for this method
			logger.Log.Error("Failed to create the os path for this method!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to create the os path for this method!", err.Error())
			return
		}
		defer dst.Close()

		//copy the file into an io reader
		copy_file, err := io.Copy(dst, file_task)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to copy the file name into an io reader",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to copy the file into an io reader!", err.Error())
			return
		}
		if copy_file == 0 {
			utils.ResponseError(w, http.StatusBadRequest, "Invalid length of file task!", false)
			return
		}

		//if the task file is exist, it will be replace the old path to the new path in folder
		//it could be more efficient for a memory in database and our systems
		ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
		defer cancle()
		task_file_data, err := h.db.GetTaskById(task_id_fix, ctx)
		if err != nil {
			//logger the response error for this method
			logger.Log.Error("Failed to get the task data from task table by id!",
				zap.String("request_id", request_id),
				zap.String("client_ip", r.RemoteAddr),
			)
			utils.ResponseError(w, http.StatusBadRequest, "Failed to get the task data by id!", err.Error())
			return
		}
		if task_file_data.File_Task != "" {
			task_file_old := task_file_data.File_Task
			if _, err := os.Stat(task_file_old); os.IsNotExist(err) {
				if err := os.Remove(task_file_old); err != nil {
					//logger the response error for this method
					logger.Log.Error("Failed to remove the old file!",
						zap.String("request_id", request_id),
						zap.String("client_ip", r.RemoteAddr),
					)
					utils.ResponseError(w, http.StatusBadGateway, "Failed to remove the old path!", err.Error())
					return
				}
				//logger the response error for this method
				logger.Log.Error("Failed to check the file is exist or not!",
					zap.String("request_id", request_id),
					zap.String("client_ip", r.RemoteAddr),
				)
				utils.ResponseError(w, http.StatusBadRequest, "Failed to check the file is exist or not", err.Error())
				return
			}

			//parsing into a payload
			payloads.File_Task = &path_file

		}

	}

	//condition that macth with store
	if name_task != "" {
		payloads.Name_Task = &name_task
	}
	if student_id != "" {
		payloads.Student_Id = &student_id_fix
	}
	if mapel_task != "" {
		payloads.MapelTask = &mapel_task
	}

	//execute the methods from store
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	if err := h.db.UpdateTask(task_id_fix, ctx, payloads); err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to update the task!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to update the task!", err.Error())
		return
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Update task has been successfully!", true)

}

// func to handle the get task by id that includes the students
func (h *HandleTaskRequest) GetByIdIncludeStudents_Bp(w http.ResponseWriter, r *http.Request) {

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

	//get the role from middleware
	role, err := middleware.GetRoleMiddleware(w, r)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to get the role from middleware!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the role from middleware", err.Error())
		return
	}
	if role != "guru" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to access this method!, invalid role!", false)
		return
	}

	//get the id from params url
	vars_id := mux.Vars(r)
	if vars_id == nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to settings the params id!", false)
		return
	}
	task_id := vars_id["id"]
	if task_id == "" {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the task parameters id!", false)
		return
	}

	//parsing into an uuid
	task_id_fix, err := uuid.Parse(task_id)
	if err != nil {
		//logger the response error for this method
		logger.Log.Error("Failed to parsing into an uuid type!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to parsing the data into an uuid type!", err.Error())
		return
	}
	if task_id_fix == uuid.Nil {
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the uuid type for this method!", false)
		return
	}

	//execute the query for get task by id
	// ctx, cancle := context.WithTimeout(r.Context(), time.Second * 10)
	// defer cancle()
	// task_id_validate, err := h.db.GetTaskById(id, ctx)
	// if err != nil {
	// 	//logger the response error for this method
	// 	logger.Log.Error("Failed to get the task data by id!",
	// 		zap.String("request_id", request_id),
	// 		zap.String("client_ip", r.RemoteAddr),
	// )
	// 	utils.ResponseError(w, http.StatusBadRequest, "Failed to get the task data by id!", err.Error())
	// 	return
	// }

	//execute the query
	ctx, cancle := context.WithTimeout(r.Context(), time.Second*10)
	defer cancle()
	tasks_data, err := h.db.GetTaskByIdIncludeStudents(task_id_fix, ctx)
	if err != nil {
		//logger the responnse error for this method
		logger.Log.Error("Failed to get the task data for this method!",
			zap.String("request_id", request_id),
			zap.String("client_ip", r.RemoteAddr),
		)
		utils.ResponseError(w, http.StatusBadRequest, "Failed to get the task data from table db!", err.Error())
		return
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Get data students has been successfully!", tasks_data)

}
