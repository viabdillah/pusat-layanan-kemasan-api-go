Ini akan fokus pada API, struktur Go, dan cara menjalankan server.

Instruksi:

Buka proyek backend-go Anda di VS Code.

Buat file baru di direktori root backend-go (sejajar dengan go.mod).

Beri nama file tersebut README.md.

Salin (copy) seluruh teks di bawah ini dan tempelkan (paste) ke dalam file README.md Anda.

Markdown

# Panel Admin - Pusat Layanan Kemasan (Backend API)

Ini adalah backend API untuk aplikasi web Panel Admin Pusat Layanan Kemasan. Dibangun menggunakan Golang dengan framework Gin untuk menyediakan API yang cepat, andal, dan aman untuk mengelola seluruh alur kerja bisnis, mulai dari pesanan hingga laporan.

## Fitur Utama

- **Otentikasi & Otorisasi Berbasis Peran:** Sistem login aman menggunakan JWT dengan peran yang berbeda (Admin, Manajer, Kasir, Desainer, Operator).
- **Manajemen Alur Kerja Pesanan:** API lengkap untuk siklus hidup pesanan: pembuatan, proses desain, proses produksi, hingga pembayaran.
- **Manajemen Pelanggan:** Endpoint CRUD untuk mengelola data pelanggan.
- **API Laporan Dinamis:** Endpoint agregasi data untuk laporan penjualan dengan filter waktu (harian, mingguan, bulanan, tahunan).
- **Ekspor ke Excel:** Kemampuan untuk menghasilkan laporan detail dalam format `.xlsx`.

## Tumpukan Teknologi (Tech Stack)

- **Bahasa:** Golang (Go)
- **Framework Web:** Gin
- **Database:** MongoDB (menggunakan Driver resmi Go)
- **Otentikasi:** JSON Web Tokens (JWT)
- **Password Hashing:** Bcrypt
- **Manajemen Dependensi:** Go Modules

## Instalasi & Menjalankan Lokal

Untuk menjalankan proyek ini di lingkungan lokal, ikuti langkah-langkah berikut:

**1. Prasyarat:**
   - [Go](https://go.dev/dl/) versi 1.20+
   - [MongoDB](https://www.mongodb.com/try/download/community)
   - [Git](https://git-scm.com/downloads)

**2. Clone Repositori:**
   ```bash
   git clone [https://github.com/](https://github.com/)[NamaAnda]/pusat-layanan-kemasan-api-go.git
   cd pusat-layanan-kemasan-api-go
3. Siapkan Environment Variable:

Buat file .env di direktori root.

Salin isi dari .env.example (jika ada) atau gunakan format berikut:

Cuplikan kode

PORT=8080
MONGO_URI=mongodb://127.0.0.1:27017
JWT_SECRET=kunci_rahasia_anda_yang_sangat_aman
4. Instal Dependensi:

Bash

go mod tidy
5. Jalankan Server:

Bash

go run cmd/server/main.go
Server akan berjalan di http://localhost:8080.

Tentang Pengembang
Hai! Saya [Nama Lengkap Anda], seorang [Jabatan Anda, misal: Full-Stack Developer] dengan passion untuk membangun aplikasi web yang efisien dan bermanfaat. Proyek ini adalah salah satu eksplorasi saya dalam menggunakan Golang untuk backend dan Vue.js untuk frontend.

LinkedIn: https://linkedin.com/in/nama-anda

GitHub: https://github.com/NamaAnda

Portfolio/Website: (Opsional)


---

### **2. README untuk Frontend (`frontend-dashboard/README.md`)**

Ini akan fokus pada antarmuka pengguna, fitur-fitur interaktif, dan cara menjalankan aplikasi Vue.

**Instruksi:**
1.  Buka proyek `frontend-dashboard` Anda.
2.  Buat file baru di direktori **root** `frontend-dashboard` (sejajar dengan `package.json`).
3.  Beri nama file tersebut `README.md`.
4.  Salin dan tempel seluruh teks di bawah ini ke dalam file tersebut.

```markdown
# Panel Admin - Pusat Layanan Kemasan (Frontend)

Ini adalah antarmuka pengguna (UI) untuk aplikasi web Panel Admin Pusat Layanan Kemasan. Dibangun menggunakan Vue.js 3 dan Vite, aplikasi ini menyediakan dasbor yang reaktif, modern, dan sadar-peran (*role-aware*) untuk berinteraksi dengan backend API.

## Fitur Utama

- **Antarmuka Berbasis Peran:** Tampilan dan navigasi berubah secara dinamis sesuai dengan peran pengguna yang login (Admin, Manajer, Kasir, dll.).
- **Manajemen State Terpusat:** Menggunakan Pinia untuk mengelola state otentikasi dan data aplikasi secara konsisten.
- **Routing yang Aman:** Menggunakan Vue Router dengan *Navigation Guards* untuk melindungi halaman-halaman sensitif.
- **Formulir Dinamis & Interaktif:** Formulir pembuatan pesanan yang canggih dengan pencarian pelanggan *real-time* dan kemampuan menambah/menghapus item pesanan.
- **Tampilan Data:** Tabel dan kartu statistik yang menampilkan data dari API, lengkap dengan status *loading* dan *error*.
- **Desain Responsif:** Didesain dengan Tailwind CSS agar terlihat bagus di semua ukuran perangkat, dari ponsel hingga desktop.

## Tumpukan Teknologi (Tech Stack)

- **Framework:** Vue.js 3 (dengan Composition API)
- **Build Tool:** Vite
- **Routing:** Vue Router
- **Manajemen State:** Pinia
- **Styling:** Tailwind CSS
- **HTTP Client:** Axios

## Instalasi & Menjalankan Lokal

Untuk menjalankan proyek ini di lingkungan lokal, ikuti langkah-langkah berikut:

**1. Prasyarat:**
   - [Node.js](https://nodejs.org/) versi 18+
   - [Git](https://git-scm.com/downloads)
   - Backend API harus sudah berjalan.

**2. Clone Repositori:**
   ```bash
   git clone [https://github.com/](https://github.com/)[NamaAnda]/pusat-layanan-kemasan-frontend-vue.git
   cd pusat-layanan-kemasan-frontend-vue
3. Siapkan Environment Variable:

Buat file .env di direktori root.

Tambahkan URL base dari backend API Anda:

Cuplikan kode

VITE_API_BASE_URL=http://localhost:8080/api
4. Instal Dependensi:

Bash

npm install
5. Jalankan Server Development:

Bash

npm run dev
Aplikasi akan berjalan di http://localhost:5173.

Tentang Pengembang
Hai! Saya [Nama Lengkap Anda], seorang [Jabatan Anda, misal: Full-Stack Developer] dengan passion untuk membangun aplikasi web yang efisien dan bermanfaat. Proyek ini adalah salah satu eksplorasi saya dalam menggunakan Golang untuk backend dan Vue.js untuk frontend.

LinkedIn: https://linkedin.com/in/nama-anda

GitHub: https://github.com/NamaAnda

Portfolio/Website: (Opsional)


### **Langkah Terakhir: Unggah README ke GitHub**

Jangan lupa untuk mengunggah file-file baru ini ke repositori GitHub Anda.

1.  Buka terminal di direktori `backend-go`:
    ```bash
    git add README.md
    git commit -m "docs: Add professional README file for backend"
    git push
    ```
2.  Buka terminal di direktori `frontend-dashboard`:
    ```bash
    git add README.md
    git commit -m "docs: Add professional README file for frontend"
    git push
    ```

Sekarang, proyek Anda di GitHub terlihat jauh lebih profesional!

Ketik **lanjutkan** untuk mulai mengerjakan **Langkah 30: Membuat halaman "Detail Pesanan"**.
