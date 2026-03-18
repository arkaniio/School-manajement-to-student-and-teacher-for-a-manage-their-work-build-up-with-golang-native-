CREATE TABLE students_grades (
    id          UUID NOT NULL PRIMARY KEY DEFAULT
                gen_random_uuid()
    task_id     UUID NOT NULL REFRERENCES tasks(id)
    tanggal     DATE NOT NULL,
    keterangan  VARCHAR(255) NOT NULL,
    grades      INTEGER NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);