package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Siswa struct {
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
	Foto  string `json:"foto"`
}

type Guru struct {
	Nama  string `json:"nama"`
	Mapel string `json:"mapel"`
	NUPTK string `json:"nuptk"`
	Email string `json:"email"`
	Foto  string `json:"foto"`
}

type Jadwal struct {
	Kelas string `json:"kelas"`
	Hari  string `json:"hari"`
	Jam   string `json:"jam"`
	Mapel string `json:"mapel"`
	Guru  string `json:"guru"`
}

type Nilai struct {
	Nama  string `json:"nama"`
	Kelas string `json:"kelas"`
	Mapel string `json:"mapel"`
	Nilai int    `json:"nilai"`
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/index.html")
	})
	http.HandleFunc("/siswa", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/siswa.html")
	})
	http.HandleFunc("/guru", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/guru.html")
	})
	http.HandleFunc("/jadwal", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/jadwal.html")
	})
	http.HandleFunc("/nilai", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "template/nilai.html")
	})

	http.HandleFunc("/upload", uploadHandler)

	http.HandleFunc("/api/siswa", handleSiswa)
	http.HandleFunc("/api/guru", handleGuru)
	http.HandleFunc("/api/jadwal", handleJadwal)
	http.HandleFunc("/api/nilai", handleNilai)

	fmt.Println("Server running di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Metode tidak diizinkan", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Gagal membaca file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		http.Error(w, "Hanya format JPG, JPEG, PNG yang diizinkan", http.StatusBadRequest)
		return
	}

	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join("static", "foto", newFileName)

	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Gagal menyimpan file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Gagal menyalin file", http.StatusInternalServerError)
		return
	}

	url := "/static/foto/" + newFileName
	resp := map[string]string{"url": url}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleSiswa(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("data/siswa.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		var siswa []Siswa
		err := json.Unmarshal(body, &siswa)
		if err != nil {
			http.Error(w, "JSON error: "+err.Error(), 400)
			return
		}
		jsonData, _ := json.MarshalIndent(siswa, "", "  ")
		os.WriteFile("data/siswa.json", jsonData, 0644)
		w.Write([]byte("OK"))
	}
}

func handleGuru(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("data/guru.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		var guru []Guru
		err := json.Unmarshal(body, &guru)
		if err != nil {
			http.Error(w, "JSON error: "+err.Error(), 400)
			return
		}
		jsonData, _ := json.MarshalIndent(guru, "", "  ")
		os.WriteFile("data/guru.json", jsonData, 0644)
		w.Write([]byte("OK"))
	}
}

func handleJadwal(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, err := os.ReadFile("data/jadwal.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else if r.Method == "POST" {
		body, _ := io.ReadAll(r.Body)
		var jadwal []Jadwal
		err := json.Unmarshal(body, &jadwal)
		if err != nil {
			http.Error(w, "JSON error: "+err.Error(), 400)
			return
		}
		jsonData, _ := json.MarshalIndent(jadwal, "", "  ")
		os.WriteFile("data/jadwal.json", jsonData, 0644)
		w.Write([]byte("OK"))
	}
}

func handleNilai(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Menampilkan semua data nilai
		data, err := os.ReadFile("data/nilai.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	} else if r.Method == "POST" {
		// Menyimpan data nilai siswa baru
		body, _ := io.ReadAll(r.Body)
		var nilai []Nilai
		err := json.Unmarshal(body, &nilai)
		if err != nil {
			http.Error(w, "JSON error: "+err.Error(), 400)
			return
		}

		// Menyimpan data nilai ke file
		jsonData, _ := json.MarshalIndent(nilai, "", "  ")
		err = os.WriteFile("data/nilai.json", jsonData, 0644)
		if err != nil {
			http.Error(w, "Gagal menyimpan data nilai", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("Data nilai berhasil disimpan"))

	} else if r.Method == "PUT" {
		// Menghitung rata-rata nilai berdasarkan nama siswa
		namaSiswa := r.URL.Query().Get("nama_siswa")
		if namaSiswa == "" {
			http.Error(w, "Nama siswa diperlukan", http.StatusBadRequest)
			return
		}

		// Membaca semua data nilai
		data, err := os.ReadFile("data/nilai.json")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var nilai []Nilai
		err = json.Unmarshal(data, &nilai)
		if err != nil {
			http.Error(w, "Gagal memproses data nilai", http.StatusInternalServerError)
			return
		}

		// Hitung rata-rata
		rataRata, err := hitungRataRata(nilai, namaSiswa)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Mengirimkan hasil rata-rata sebagai JSON
		w.Header().Set("Content-Type", "application/json")
		response := map[string]float64{
			"rata_rata": rataRata,
		}
		json.NewEncoder(w).Encode(response)
	}

}
func hitungRataRata(nilaiList []Nilai, namaSiswa string) (float64, error) {
	var total int
	var count int

	for _, n := range nilaiList {
		if n.Nama == namaSiswa {
			total += n.Nilai
			count++
		}
	}

	if count == 0 {
		return 0, fmt.Errorf("Siswa dengan nama '%s' tidak ditemukan", namaSiswa)
	}

	rataRata := float64(total) / float64(count)
	return rataRata, nil
}
