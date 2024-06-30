package kuburan

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"kuburan/database"
	"kuburan/model/kuburan"
)

func GetKuburan(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM kuburan")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var kuburans []kuburan.Kuburan
	for rows.Next() {
		var c kuburan.Kuburan
		if err := rows.Scan(&c.KuburanId, &c.PlotId, &c.TentaraId,&c.NomorKuburan, &c.TanggalDikubur, &c.StatusKuburan, &c.Keterangan); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		kuburans = append(kuburans, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kuburans)
}

func PostKuburan(w http.ResponseWriter, r *http.Request) {
	var pc kuburan.Kuburan
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Periksa apakah PlotId dan TentaraId ada di tabel yang terkait
	var plotExists, tentaraExists bool

	err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM plot_pemakaman WHERE plot_id = ?)", pc.PlotId).Scan(&plotExists)
	if err != nil || !plotExists {
		http.Error(w, "Plot ID not found", http.StatusBadRequest)
		return
	}

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tentara WHERE tentara_id = ?)", pc.TentaraId).Scan(&tentaraExists)
	if err != nil || !tentaraExists {
		http.Error(w, "Tentara ID not found", http.StatusBadRequest)
		return
	}

	// Siapkan pernyataan SQL untuk menyisipkan kuburan baru
	query := `
		INSERT INTO kuburan (plot_id, tentara_id, nomor_kuburan, tanggal_dikubur, status_kuburan, keterangan)
		VALUES (?, ?, ?, ?, ?, ?)`

	// Jalankan pernyataan SQL
	res, err := database.DB.Exec(query, pc.PlotId, pc.TentaraId, pc.NomorKuburan, pc.TanggalDikubur, pc.StatusKuburan, pc.Keterangan)
	if err != nil {
		http.Error(w, "Failed to insert kuburan: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Ambil ID terakhir yang dimasukkan
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Kembalikan ID yang baru dibuat dalam respons
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Kuburan added successfully",
		"id":      id,
	})
}

func PutKuburan(w http.ResponseWriter, r *http.Request) {
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
	var pc kuburan.Kuburan
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Periksa apakah PlotId dan TentaraId ada di tabel yang terkait
	var plotExists, tentaraExists bool

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM plot_pemakaman WHERE plot_id = ?)", pc.PlotId).Scan(&plotExists)
	if err != nil || !plotExists {
		http.Error(w, "Plot ID not found", http.StatusBadRequest)
		return
	}

	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tentara WHERE tentara_id = ?)", pc.TentaraId).Scan(&tentaraExists)
	if err != nil || !tentaraExists {
		http.Error(w, "Tentara ID not found", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating the kuburan
	query := `
		UPDATE kuburan
		SET plot_id = ?, tentara_id = ?, nomor_kuburan = ?, tanggal_dikubur = ?, status_kuburan = ?, keterangan = ?
		WHERE kuburan_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, pc.PlotId, pc.TentaraId, pc.NomorKuburan, pc.TanggalDikubur, pc.StatusKuburan, pc.Keterangan, id)
	if err != nil {
		http.Error(w, "Failed to update kuburan: "+err.Error(), http.StatusInternalServerError)
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
		"message": "Kuburan updated successfully",
	})
}

func DeleteKuburan(w http.ResponseWriter, r *http.Request) {
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
		DELETE FROM kuburan
		WHERE kuburan_id = ?`

	// Execute the SQL statement
	result, err := database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete kuburan: "+err.Error(), http.StatusInternalServerError)
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
		"message": "Kuburan deleted successfully",
	})
}
