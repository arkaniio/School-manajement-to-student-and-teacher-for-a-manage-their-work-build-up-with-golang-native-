CREATE TABLE tasks (
    id          UUID NOT NULL PRIMARY KEY DEFAULT 
                gen_random_uuid(),
    name_task   VARCHAR(50) NOT NULL,
    file_task   VARCHAR(50) NOT NULL,
    date_task   DATE NOT NULL,
    student_id  UUID NOT NULL REFERENCES students(id),
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL
);