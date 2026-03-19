# Student Grades CRUD Implementation Plan

Status: In Progress

## Steps:
1. **[DONE]** Update `types/students_grade.go`: Add interface methods (Update, Delete, GetAllWithTask), new response struct StudentsGradesWithTaskResponse with task fields (name_task, mapel_task), fix Grades json tag.

2. **[DONE]** Update `service/students_grades/students_grade_repository.go`: Implement UpdateStudentsGrade (dynamic partial update with tx), DeleteStudentsGrade (tx), GetAllStudentsGradesWithTask (JOIN query scan to []StudentsGradesWithTaskResponse).

3. **[DONE]** Update `service/students_grades/students_grade_handler.go`: Add handlers UpdateStudentsGrade_Bp (PATCH {id}), DeleteStudentsGrade_Bp (DELETE {id}), GetAllStudentsGradesWithTask_Bp (GET), following existing patterns (role check, logging, validation).

4. **[DONE]** Update `cmd/api/router_api.go`: Register new routes PATCH/DELETE /grades/{id}, GET /grades with middleware.TokenIdMiddleware.

5. **[DONE]** ✅ All features implemented and routes registered: POST/PATCH/DELETE/GET /api/v1/grades{/{id}} with task relations in GetAll.

**Test with curl (replace YOUR_TOKEN with valid guru JWT):**
```
# GET all grades with tasks
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/grades

# POST new grade
curl -X POST -H "Authorization: Bearer YOUR_TOKEN" -H "Content-Type: application/json" -d '{"task_id":"uuid","tanggal":"2024-01-01","keterangan":"Good","grades":"A"}' http://localhost:8080/api/v1/grades

# PATCH update grade
curl -X PATCH -H "Authorization: Bearer YOUR_TOKEN" -H "Content-Type: application/json" -d '{"keterangan":"Excellent"}' http://localhost:8080/api/v1/grades/{id}

# DELETE grade
curl -X DELETE -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/grades/{id}
```

Run `go run main.go` to start server. All SQL queries use PostgreSQL placeholders ($1 etc.). Dynamic partial updates in PATCH handle optional fields with tx.


## Completed:

