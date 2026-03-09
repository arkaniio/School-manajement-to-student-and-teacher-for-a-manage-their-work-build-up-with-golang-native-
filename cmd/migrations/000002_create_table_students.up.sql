CREATE TABLE students (
    id              UUID NOT NULL PRIMARY KEY DEFAULT
                    gen_random_uuid(),
    full_name       VARCHAR(100) NOT NULL,
    kelas           VARCHAR(50) NOT NULL,
    jurusan         VARCHAR(100) NOT NULL,
    absen           INTEGER NOT NULL,
    student_profile TEXT NOT NULL,
    wali_kelas      VARCHAR(100) NOT NULL,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL
);