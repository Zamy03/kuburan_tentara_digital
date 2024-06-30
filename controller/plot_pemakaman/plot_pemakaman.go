package plotpemakaman

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"kuburan/database"
	"kuburan/model/plot_pemakaman"
)

func GetPlotPemakaman(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query("SELECT * FROM plot_pemakaman")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var plot_pemakamans []plotpemakaman.PlotPemakaman
	for rows.Next() {
		var c plotpemakaman.PlotPemakaman
		if err := rows.Scan(&c.PlotId, &c.NomorPlot, &c.StatusPlot, &c.Keterangan); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		plot_pemakamans = append(plot_pemakamans, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plot_pemakamans)
}

func PostPlotPemakaman(w http.ResponseWriter, r *http.Request) {
	var pc plotpemakaman.PlotPemakaman
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for inserting a new plot pemakaman
	query := `
		INSERT INTO plot_pemakaman (nomor_plot, status_plot, keterangan)
		VALUES (?, ?, ?)`
	// Execute the SQL statement
	res, err :=
		database.DB.Exec(query, pc.NomorPlot, pc.StatusPlot, pc.Keterangan)
	if err != nil {
		http.Error(w, "Failed to insert plot pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id, 
		"message": "Plot pemakaman created successfully",
	})
}

func PutPlotPemakaman(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid plot pemakaman ID", http.StatusBadRequest)
		return
	}

	var pc plotpemakaman.PlotPemakaman
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for updating a plot pemakaman
	query := `
		UPDATE plot_pemakaman
		SET nomor_plot=?, status_plot=?, keterangan=?
		WHERE plot_id=?`

	// Execute the SQL statement
	_, err = database.DB.Exec(query, pc.NomorPlot, pc.StatusPlot, pc.Keterangan, id)
	if err != nil {
		http.Error(w, "Failed to update plot pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Plot pemakaman updated successfully"})
}

func DeletePlotPemakaman(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid plot pemakaman ID", http.StatusBadRequest)
		return
	}

	// Prepare the SQL statement for deleting a plot pemakaman
	query := `DELETE FROM plot_pemakaman WHERE plot_id=?`

	// Execute the SQL statement
	_, err = database.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete plot pemakaman: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Plot pemakaman deleted successfully"})
}