CREATE TABLE tasks (
    id          UUID NOT NULL PRIMARY KEY DEFAULT 
                gen_random_uuid(),
    name_task   VARCHAR(100) NOT NULL,
    file_task   TEXT NOT NULL,
    date_task   TIMESTAMP NOT NULL,
    student_id  UUID NOT NULL REFERENCES students(id),
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
    mapel_task  VARCHAR(255) NOT NULL,
);