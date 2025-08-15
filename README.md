# 🚀 Pusat Layanan Kemasan - Backend API (Versi Go)

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=for-the-badge&logo=go)
![Gin](https://img.shields.io/badge/Gin-v1.9-007EC6?style=for-the-badge&logo=gin)
![MongoDB](https://img.shields.io/badge/MongoDB-4.4+-47A248?style=for-the-badge&logo=mongodb)
![License](https://img.shields.io/github/license/mashape/apistatus.svg?style=for-the-badge)

Selamat datang di repositori backend untuk aplikasi **Pusat Layanan Kemasan**. Proyek ini adalah implementasi API yang tangguh dan berperforma tinggi yang dibangun menggunakan Golang dan framework Gin, dirancang untuk mengelola alur kerja bisnis yang kompleks dari awal hingga akhir.

---

## ✨ Fitur Utama

-   🔐 **Otentikasi & Otorisasi Berbasis Peran:** Sistem login aman menggunakan JWT dan *middleware* untuk membatasi akses berdasarkan peran pengguna (Admin, Manajer, Kasir, Desainer, Operator).
-   🔄 **Alur Kerja Pesanan Dinamis:** Mengelola siklus hidup pesanan dari pembuatan oleh Kasir, proses desain oleh Desainer, produksi oleh Operator, hingga finalisasi pembayaran.
-   📈 **Manajemen Data Terpusat:** API CRUD untuk mengelola data Pelanggan dan Pesanan multi-item.
-   📊 **Pelaporan Dinamis:** Endpoint agregasi data yang kuat untuk menghasilkan laporan penjualan berdasarkan rentang waktu (harian, mingguan, bulanan, tahunan).
-   📄 **Ekspor ke Excel:** Kemampuan untuk mengunduh laporan penjualan yang telah difilter langsung dalam format `.xlsx`.
-   🏢 **Struktur Proyek Profesional:** Mengikuti konvensi proyek Go modern dengan struktur direktori yang jelas (`cmd`, `internal`).

---

## 🛠️ Teknologi yang Digunakan

-   **Bahasa:** [Golang](https://go.dev/)
-   **Web Framework:** [Gin](https://gin-gonic.com/)
-   **Database:** [MongoDB](https://www.mongodb.com/) (dengan [Official Go Driver](https://github.com/mongodb/mongo-go-driver))
-   **Password Hashing:** [Bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt)
-   **Token JWT:** [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
-   **Konfigurasi:** [godotenv](https://github.com/joho/godotenv)
-   **Excel Generation:** [Excelize](https://github.com/xuri/excelize)

---

## ⚙️ Cara Menjalankan Secara Lokal

Untuk menjalankan proyek ini di lingkungan lokal Anda, ikuti langkah-langkah berikut:

1.  **Clone Repositori**
    ```bash
    git clone [https://github.com/NAMA_ANDA/pusat-layanan-kemasan-api-go.git](https://github.com/NAMA_ANDA/pusat-layanan-kemasan-api-go.git)
    cd pusat-layanan-kemasan-api-go
    ```

2.  **Siapkan Variabel Lingkungan**
    Buat file bernama `.env` di direktori utama dan isi dengan konfigurasi berikut:
    ```.env
    # Port untuk server backend
    PORT=8080

    # Alamat koneksi ke MongoDB lokal Anda
    MONGO_URI=mongodb://127.0.0.1:27017

    # Kunci rahasia untuk JSON Web Token (buat yang unik dan acak)
    JWT_SECRET=kunci_rahasia_golang_yang_sangat_aman
    ```

3.  **Instal Dependensi**
    Go akan secara otomatis menangani ini saat Anda menjalankan aplikasi, tetapi Anda bisa melakukannya secara manual dengan:
    ```bash
    go mod tidy
    ```

4.  **Jalankan Server**
    ```bash
    go run cmd/server/main.go
    ```
    Server akan berjalan di `http://localhost:8080`.

---

## 📚 Dokumentasi API

Berikut adalah ringkasan dari *endpoint* utama yang tersedia.

| Method | Endpoint                      | Deskripsi                                        | Akses               |
| :----- | :---------------------------- | :----------------------------------------------- | :------------------ |
| `POST` | `/api/users/register`         | Mendaftarkan user baru (default role: kasir)     | Publik              |
| `POST` | `/api/users/login`            | Login pengguna dan mendapatkan token JWT           | Publik              |
| `GET`  | `/api/users/profile`          | Mendapatkan profil pengguna yang sedang login      | Terotentikasi       |
| `GET`  | `/api/users`                  | Mendapatkan semua data pengguna                    | Admin               |
| `POST` | `/api/customers`              | Membuat data pelanggan baru                      | Kasir, Admin        |
| `GET`  | `/api/customers`              | Mencari/mendapatkan semua pelanggan              | Kasir, Admin        |
| `POST` | `/api/orders`                 | Membuat pesanan baru dengan multi-item           | Kasir, Admin        |
| `GET`  | `/api/orders/queue/designer`  | Mendapatkan antrian kerja untuk Designer         | Designer, Admin     |
| `GET`  | `/api/orders/queue/operator`  | Mendapatkan antrian kerja untuk Operator         | Operator, Admin     |
| `GET`  | `/api/orders/queue/kasir`     | Mendapatkan antrian pesanan siap bayar           | Kasir, Admin        |
| `PATCH`| `/api/orders/:id/status`      | Memperbarui status sebuah pesanan                | Designer, Operator, Kasir, Admin |
| `GET`  | `/api/orders/monitoring`      | Mendapatkan semua pesanan yang sedang diproses     | Manajer, Admin      |
| `GET`  | `/api/reports/sales-summary`  | Mendapatkan ringkasan penjualan (JSON)           | Manajer, Admin      |
| `GET`  | `/api/reports/sales-summary/export`| Mengunduh laporan penjualan (Excel)        | Manajer, Admin      |

---

## 🚀 Rencana Pengembangan

-   [ ] **Langkah 30:** Membuat Halaman Detail Pesanan di Frontend.
-   [ ] **Langkah 31:** Menambahkan Grafik & Bagan pada Halaman Laporan.
-   [ ] **Langkah 32:** Membangun Sistem Notifikasi Dasar.
-   [ ] Implementasi UI penuh untuk Manajemen Pengguna (CRUD).
-   [ ] Fungsionalitas unggah file untuk Desain.
-   [ ] Deployment ke server produksi (misalnya: Render, Fly.io).

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah [MIT License](https://choosealicense.com/licenses/mit/).
MIT License

Copyright (c) 2025 viabdillah

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
