package kunjungan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"kuburan/database"
	"kuburan/model/kunjungan"
)

func GetKunjungan(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM kunjungan")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var kunjungans []kunjungan.Kunjungan
	for rows.Next() {
		var c kunjungan.Kunjungan
		if err := rows.Scan(&c.KunjunganId, &c.KuburanId, &c.TanggalKunjungan, &c.NamaPengunjung, &c.Hubungan, &c.Keterangan); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kunjungans = append(kunjungans, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kunjungans)
}

func PostKunjungan(w http.ResponseWriter, r *http.Request) {
	var pc kunjungan.Kunjungan
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Periksa apakah KuburanId ada di tabel yang terkait
	var kuburanExists bool
	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kuburan WHERE kuburan_id = ?)", pc.KuburanId).Scan(&kuburanExists)
	if err != nil || !kuburanExists {
		http.Error(w, "Kuburan ID not found", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new kunjungan
	query := `
		INSERT INTO kunjungan (kuburan_id, tanggal_kunjungan, nama_pengunjung, hubungan, keterangan)
		VALUES (?, ?, ?, ?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.KuburanId, pc.TanggalKunjungan, pc.NamaPengunjung, pc.Hubungan, pc.Keterangan)
	if err != nil {
		http.Error(w, "Failed to insert kunjungan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created ID in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kunjungan added successfully",
		"id":      id,
	})
}

func PutKunjungan(w http.ResponseWriter, r *http.Request) {
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

	// Decode JSON body
	var pc kunjungan.Kunjungan
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Periksa apakah KuburanId ada di tabel yang terkait
	var kuburanExists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM kuburan WHERE kuburan_id = ?)", pc.KuburanId).Scan(&kuburanExists)
	if err != nil || !kuburanExists {
		http.Error(w, "Kuburan ID not found", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the category admin
	query := `
		UPDATE kunjungan
		SET kuburan_id = ?, tanggal_kunjungan = ?, nama_pengunjung = ?, hubungan = ?, keterangan = ?
		WHERE kunjungan_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.KuburanId, pc.TanggalKunjungan, pc.NamaPengunjung, pc.Hubungan, pc.Keterangan, id)
	if err != nil {
		http.Error(w, "Failed to update kunjungan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were updated
	if rowsAffected == 0 {
		http.Error(w, "No rows were updated", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kunjungan updated successfully",
	})
}

func DeleteKunjungan(w http.ResponseWriter, r *http.Request) {
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

	// Prepare the SQL statement for deleting the category admin
	query := `DELETE FROM kunjungan WHERE kunjungan_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete kunjungan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to retrieve affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if any rows were deleted
	if rowsAffected == 0 {
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	// Return success message
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kunjungan deleted successfully",
	})
}