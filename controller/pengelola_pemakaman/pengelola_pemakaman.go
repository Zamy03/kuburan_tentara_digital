package pengelolapemakaman

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"kuburan/database"
	"kuburan/model/pengelola_pemakaman"
)

func GetPengelolaPemakaman(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM pengelola_pemakaman")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pengelola_pemakamans []pengelolapemakaman.PengelolaPemakaman
	for rows.Next() {
		var c pengelolapemakaman.PengelolaPemakaman
		if err := rows.Scan(&c.PengelolaId, &c.NamaLengkap, &c.Jabatan, &c.NomorTelepon, &c.Email, &c.Alamat); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pengelola_pemakamans = append(pengelola_pemakamans, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pengelola_pemakamans)
}

func PostPengelolaPemakaman(w http.ResponseWriter, r *http.Request) {
	var pc pengelolapemakaman.PengelolaPemakaman
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new pengelola pemakaman
	query := `
		INSERT INTO pengelola_pemakaman (nama_lengkap, jabatan, nomor_telepon, email, alamat)
		VALUES (?, ?, ?, ?, ?)`

	// Execute the SQL statement
	res, err := database.DB.Exec(query, pc.NamaLengkap, pc.Jabatan, pc.NomorTelepon, pc.Email, pc.Alamat)
	if err != nil {
		http.Error(w, "Failed to insert pengelola pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pc.PengelolaId = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pc)
}

func PutPengelolaPemakaman(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var pc pengelolapemakaman.PengelolaPemakaman
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating pengelola pemakaman
	query := `
		UPDATE pengelola_pemakaman
		SET nama_lengkap = ?, jabatan = ?, nomor_telepon = ?, email = ?, alamat = ?
		WHERE pengelola_id = ?`

	// Execute the SQL statement
	_, err = database.DB.Exec(query, pc.NamaLengkap, pc.Jabatan, pc.NomorTelepon, pc.Email, pc.Alamat, id)
	if err != nil {
		http.Error(w, "Failed to update pengelola pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pc)
}

func DeletePengelolaPemakaman(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting pengelola pemakaman
	query := `DELETE FROM pengelola_pemakaman WHERE pengelola_id = ?`

	// Execute the SQL statement
	_, err = database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete pengelola pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Pengelola pemakaman deleted successfully",
	})
}