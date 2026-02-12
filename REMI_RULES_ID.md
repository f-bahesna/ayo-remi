# Peraturan Permainan Remi (Remi Rules)

## Pengaturan Permainan (Setup)
- **Pemain**: 4 orang.
- **Kartu**: Dek standar 52 kartu + Joker.
- **Pembagian Kartu**:
  - Pemain "Master" (pembuat room) mendapatkan **8 kartu** dan jalan pertama.
  - Pemain lainnya mendapatkan **7 kartu**.
- **Tujuan**: Mengosongkan tangan dengan membentuk set (kombinasi kartu) yang valid dan membuang kartu terakhir, atau memiliki poin terendah saat dek habis.

## Jalannya Permainan (Gameplay)

Permainan berlangsung searah jarum jam. Setiap giliran terdiri dari 3 fase:

### 1. Fase Ambil (Draw Phase)
Pemain harus mengambil kartu dari salah satu sumber:
- **Dek (Tumpukan Tertutup)**: Ambil 1 kartu teratas.
- **Pile (Tumpukan Buangan)** — hanya dari **3 kartu terakhir**:
  - **Ambil biasa**: Klik kartu di 3 teratas pile → kartu tersebut dan semua kartu di atasnya diambil ke tangan. Semua kartu yang diambil **harus** bisa membentuk set dengan kartu di tangan.
  - **Ambil spesifik (Pile Pick)**: Pilih 2+ kartu dari tangan terlebih dahulu, lalu klik 1 kartu dari 3 terakhir pile yang melengkapi set. Set langsung dimainkan ke meja.

### 2. Fase Main (Play Phase)
Setelah mengambil kartu, pemain dapat (opsional):
- Menurunkan set kartu yang valid dari tangan ke meja (Open Card).
- Set yang sudah turun tidak bisa ditarik kembali ke tangan.

### 3. Fase Buang (Discard Phase)
- Pemain **HARUS** membuang 1 kartu dari tangan ke tumpukan buangan (Pile) untuk mengakhiri giliran.
- Giliran berpindah ke pemain berikutnya.

---

## Kombinasi Kartu Valid (Valid Sets)

Ada dua jenis kombinasi yang diperbolehkan (minimal 3 kartu):

### A. Seri / Run (Urutan)
Kartu dengan **bunga (suit) yang sama** dan **angka berurutan**.
- Contoh: `4♠ - 5♠ - 6♠`
- **Ace (As)** bisa menjadi kartu rendah (`A-2-3`) atau tinggi (`Q-K-A`).
- **Tidak Boleh Putar Balik**: `K-A-2` tidak valid.

### B. Group / Set (Angka Sama)
Kartu dengan **angka/rank yang sama** tapi **bunga berbeda**.
- Contoh: `7♠ - 7♥ - 7♣`
- **Syarat Ketat**: Bunga harus unik dalam satu set (tidak boleh ada dua `7♠` dalam set yang sama).

### Joker (Wildcard)
- Kartu Joker bisa menggantikan kartu apa saja untuk melengkapi set.

---

## Cara Menang (Winning Conditions)

### Nutup (Declare Win)
Pemain bisa menekan tombol **"Nutup"** jika:
1. Seluruh kartu di tangan sudah menjadi set yang valid (habis atau sisa 1 untuk dibuang).
2. Atau, pemain telah menurunkan semua set ke meja dan menyisakan 1 kartu terakhir untuk dibuang.

### Dek Habis (Game Over)
Jika tumpukan kartu (Dek) habis dan tidak ada yang menang, permainan berakhir seri atau dihitung berdasarkan poin (lihat Scoring).

---

## Perhitungan Poin (Scoring)

Jika ada pemenang, atau permainan berakhir:
- **Pemenang**: 0 poin.
- **Pemain Kalah**: Dihitung dari total nilai kartu yang tersisa di tangan.

**Nilai Kartu:**
- **Joker**: -250 Poin (Mengurangi skor, sangat bagus jika dipegang saat kalah? *Sesuai implementasi saat ini*).
- **As (Ace)**: 15 Poin.
- **Raja/Gambar (J, Q, K)**: 10 Poin.
- **Angka (2-10)**: 5 Poin.

> **Catatan Strategi**: Usahakan membuang kartu bernilai tinggi (As, J, Q, K) jika tidak bisa dijadikan set, untuk meminimalisir poin jika lawan menang.
