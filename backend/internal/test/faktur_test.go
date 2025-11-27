package test

import (
	"api-distributor/helper"
	"api-distributor/internal/dto"
	"fmt"
	"sort"
	"testing"
)

func TestCetakFaktur(t *testing.T) {
	fmt.Println("Test cetak faktur")
	// dummy data detil penjualan
	detilJualDTO := []dto.PenjualanDetailResponse{
		{
			BarangID:       1,
			BarangKode:     "JT-001",
			BarangArtikel:  "Lembata",
			BarangWarna:    "Hitam",
			BarangUkuran:   "38-41",
			BarangKategori: "D",
			Size:           41,
			Qty:            3,
			Harga:          140000,
			Subtotal:       420000,
		},
		{
			BarangID:       1,
			BarangKode:     "JT-001",
			BarangArtikel:  "Lembata",
			BarangWarna:    "Hitam",
			BarangUkuran:   "38-41",
			BarangKategori: "D",
			Size:           40,
			Qty:            2,
			Harga:          140000,
			Subtotal:       280000,
		},
		{
			BarangID:       1,
			BarangKode:     "JT-001",
			BarangArtikel:  "Lembata",
			BarangWarna:    "Hitam",
			BarangUkuran:   "38-41",
			BarangKategori: "D",
			Size:           39,
			Qty:            5,
			Harga:          140000,
			Subtotal:       700000,
		},
		{
			BarangID:       2,
			BarangKode:     "JT-002",
			BarangArtikel:  "Lembata",
			BarangWarna:    "Putih",
			BarangUkuran:   "38-41",
			BarangKategori: "D",
			Size:           39,
			Qty:            4,
			Harga:          140000,
			Subtotal:       560000,
		},
		{
			BarangID:       3,
			BarangKode:     "SB-190",
			BarangArtikel:  "Chris Tomlin",
			BarangWarna:    "Putih",
			BarangUkuran:   "33-36",
			BarangKategori: "A",
			Size:           33,
			Qty:            10,
			Harga:          150000,
			Subtotal:       1500000,
		},
		{
			BarangID:       2,
			BarangKode:     "SB-190",
			BarangArtikel:  "Chris Tomlin",
			BarangWarna:    "Putih",
			BarangUkuran:   "33-36",
			BarangKategori: "A",
			Size:           34,
			Qty:            1,
			Harga:          150000,
			Subtotal:       150000,
		},
	}

	// debug (optional)
	for _, detil := range detilJualDTO {
		fmt.Println(detil)
	}

	// 1. Ambil semua size unik
	sizeSet := map[int]bool{}
	for _, d := range detilJualDTO {
		sizeSet[d.Size] = true
	}

	sizes := make([]int, 0, len(sizeSet))
	for s := range sizeSet {
		sizes = append(sizes, s)
	}
	sort.Ints(sizes)

	// 2. Buat pivot: "Kode + warna" -> size -> qty dan map harga
	hargaMap := map[string]int{}
	pivot := map[string]map[int]int{}

	for _, d := range detilJualDTO {
		// key pivot: kode + warna
		key := fmt.Sprintf("%-7s %-13s %-10s %-1s", d.BarangKode, d.BarangArtikel, d.BarangWarna, d.BarangKategori)

		if _, ok := pivot[key]; !ok {
			pivot[key] = map[int]int{}
		}
		// key = barang kode + warna
		pivot[key][d.Size] += d.Qty
		// simpan harga (ambil dari dto)
		hargaMap[key] = d.Harga
	}

	fmt.Println("=================================================================================")
	// 3. Print pivot table
	fmt.Printf("%-7s %-13s %-10s %-1s", "Kode", "Artikel", "Warna", "K")
	for _, s := range sizes {
		fmt.Printf("%5d", s)
	}

	// header kolom total
	fmt.Printf("%5s", "Qty")

	// header kolom harga
	fmt.Printf("%8s", "Harga")

	// header kolom total
	fmt.Printf("%9s", "Total")

	fmt.Println("\n=================================================================================")

	for key, row := range pivot {
		fmt.Printf("%-20s", key)
		// initial value total per row
		total := 0

		for _, s := range sizes {
			val := row[s]
			if val == 0 {
				fmt.Printf("%5s", " ") // kosong jika 0
			} else {
				fmt.Printf("%5d", val) // qty sebenarnya
				total += val
			}
		}
		// tampilkan total di kolom terakhir
		fmt.Printf("%5d", total)
		fmt.Printf("%8d", hargaMap[key])       // harga dari dto
		fmt.Printf("%9d", total*hargaMap[key]) // hitung total harga * qty

		fmt.Println()
	}
	fmt.Println("=================================================================================")
}

func TestCetakFakturMultiHeader(t *testing.T) {
	data := []dto.PenjualanDetailResponse{
		{
			BarangID:       1,
			BarangKode:     "VP001",
			BarangArtikel:  "BLAZE",
			BarangWarna:    "ALL BLACK",
			BarangUkuran:   "37-44",
			BarangKategori: "D",
			Size:           41,
			Qty:            5,
			Harga:          139800,
			Subtotal:       699000,
		},
		{
			BarangID:       1,
			BarangKode:     "VP001",
			BarangArtikel:  "BLAZE",
			BarangWarna:    "ALL BLACK",
			BarangUkuran:   "37-44",
			BarangKategori: "D",
			Size:           40,
			Qty:            4,
			Harga:          139800,
			Subtotal:       559200,
		},
		{
			BarangID:       258,
			BarangKode:     "CAV001",
			BarangArtikel:  "LEMBATA",
			BarangWarna:    "HTM/HTM",
			BarangUkuran:   "38-43",
			BarangKategori: "D",
			Size:           39,
			Qty:            3,
			Harga:          140000,
			Subtotal:       420000,
		},
		{
			BarangID:       258,
			BarangKode:     "CAV001",
			BarangArtikel:  "LEMBATA",
			BarangWarna:    "HTM/HTM",
			BarangUkuran:   "38-43",
			BarangKategori: "D",
			Size:           40,
			Qty:            1,
			Harga:          140000,
			Subtotal:       140000,
		},
		{
			BarangID:       258,
			BarangKode:     "CAV001",
			BarangArtikel:  "LEMBATA",
			BarangWarna:    "HTM/HTM",
			BarangUkuran:   "38-43",
			BarangKategori: "D",
			Size:           42,
			Qty:            6,
			Harga:          140000,
			Subtotal:       840000,
		},
		{
			BarangID:       258,
			BarangKode:     "CAV001",
			BarangArtikel:  "LEMBATA",
			BarangWarna:    "HTM/HTM",
			BarangUkuran:   "38-43",
			BarangKategori: "D",
			Size:           43,
			Qty:            8,
			Harga:          140000,
			Subtotal:       1120000,
		},
	}

	// dummy data helper.Produk
	// data := []helper.Produk{
	// 	{Nama: "Sepatu A", Kategori: "D", Size: 40, Stok: 13, Harga: 10000},
	// 	{Nama: "Sepatu A", Kategori: "D", Size: 41, Stok: 41, Harga: 10000},
	// 	{Nama: "Sepatu A", Kategori: "D", Size: 42, Stok: 42, Harga: 10000},
	// 	{Nama: "Sepatu B", Kategori: "A", Size: 33, Stok: 25, Harga: 20000},
	// 	{Nama: "Sepatu B", Kategori: "A", Size: 34, Stok: 34, Harga: 21000},
	// 	{Nama: "Sepatu C", Kategori: "D", Size: 42, Stok: 10, Harga: 30000},
	// 	{Nama: "Sepatu D", Kategori: "A", Size: 35, Stok: 7, Harga: 40000},
	// 	{Nama: "Sepatu E", Kategori: "K", Size: 28, Stok: 28, Harga: 50000},
	// 	{Nama: "Sepatu F", Kategori: "K", Size: 29, Stok: 29, Harga: 60000},
	// 	{Nama: "Sepatu F", Kategori: "K", Size: 26, Stok: 26, Harga: 159000},
	// }

	kategoriSizeRange := map[string][]int{
		"K": {26, 27, 28, 29, 30, 31},
		"A": {32, 33, 34, 35, 36, 37},
		"D": {36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
	}

	kategoriLabel := map[string]string{"K": "|K|", "A": "|A|", "D": "|D|"}

	totalSizes := len(kategoriSizeRange["K"]) + len(kategoriSizeRange["A"]) + len(kategoriSizeRange["D"])

	pivot := helper.BuildPivot(data, kategoriSizeRange)

	output := helper.PrintHeader(kategoriSizeRange, kategoriLabel, totalSizes)
	output += helper.PrintData(pivot, kategoriSizeRange)
	output += helper.PrintSumarry()

	err := helper.WriteToFile("faktur.txt", output)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Berhasil membuat file faktur.txt")
	}
}
