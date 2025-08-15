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
