package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

// Membuat struct Mahasiswa
type Mahasiswa struct {
	ID        int       `form:"id" json:"id"`
	NIM       int       `json:"nim"`
	Name      string    `json:"name"`
	Semester  int       `json:"semester"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Membuat array dari Mahasiswa
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/db_apigo2?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection Success")
	}
	db.AutoMigrate(&Mahasiswa{})
	handleRequests()
}

func handleRequests() {
	log.Println("Start the development server at localhost:1991")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := Result{Code: 404, Message: "Method not found"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	myRouter.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := Result{Code: 403, Message: "Method not allowed"}
		response, _ := json.Marshal(res)
		w.Write(response)
	})

	//Martin Halomoan Lumbangaol
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/Mahasiswas", createMahasiswa).Methods("POST")
	myRouter.HandleFunc("/api/Mahasiswas", getMahasiswas).Methods("GET")
	myRouter.HandleFunc("/api/Mahasiswas/{id}", getMahasiswa).Methods("GET")
	myRouter.HandleFunc("/api/Mahasiswas/{id}", updateMahasiswa).Methods("PUT")
	myRouter.HandleFunc("/api/Mahasiswas/{id}", deleteMahasiswa).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":1991", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the Martin Website!")
}

func createMahasiswa(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var Mahasiswa Mahasiswa
	json.Unmarshal(payloads, &Mahasiswa)

	db.Create(&Mahasiswa)

	res := Result{Code: 200, Data: Mahasiswa, Message: "Berhasil membuat data mahasiswa"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Martin Halomoan Lumbangaol
func getMahasiswas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get Mahasiswas")

	Mahasiswas := []Mahasiswa{}
	db.Find(&Mahasiswas)

	res := Result{Code: 200, Data: Mahasiswas, Message: "Berhasil"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

//Martin Halomoan Lumbangaol
func getMahasiswa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	MahasiswaID := vars["id"]

	var Mahasiswa Mahasiswa

	db.First(&Mahasiswa, MahasiswaID)

	res := Result{Code: 200, Data: Mahasiswa, Message: "Berhasil mengambil data mahasiswa"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Martin Halomoan Lumbangaol
func updateMahasiswa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	MahasiswaID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)

	var MahasiswaUpdates Mahasiswa
	json.Unmarshal(payloads, &MahasiswaUpdates)

	var Mahasiswa Mahasiswa
	db.First(&Mahasiswa, MahasiswaID)
	db.Model(&Mahasiswa).Updates(MahasiswaUpdates)

	res := Result{Code: 200, Data: Mahasiswa, Message: "update data mahasiswa berhasil"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

//Martin Halomoan Lumbangaol
func deleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	MahasiswaID := vars["id"]

	var Mahasiswa Mahasiswa

	db.First(&Mahasiswa, MahasiswaID)
	db.Delete(&Mahasiswa)

	res := Result{Code: 200, Message: "data mahasiswa berhasil dihapus"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
