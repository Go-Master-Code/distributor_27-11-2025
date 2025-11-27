package helper

import (
	"api-distributor/internal/dto"
	"api-distributor/internal/model"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// ======================= STRUCT =======================
type Produk struct {
	Artikel  string
	Warna    string
	Kategori string
	Size     int
	Stok     int
	Harga    int
}

type PivotRow struct {
	Artikel  string
	Warna    string
	Kategori string
	Harga    int
	QtyMap   map[string]map[int]int
}

type SizeMapping struct {
	KatAsal    string
	SizeAsal   int
	KatTujuan  string
	SizeTujuan int
}

// ======================= BUILD PIVOT =======================
func BuildPivot(data []dto.PenjualanDetailResponse, kategoriSizeRange map[string][]int) map[string]*PivotRow {
	pivot := map[string]*PivotRow{}

	for _, p := range data {
		if _, ok := pivot[p.BarangArtikel]; !ok {
			pivot[p.BarangArtikel] = &PivotRow{
				Artikel:  p.BarangArtikel,
				Kategori: p.BarangKategori,
				Harga:    p.Harga,
				QtyMap:   map[string]map[int]int{},
			}
			for kat, sizes := range kategoriSizeRange {
				pivot[p.BarangArtikel].QtyMap[kat] = map[int]int{}
				for _, s := range sizes {
					pivot[p.BarangArtikel].QtyMap[kat][s] = 0
				}
			}
		}
		pivot[p.BarangArtikel].QtyMap[p.BarangKategori][p.Size] = p.Qty
	}

	// jalankan auto size mapping
	for _, row := range pivot {
		for kat, sizes := range row.QtyMap {
			for size, stok := range sizes {
				if stok == 0 {
					continue
				}
				mappings := AutoConvertSize(kat, size)
				for _, m := range mappings {
					row.QtyMap[m.KatTujuan][m.SizeTujuan] = stok
					row.QtyMap[m.KatAsal][m.SizeAsal] = 0
				}
			}
		}
	}

	return pivot
}

// ======================= AUTO SIZE =======================
func AutoConvertSize(kat string, size int) []SizeMapping {
	var res []SizeMapping

	if kat == "D" {
		dest := size - 10
		if dest >= 26 && dest <= 31 {
			res = append(res, SizeMapping{"D", size, "K", dest})
		}
		if dest >= 32 && dest <= 37 {
			res = append(res, SizeMapping{"D", size, "A", dest})
		}
	}

	if kat == "A" {
		dest := size - 6
		if dest >= 26 && dest <= 31 {
			res = append(res, SizeMapping{"A", size, "K", dest})
		}
	}

	return res
}

// ======================= PRINT HEADER =======================
func PrintHeader(kategoriSizeRange map[string][]int, kategoriLabel map[string]string, totalSizes int) string {
	var b strings.Builder

	// dummy data penjualan
	penjualan := model.Penjualan{
		ID:            111,
		NoFaktur:      "JT-001",
		TglPenjualan:  time.Now(),
		TglJatuhTempo: time.Now(),
		TokoID:        1,
		Toko: model.Toko{
			Nama:   "AGRIS SPORT",
			Alamat: "JL. KH Zainal Alim No. 31. Kemayoran.",
			KotaID: 66,
			Kota: model.Kota{
				Nama: "BANGKALAN",
			},
		},
		Total:      2500000,
		TotalNetto: 1789000,
		CreatedAt:  time.Now(),
	}

	b.WriteString(strings.Repeat("─", 82) + "\n")
	b.WriteString(fmt.Sprintf("MITRA SUKSES BERSAMA%10sNama Toko: %s\n", "", penjualan.Toko.Nama))
	b.WriteString(fmt.Sprintf("FAKTUR PENJUALAN%13s %s (%s)\n", "", penjualan.Toko.Alamat, penjualan.Toko.Kota.Nama))
	b.WriteString(fmt.Sprintf("No.Faktur: %-10s         Tanggal Faktur: %s\n",
		penjualan.NoFaktur,
		penjualan.TglPenjualan.Format("02-01-2006")))
	b.WriteString(fmt.Sprintf("%30sJatuh Tempo: %s\n", "", penjualan.TglJatuhTempo.Format("02-01-2006")))

	b.WriteString(strings.Repeat("─", 82) + "\n")
	fmt.Fprintf(&b, "%-10s", "Artikel")
	for i, kat := range []string{"K", "A", "D"} {
		if i > 0 {
			if i == 2 {
				fmt.Fprintf(&b, "%10s%-4s", "", kategoriLabel[kat])
				for _, s := range kategoriSizeRange[kat] {
					fmt.Fprintf(&b, "%-4d", s)
				}
				fmt.Fprintf(&b, "%3s%8s%9s\n", "Qty", "Harga", "Total")
			} else {
				fmt.Fprintf(&b, "%10s%-4s", "", kategoriLabel[kat])
				for _, s := range kategoriSizeRange[kat] {
					fmt.Fprintf(&b, "%-4d", s)
				}
				b.WriteString("\n")
			}
		} else {
			fmt.Fprintf(&b, "%-4s", kategoriLabel[kat])
			for _, s := range kategoriSizeRange[kat] {
				fmt.Fprintf(&b, "%-4d", s)
			}
			pad := totalSizes - len(kategoriSizeRange[kat])
			for i := 0; i < pad; i++ {
				fmt.Fprintf(&b, "%-6s", "")
			}
			b.WriteString("\n")
		}
	}
	b.WriteString(strings.Repeat("─", 82) + "\n")
	return b.String()
}

// ======================= PRINT DATA =======================
func PrintData(pivot map[string]*PivotRow, kategoriSizeRange map[string][]int) string {
	var b strings.Builder
	names := []string{}
	for nama := range pivot {
		names = append(names, nama)
	}
	sort.Strings(names)

	for _, nama := range names {
		p := pivot[nama]
		fmt.Fprintf(&b, "%-10s %-3s", p.Artikel, p.Kategori)

		row := []int{}
		for _, kat := range []string{"K", "A", "D"} {
			for _, s := range kategoriSizeRange[kat] {
				row = append(row, p.QtyMap[kat][s])
			}
		}

		if len(row) > 12 {
			row = row[:len(row)-12]
		}

		totalQty := 0
		for _, v := range row {
			if v == 0 {
				fmt.Fprintf(&b, "%-4s", "")
			} else {
				fmt.Fprintf(&b, "%-4d", v)
				totalQty += v
			}
		}

		totalUang := totalQty * p.Harga
		fmt.Fprintf(&b, "%3d %7d %8d\n", totalQty, p.Harga, totalUang)
		b.WriteString(strings.Repeat("─", 82) + "\n")
	}

	return b.String()
}

// ======================= PRINT HEADER =======================
func PrintSumarry() string {
	var b strings.Builder

	// dummy data penjualan
	penjualan := model.Penjualan{
		ID:            111,
		NoFaktur:      "JT-001",
		TglPenjualan:  time.Now(),
		TglJatuhTempo: time.Now(),
		TokoID:        1,
		Toko: model.Toko{
			Nama:   "AGRIS SPORT",
			Alamat: "JL. KH Zainal Alim No. 31. Kemayoran.",
			KotaID: 66,
			Kota: model.Kota{
				Nama: "BANGKALAN",
			},
			Disc1: 0.3,
			Disc2: 0,
			Disc3: 0,
		},
		Total:      2500000,
		TotalNetto: 2644740,
		CreatedAt:  time.Now(),
	}

	total := 3778200

	// perhitungan diskon 1,2, dan 3
	var grandTotal float64
	grandTotal = float64(total)

	// fmt.Println("Total bruto:", grandTotal)
	diskon1 := penjualan.Toko.Disc1 * grandTotal
	// fmt.Println("Diskon 1:", diskon1)
	grandTotal -= diskon1
	// fmt.Println("Grand total after diskon 1:", grandTotal)
	diskon2 := penjualan.Toko.Disc2 * grandTotal
	// fmt.Println("Diskon 2:", diskon2)
	grandTotal -= diskon2
	// fmt.Println("Grand total after diskon 2:", grandTotal)
	diskon3 := penjualan.Toko.Disc3 * grandTotal
	// fmt.Println("Diskon 3:", diskon3)
	grandTotal -= diskon3
	// fmt.Println("Grand total after diskon 3:", grandTotal)

	//grandTotal = grandTotal - diskon3
	b.WriteString(fmt.Sprintf("%73s %8d\n", "Total", total))

	b.WriteString(fmt.Sprintf("%73s %8.0f\n", "Disc 1", diskon1)) // %9.0F ARTINYA pembualan ke bilangan bulat tanpa desimal
	b.WriteString(fmt.Sprintf("%73s %8.0f\n", "Disc 2", diskon2))
	b.WriteString(fmt.Sprintf("%73s %8.0f\n", "Disc 3", diskon3))

	b.WriteString(fmt.Sprintf("%73s %8d\n", "Grand Total", penjualan.TotalNetto))
	b.WriteString(strings.Repeat("─", 82) + "\n")
	b.WriteString(centerText("TERIMA KASIH ATAS PEMBELIAN ANDA", 80) + "\n")
	b.WriteString(strings.Repeat("─", 82) + "\n")

	return b.String()
}

// agar text terima kasih berada di tengah
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text
}

// ======================= WRITE FILE =======================
func WriteToFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// ======================= FORMAT RIBUAN =======================
func FormatRibuan(n int) string {
	s := fmt.Sprintf("%d", n)
	l := len(s)
	if l <= 3 {
		return s
	}

	var result strings.Builder
	pre := l % 3
	if pre > 0 {
		result.WriteString(s[:pre])
		if l > pre {
			result.WriteString(".")
		}
	}

	for i := pre; i < l; i += 3 {
		result.WriteString(s[i : i+3])
		if i+3 < l {
			result.WriteString(".")
		}
	}

	return result.String()
}
