package tentara

import (
	"encoding/json"
	"net/http"
	"strconv"
	"os"
	"io"

	"kuburan/database"
	"path/filepath"
	"kuburan/model/tentara"

	"github.com/gorilla/mux"
)

func GetTentara(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM tentara")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tentaras []tentara.Tentara
	for rows.Next() {
		var c tentara.Tentara
		if err := rows.Scan(&c.TentaraId, &c.NamaLengkap, &c.Pangkat, &c.TanggalLahir, &c.TanggalWafat, &c.NomorIdentitas, &c.FotoTentara); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tentaras = append(tentaras, c)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tentaras)
}

// func PostTentara(w http.ResponseWriter, r *http.Request) {
// 	var ps tentara.Tentara
// 	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}
	
// 	// Limit the size of incoming file uploads to 10 MB
// 	var err error
// 	err = r.ParseMultipartForm(10 << 20)
// 	if err != nil {
// 		http.Error(w, "Error parsing form data", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve form values
// 	ps = tentara.Tentara{
// 		NamaLengkap:    r.FormValue("nama_lengkap"),
// 		Pangkat:        r.FormValue("pangkat"),
// 		TanggalLahir:   r.FormValue("tanggal_lahir"),
// 		TanggalWafat:   r.FormValue("tanggal_wafat"),
// 		NomorIdentitas: r.FormValue("nomor_identitas"),
// 	}

// 	// Retrieve file from posted form-data
// 	file, handler, err := r.FormFile("foto_tentara")
// 	if err != nil {
// 		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer file.Close()

// 	// Create a directory to save the uploaded file
// 	dir := "./uploads"
// 	if _, err := os.Stat(dir); os.IsNotExist(err) {
// 		err = os.Mkdir(dir, os.ModePerm)
// 		if err != nil {
// 			http.Error(w, "Error creating directory", http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Create a new file in the uploads directory
// 	filePath := filepath.Join(dir, handler.Filename)
// 	dst, err := os.Create(filePath)
// 	if err != nil {
// 		http.Error(w, "Error creating file", http.StatusInternalServerError)
// 		return
// 	}
// 	defer dst.Close()

// 	// Copy the uploaded file's content to the new file
// 	if _, err := io.Copy(dst, file); err != nil {
// 		http.Error(w, "Error saving file", http.StatusInternalServerError)
// 		return
// 	}

// 	// Set the file path in the struct
// 	ps.FotoTentara = filePath

// 	// Query untuk memasukan mahasiswa ke dalam table
// 	query := `
// 		INSERT INTO tentara (nama_lengkap, tanggal_lahir, tanggal_wafat, nomor_identitas, foto_tentara)
// 		VALUES (?, ?, ?, ?, ?)`

// 	// Mengeksekusi query
// 	res, err := database.DB.Exec(query, ps.NamaLengkap, ps.TanggalLahir, ps.TanggalWafat, ps.NomorIdentitas, ps.FotoTentara)
// 	if err != nil {
// 		http.Error(w, "Failed to insert tentara: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Ambil id terakhir
// 	id, err := res.LastInsertId()
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Return the newly created ID in the response
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Mayat added successfully",
// 		"id":      id,
// 	})
// }

func PostTentara(w http.ResponseWriter, r *http.Request) {
	// Batasi ukuran unggahan file hingga 10 MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Ambil nilai dari form-data
	ps := tentara.Tentara{
		NamaLengkap:    r.FormValue("nama_lengkap"),
		Pangkat:        r.FormValue("pangkat"),
		TanggalLahir:   r.FormValue("tanggal_lahir"),
		TanggalWafat:   r.FormValue("tanggal_wafat"),
		NomorIdentitas: r.FormValue("nomor_identitas"),
	}

	// Ambil file dari form-data
	file, handler, err := r.FormFile("foto_tentara")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Buat direktori untuk menyimpan file yang diunggah
	dir := "./img"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			http.Error(w, "Error creating directory", http.StatusInternalServerError)
			return
		}
	}

	// Buat file baru di direktori img
	filePath := filepath.Join(dir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Salin konten file yang diunggah ke file baru
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Setel jalur file di struct
	ps.FotoTentara = filePath

	// Sisipkan data ke dalam database
	query := `
		INSERT INTO tentara (nama_lengkap, pangkat, tanggal_lahir, tanggal_wafat, nomor_identitas, foto_tentara)
		VALUES (?, ?, ?, ?, ?, ?)`

	res, err := database.DB.Exec(query, ps.NamaLengkap, ps.Pangkat, ps.TanggalLahir, ps.TanggalWafat, ps.NomorIdentitas, ps.FotoTentara)
	if err != nil {
		http.Error(w, "Failed to insert tentara: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil ID yang baru dimasukkan
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Kembalikan ID yang baru dibuat dalam respon
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tentara added successfully",
		"id":      id,
	})
}

// func PutTentara(w http.ResponseWriter, r *http.Request) {
// 	// Ambil ID dari URL
// 	vars := mux.Vars(r)
// 	idStr, ok := vars["id"]
// 	if !ok {
// 		http.Error(w, "ID not provided", http.StatusBadRequest)
// 		return
// 	}
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	// Decode JSON body
// 	var ps tentara.Tentara
// 	if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
// 		http.Error(w, "Invalid request payload", http.StatusBadRequest)
// 		return
// 	}

// 	// Query ubah mahasiswa
// 	query := `
// 		UPDATE tentara
// 		SET nama_lengkap = ?, tanggal_lahir = ?, tanggal_wafat = ?, nomor_identitas = ?, foto_tentara = ?
// 		WHERE tentara_id = ?`

// 	// Execute the SQL statement
// 	result, err := database.DB.Exec(query, ps.NamaLengkap, ps.TanggalLahir, ps.TanggalWafat, ps.NomorIdentitas, ps.FotoTentara, id)
// 	if err != nil {
// 		http.Error(w, "Failed to update tentara: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Get the number of rows affected
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Check if any rows were updated
// 	if rowsAffected == 0 {
// 		http.Error(w, "No rows were updated", http.StatusNotFound)
// 		return
// 	}

// 	// Return success message
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(map[string]interface{}{
// 		"message": "Mayat updated successfully",
// 	})
// }

func PutTentara(w http.ResponseWriter, r *http.Request) {
	// Ambil ID dari URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Batasi ukuran unggahan file hingga 10 MB
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	// Ambil nilai dari form-data
	ps := tentara.Tentara{
		NamaLengkap:    r.FormValue("nama_lengkap"),
		Pangkat:        r.FormValue("pangkat"),
		TanggalLahir:   r.FormValue("tanggal_lahir"),
		TanggalWafat:   r.FormValue("tanggal_wafat"),
		NomorIdentitas: r.FormValue("nomor_identitas"),
	}

	// Ambil file dari form-data
	file, handler, err := r.FormFile("foto_tentara")
	if err != nil {
		if err == http.ErrMissingFile {
			// Jika tidak ada file yang diunggah, hanya perbarui field selain foto_tentara
			query := `
				UPDATE tentara
				SET nama_lengkap = ?, pangkat = ?, tanggal_lahir = ?, tanggal_wafat = ?, nomor_identitas = ?
				WHERE tentara_id = ?`
			result, err := database.DB.Exec(query, ps.NamaLengkap, ps.Pangkat, ps.TanggalLahir, ps.TanggalWafat, ps.NomorIdentitas, id)
			if err != nil {
				http.Error(w, "Failed to update tentara: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Ambil jumlah baris yang terpengaruh
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
				return
			}

			// Cek apakah ada baris yang diperbarui
			if rowsAffected == 0 {
				http.Error(w, "No rows were updated", http.StatusNotFound)
				return
			}

			// Kembalikan pesan sukses
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Tentara updated successfully",
			})
			return
		} else {
			http.Error(w, "Error retrieving the file", http.StatusInternalServerError)
			return
		}
	}
	defer file.Close()

	// Buat direktori untuk menyimpan file yang diunggah
	dir := "./img"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			http.Error(w, "Error creating directory", http.StatusInternalServerError)
			return
		}
	}

	// Buat file baru di direktori img
	filePath := filepath.Join(dir, handler.Filename)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Salin konten file yang diunggah ke file baru
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Setel jalur file di struct
	ps.FotoTentara = filePath

	// Query ubah data tentara
	query := `
		UPDATE tentara
		SET nama_lengkap = ?, pangkat = ?, tanggal_lahir = ?, tanggal_wafat = ?, nomor_identitas = ?, foto_tentara = ?
		WHERE tentara_id = ?`

	// Eksekusi statement SQL
	result, err := database.DB.Exec(query, ps.NamaLengkap, ps.Pangkat, ps.TanggalLahir, ps.TanggalWafat, ps.NomorIdentitas, ps.FotoTentara, id)
	if err != nil {
		http.Error(w, "Failed to update tentara: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil jumlah baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Cek apakah ada baris yang diperbarui
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Kembalikan pesan sukses
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Tentara updated successfully",
	})
}

func DeleteTentara(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "ID not provided", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting a category admin
	query := `
		DELETE FROM tentara
		WHERE tentara_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete tentara: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Mayat deleted successfully",
	})
}
