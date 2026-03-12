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
	if name_task == "" && student_id == "" {
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
				utils.ResponseError(w, http.StatusBadRequest, "Failed to remove the data!", err.Error())
				return
			}
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
	}

	//return final result
	utils.ResponseSuccess(w, http.StatusOK, "Create a new task has been successfully!", task_response)

}
