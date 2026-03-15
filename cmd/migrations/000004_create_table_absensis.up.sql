CREATE TABLE absensis (
    id                  UUID NOT NULL PRIMARY KEY DEFAULT
                        gen_random_uuid(),
    name_lengkap        VARCHAR(255) NOT NULL,
    kelas               VARCHAR(255) NOT NULL,
    jurusan             VARCHAR(255) NOT NULL,
    hari                VARCHAR(255) NOT NULL,
    tanggal             VARCHAR(255) NOT NULL,
    status              VARCHAR(50) NOT NULL,
    keterangan_izin     VARCHAR(50),
    created_at          TIMESTAMP NOT NULL,
    updated_at          TIMESTAMP NOT NULL
);