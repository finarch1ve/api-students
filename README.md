# API Students

REST API untuk mengelola data mahasiswa. Dibuat dengan Go dan Fiber v2.

## Cara Menjalankan
```bash
go mod tidy
go run .
```
Server berjalan di `http://localhost:3000`

## Base URL
`http://localhost:3000/api/v1`

## Kontrak API

| Metode | Endpoint | Parameter | Contoh Body | Status | Contoh Respons |
|---|---|---|---|---|---|
| GET | /students | `page`, `limit`, `search`, `sort`, `order`, `is_active` | - | 200 | `{"success":true,"data":[...],"meta":{...}}` |
| GET | /students/:id | - | - | 200 / 404 | `{"success":true,"data":{...}}` |
| POST | /students | - | `{"nim":"434241001","name":"Budi","grade":85}` | 201 / 409 / 422 / 415 | `{"success":true,"data":{...}}` |
| PUT | /students/:id | - | `{"nim":"434241001","name":"Budi Baru","grade":95,"is_active":false}` | 200 / 404 / 422 | `{"success":true,"data":{...}}` |
| PATCH | /students/:id | - | `{"is_active":true}` | 200 / 404 / 422 | `{"success":true,"data":{...}}` |
| DELETE | /students/:id | - | - | 204 / 404 | (tanpa body) |

## Aturan Query String
- `limit` maksimal 50 (data mahasiswa relatif sedikit per halaman, biar tetap enak dibaca)
- `sort` hanya menerima: id, name, nim, grade, created_at (daftar putih, mencegah field sembarangan)
- `search` mencari berdasarkan nama, tidak case-sensitive