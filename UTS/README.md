# SIAKAD Mini API

RESTful API back end layanan akademik sederhana (mahasiswa, mata kuliah, KRS).
Teknologi: Go, Fiber v2, GORM, PostgreSQL, JWT, bcrypt.

## Cara menjalankan
1. Buat database PostgreSQL bernama `db_siakad`.
2. Salin `.env.example` menjadi `.env`, lalu isi nilainya (koneksi database, `JWT_SECRET`, `JWT_EXPIRES_HOURS`).
3. Jalankan: go mod tidy go run .
   Migration dan seeder berjalan otomatis saat aplikasi dimulai. Server aktif di `http://localhost:3000`.

## Akun seeder
| Role | Email | Password |
|---|---|---|
| admin | admin@siakad.test | admin12345 |
| mahasiswa | mhs01@siakad.test ... mhs20@siakad.test | NIM masing-masing |

Seeder membuat 1 admin, 20 mahasiswa, dan 10 mata kuliah. Contoh: `mhs01@siakad.test` dengan password `434241000001`.

## Endpoint
| No | Method | Endpoint | Akses |
|---|---|---|---|
| 1 | POST | /api/v1/auth/login | Publik |
| 2 | GET | /api/v1/auth/me | Semua role |
| 3 | GET | /api/v1/students | Admin |
| 4 | POST | /api/v1/students | Admin |
| 5 | GET | /api/v1/students/{id} | Admin, mahasiswa (data sendiri) |
| 6 | PUT | /api/v1/students/{id} | Admin |
| 7 | DELETE | /api/v1/students/{id} | Admin |
| 8 | GET | /api/v1/courses | Semua role |
| 9 | POST | /api/v1/enrollments | Mahasiswa |
| 10 | DELETE | /api/v1/enrollments/{id} | Mahasiswa (milik sendiri) |

Semua endpoint selain login memerlukan header `Authorization: Bearer <token>`.

## Laporan
Laporan pengujian lengkap ada pada file `Laporan-UTS.pdf` di folder ini.