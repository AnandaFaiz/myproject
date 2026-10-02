# 📰 Portal Berita

Aplikasi portal berita fullstack dengan arsitektur modern:

- **Frontend:** Nuxt 4 + Vue 3 + TypeScript
- **Backend:** Go + Gin
- **Database:** PostgreSQL

Fitur utama: CRUD berita, autentikasi JWT, komentar, upload gambar, dan dashboard admin.

---

## 📁 Struktur Proyek

```
portal-berita/
├── newsapp/                  # Frontend (Nuxt 4 + Vue 3)
│   ├── app/
│   │   ├── pages/            # Halaman (beranda, detail, admin, dll.)
│   │   ├── composables/      # useAuth, dll.
│   │   └── middleware/       # Proteksi route
│   ├── public/               # Static assets
│   ├── nuxt.config.ts
│   └── package.json
│
├── portal-berita-api/        # Backend (Go + Gin)
│   ├── config/               # Koneksi database
│   ├── handlers/             # Handler HTTP
│   ├── middleware/           # Auth, Admin
│   ├── models/               # Struct model
│   ├── utils/                # Helper (JWT, slug)
│   ├── migrations/           # File migrasi database (goose)
│   ├── public/uploads/       # Upload gambar
│   ├── main.go
│   └── go.mod
│
└── README.md
```

---

## ✅ Prasyarat

Pastikan sudah terinstall:

| Tool           | Versi Minimal | Cara Cek          |
| -------------- | :-----------: | ----------------- |
| **Node.js**    |      18+      | `node --version`  |
| **Go**         |     1.21+     | `go version`      |
| **PostgreSQL** |      14+      | `psql --version`  |
| **Goose**      |      3+       | `goose --version` |

### Install yang Belum Ada

- **Node.js** → [nodejs.org](https://nodejs.org)
- **Go** → [go.dev/dl](https://go.dev/dl/)
- **PostgreSQL** → [postgresql.org/download](https://www.postgresql.org/download/)
- **Goose:**
  ```bash
  go install github.com/pressly/goose/v3/cmd/goose@latest
  ```
  Pastikan `$GOPATH/bin` sudah masuk ke `PATH`.

---

## 🚀 Setup

### 1. Clone Repository

```bash
git clone https://github.com/AnandaFaiz/myproject.git
cd myproject
```

### 2. Setup Database

Buka **pgAdmin** atau **psql**, buat database baru:

```sql
CREATE DATABASE portal_berita;
```

### 3. Jalankan Migrasi

Migrasi akan membuat semua tabel (users, berita, komentar) dan menambahkan akun admin default.

Masuk ke folder backend:

```bash
cd portal-berita-api
```

Jalankan migrasi (ganti `PASSWORD_KAMU` dengan password PostgreSQL-mu):

```bash
goose -dir migrations postgres "postgresql://postgres:PASSWORD_KAMU@localhost:5432/portal_berita?sslmode=disable" up
```

**Kalau berhasil**, output-nya akan seperti:

```
OK   202610020426_create_users.sql
OK   202610020504_create_berita.sql
OK   202610020512_create_komentar.sql
OK   202610020550_seed_admin.sql
goose: no migrations to run. current version: 202610020550
```

### 4. Setup Environment Backend

Copy file `.env.example` menjadi `.env`:

```bash
cp .env.example .env
```

Edit file `.env` dan sesuaikan:

```env
DATABASE_URL=postgresql://postgres:PASSWORD_KAMU@localhost:5432/portal_berita?sslmode=disable
JWT_SECRET=ganti-dengan-string-acak-panjang-minimal-32-karakter
PORT=8080
```

**Cara bikin JWT_SECRET acak:** buka [random.org/strings](https://www.random.org/strings/) atau ketik sembarang string panjang (32+ karakter).

### 5. Jalankan Backend

Masih di folder `portal-berita-api`:

```bash
go mod download
go run main.go
```

**Kalau berhasil:**

```
Berhasil koneksi ke PostgreSQL
Server jalan di http://localhost:8080
```

**Jangan tutup terminal ini**, biarkan backend tetap jalan.

### 6. Jalankan Frontend

Buka **terminal baru** (jangan tutup yang jalan backend), lalu:

```bash
cd newsapp
npm install
npm run dev
```

**Kalau berhasil:**

```
Nuxt 4.x.x
Local: http://localhost:3000
```

### 7. Buka Aplikasi

Buka browser:

👉 **http://localhost:3000**

---

## 🔐 Login Admin Default

Setelah migrasi dijalankan, kamu bisa login sebagai admin:

| Field    | Value              |
| -------- | ------------------ |
| Email    | `admin@portal.com` |
| Password | `admin000`         |

⚠️ **Ganti password ini setelah login pertama di production!**

---

## 🌐 Akses Aplikasi

| URL                                | Fungsi                         |
| ---------------------------------- | ------------------------------ |
| http://localhost:3000              | Beranda (publik)               |
| http://localhost:3000/berita/:slug | Detail berita                  |
| http://localhost:3000/login        | Login                          |
| http://localhost:3000/register     | Daftar user baru               |
| http://localhost:3000/profil       | Profil user                    |
| http://localhost:3000/admin/berita | Dashboard admin (khusus admin) |

---

## 🔌 API Endpoints

Backend berjalan di `http://localhost:8080`.

### Berita

| Method | Endpoint                 | Akses  | Deskripsi      |
| ------ | ------------------------ | :----: | -------------- |
| GET    | `/api/berita`            | Publik | Semua berita   |
| GET    | `/api/berita/slug/:slug` | Publik | Berita by slug |
| POST   | `/api/berita`            | Admin  | Tambah berita  |
| PUT    | `/api/berita/:id`        | Admin  | Update berita  |
| DELETE | `/api/berita/:id`        | Admin  | Hapus berita   |

### Autentikasi

| Method | Endpoint             | Akses  | Deskripsi            |
| ------ | -------------------- | :----: | -------------------- |
| POST   | `/api/auth/register` | Publik | Daftar user baru     |
| POST   | `/api/auth/login`    | Publik | Login (return JWT)   |
| POST   | `/api/auth/logout`   | Publik | Logout               |
| GET    | `/api/auth/me`       | Login  | Info user yang login |
| PUT    | `/api/auth/profil`   | Login  | Update profil        |
| PUT    | `/api/auth/password` | Login  | Ganti password       |

### Komentar

| Method | Endpoint                   |     Akses     | Deskripsi           |
| ------ | -------------------------- | :-----------: | ------------------- |
| GET    | `/api/komentar/:berita_id` |    Publik     | Komentar per berita |
| POST   | `/api/komentar`            |     Login     | Tambah komentar     |
| DELETE | `/api/komentar/:id`        | Admin/Pemilik | Hapus komentar      |

### Upload

| Method | Endpoint      | Akses | Deskripsi            |
| ------ | ------------- | :---: | -------------------- |
| POST   | `/api/upload` | Admin | Upload gambar berita |

**Autentikasi:** Kirim JWT di header:

```
Authorization: Bearer <token>
```

---

## 🛠️ Catatan Teknis

### Backend (Go)

- **Framework:** Gin
- **Database Driver:** pgx/v5
- **JWT:** golang-jwt/jwt/v5
- **Password Hash:** bcrypt
- **Migrasi:** goose

### Frontend (Nuxt 4)

- **Framework:** Nuxt 4 + Vue 3 + TypeScript
- **Styling:** CSS murni (scoped)
- **Icons:** Material Design Icons (MDI)
- **Proxy:** `nitro.devProxy` meneruskan `/api` ke backend Go

### Struktur Komunikasi

```
┌─────────────────┐         ┌─────────────────┐
│  Nuxt (port     │  HTTP   │  Go (port       │
│  3000)          │◄───────►│  8080)          │
└─────────────────┘         └────────┬────────┘
                                     │
                                     ▼
                            ┌─────────────────┐
                            │  PostgreSQL     │
                            │  (port 5432)    │
                            └─────────────────┘
```

---

## 🧪 Testing Cepat

Setelah semua jalan, coba:

1. **Buka beranda** → harusnya tampil daftar berita (kosong kalau belum ada)
2. **Login sebagai admin** → masuk dashboard
3. **Tambah berita** → isi form + upload gambar
4. **Lihat beranda** → berita baru muncul
5. **Klik berita** → halaman detail + form komentar
6. **Logout, daftar sebagai user biasa, login, komentar**

---

## 📚 Migrasi Database

Kalau ada perubahan skema database:

### Buat Migration Baru

```bash
cd portal-berita-api
goose -dir migrations create nama_migration sql
```

Edit file yang di-generate di `migrations/`, isi dengan SQL.

### Jalankan Migration

```bash
goose -dir migrations postgres "postgresql://postgres:PASSWORD@localhost:5432/portal_berita?sslmode=disable" up
```

### Rollback

```bash
# Rollback 1 migration terakhir
goose -dir migrations postgres "..." down

# Rollback semua
goose -dir migrations postgres "..." down-to 0
```

### Cek Status

```bash
goose -dir migrations postgres "..." status
```

---

## 🐛 Troubleshooting

### Backend Error "Connection Refused"

PostgreSQL belum jalan. Cek service PostgreSQL, atau restart komputer.

### Backend Error "password authentication failed"

Password di `.env` salah. Cek kembali `DATABASE_URL`.

### Frontend Error "404 /api/..."

Backend Go tidak jalan, atau proxy salah. Cek `nitro.devProxy` di `nuxt.config.ts`.

### Migrasi Error "relation already exists"

Database sudah punya tabel. Reset:

```bash
goose -dir migrations postgres "..." down-to 0
goose -dir migrations postgres "..." up
```

Atau drop database dan buat ulang.

### Upload Gambar Gagal

Pastikan folder `portal-berita-api/public/uploads/` ada dan writable.

---

## 🎯 Status Proyek

**Sudah selesai:**

- ✅ CRUD Berita
- ✅ Autentikasi (register, login, JWT)
- ✅ Role user (admin & user biasa)
- ✅ Update profil & ganti password
- ✅ Upload gambar
- ✅ Komentar
- ✅ Proxy Nuxt ↔ Go

**Dalam pengembangan:**

- ⏳ Bookmark berita
- ⏳ Like / reaksi
- ⏳ Search & filter kategori
- ⏳ Pagination
- ⏳ SEO meta tags
- ⏳ Deployment

---

## 🤝 Kontribusi

Ini proyek untuk review internal. Kalau ada saran atau masukan, silakan buka **Issue** di GitHub.

---

## 📄 Lisensi

MIT

---

## 👤 Kontak

**Nama:** Fais
**Email:** anandafaiz@gmail.com
**GitHub:** [@AnandaFaiz](https://github.com/usernameAnandaFaiz)

---

_Dibuat dengan ❤️ menggunakan Nuxt 4 + Go_
