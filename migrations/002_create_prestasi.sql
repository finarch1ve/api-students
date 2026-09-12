-- tabel prestasi
CREATE TABLE prestasi (
    id SERIAL PRIMARY KEY,
    nama_prestasi VARCHAR(255) NOT NULL,
    id_student INT NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    juara VARCHAR(100) NOT NULL
);

