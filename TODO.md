# Attendance Statistics Feature Implementation

## Steps:

- [x] **Step 1:** Update `types/absensi.go` - Add AbsensiStats, AbsensiWithStudent structs and extend AbsensiStore interface
- [x] **Step 2:** Create `service/absensis/absensi_service.go` - Business logic layer
- [x] **Step 3:** Update `service/absensis/absensi_repository.go` - Implement stats and JOIN repo methods
- [x] **Step 4:** Update `service/absensis/absensi_handler.go` - Add service dependency and new handlers (GetWeeklyStats_Bp, GetMonthlyStats_Bp, GetAllAbsensiWithStudents_Bp)
- [ ] **Step 5:** Update `cmd/api/router_api.go` - Register new routes with auth middleware
- [ ] **Step 6:** Test compilation and endpoints

**Completed Steps:** 0/6

**Notes:**
- Status mapping: Hadir=&#34;accepted&#34;, Tidak hadir=&#34;not accepted&#34;, Izin=&#34;permissions&#34;
- tanggal: VARCHAR (YYYY-MM-DD format for filters)
- Auth: TokenIdMiddleware + role=&#34;siswa&#34; (matching CRUD)
- Routes: /api/v1/absensi/stats/weekly (GET), /api/v1/absensi/stats/monthly (GET), /api/v1/absensi/all-with-students (GET)

