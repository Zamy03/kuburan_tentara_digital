package main

import (
	"kuburan/database"
	"fmt"
	"kuburan/controller/auth"
	"kuburan/controller/kuburan"
	"kuburan/controller/plot_pemakaman"
	"kuburan/controller/pengelola_pemakaman"
	"kuburan/controller/tentara"
	"kuburan/controller/kunjungan"
	"log"
	"net/http"	

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()	
	fmt.Println("Hello World")

	router := mux.NewRouter()

	router.HandleFunc("/regis", auth.Registration).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")

	// Router handler kuburan
	router.HandleFunc("/kuburan", kuburan.GetKuburan).Methods("GET")
	router.HandleFunc("/kuburan", auth.JWTAuth(kuburan.PostKuburan)).Methods("POST")
	router.HandleFunc("/kuburan/{id}", auth.JWTAuth(kuburan.PutKuburan)).Methods("PUT")
	router.HandleFunc("/kuburan/{id}", auth.JWTAuth(kuburan.DeleteKuburan)).Methods("DELETE")

	// Router handler kunjungan
	router.HandleFunc("/kunjungan", kunjungan.GetKunjungan).Methods("GET")
	router.HandleFunc("/kunjungan", auth.JWTAuth(kunjungan.PostKunjungan)).Methods("POST")
	router.HandleFunc("/kunjungan/{id}", auth.JWTAuth(kunjungan.PutKunjungan)).Methods("PUT")
	router.HandleFunc("/kunjungan/{id}", auth.JWTAuth(kunjungan.DeleteKunjungan)).Methods("DELETE")

	// Router handler pengelola pemakaman
	router.HandleFunc("/pengelolapemakaman", pengelolapemakaman.GetPengelolaPemakaman).Methods("GET")
	router.HandleFunc("/pengelolapemakaman", auth.JWTAuth(pengelolapemakaman.PostPengelolaPemakaman)).Methods("POST")
	router.HandleFunc("/pengelolapemakaman/{id}", auth.JWTAuth(pengelolapemakaman.PutPengelolaPemakaman)).Methods("PUT")
	router.HandleFunc("/pengelolapemakaman/{id}", auth.JWTAuth(pengelolapemakaman.DeletePengelolaPemakaman)).Methods("DELETE")

	// Router handler plot pemakaman
	router.HandleFunc("/plotpemakaman", plotpemakaman.GetPlotPemakaman).Methods("GET")
	router.HandleFunc("/plotpemakaman", auth.JWTAuth(plotpemakaman.PostPlotPemakaman)).Methods("POST")
	router.HandleFunc("/plotpemakaman/{id}", auth.JWTAuth(plotpemakaman.PutPlotPemakaman)).Methods("PUT")
	router.HandleFunc("/plotpemakaman/{id}", auth.JWTAuth(plotpemakaman.DeletePlotPemakaman)).Methods("DELETE")

	// Router handler data mayat
	router.HandleFunc("/datamayat", tentara.GetTentara).Methods("GET")
	router.HandleFunc("/datamayat", auth.JWTAuth(tentara.PostTentara)).Methods("POST")
	router.HandleFunc("/datamayat/{id}", auth.JWTAuth(tentara.PutTentara)).Methods("PUT")
	router.HandleFunc("/datamayat/{id}", auth.JWTAuth(tentara.DeleteTentara)).Methods("DELETE")

	c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://127.0.0.1:5500"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"Content-Type", "Authorization"},
        Debug: true,
    })
	
    handler := c.Handler(router)
	
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))


}