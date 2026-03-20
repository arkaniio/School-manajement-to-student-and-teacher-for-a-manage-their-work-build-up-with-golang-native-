CREATE TABLE absensis (
    id                      UUID NOT NULL PRIMARY KEY DEFAULT
                            gen_random_uuid(),
    name_lengkap            VARCHAR(255) NOT NULL,
    kelas                   VACHAR(255) NOT NULL,
    jurusan                 VARCHAR(255) NOT NULL,
    hari                    VARCHAR(255) NOT NULL,
    tanggal                 VARCHAR(255) NOT NULL,
    status                  VARCHAR(50),
    keterangan              VARCHAR(50) NOT NULL,
    created_at              TIMESTAMP NOT NULL,
    updated_at              TIMESTAMP NOT NULL,
    keterangan_tidak_hadir  VARCHAR(255),
    keterangan_dispen       VARCHAR(255),
    file_dispen             VARCHAR(255),
    student_id              UUID NOT NULL REFERENCES students(id)
);