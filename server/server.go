package server

import (
	"fmt"
	"github.com/plitn/wb_school_l0/storage"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type MyServer struct {
	mux   *http.ServeMux
	cache *storage.CacheStruct
}

// NewServer конструктор
func NewServer(cache *storage.CacheStruct) *MyServer {
	server := &MyServer{
		mux:   http.NewServeMux(),
		cache: cache,
	}
	return server
}

// Init инит, подключаем хендлеры, начинаем работу
func (ms *MyServer) Init() {
	ms.mux.HandleFunc("/", ms.indexHandle)
	ms.mux.HandleFunc("/orders", ms.orders)
	err := http.ListenAndServe(":3333", ms.mux)
	if err != nil {
		log.Println(err)
	}
}

// обработка главной страницы
func (ms *MyServer) indexHandle(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("public", "index.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("template error: %v", err)
	}
	err = tmpl.Execute(w, nil)
}

// обработка страницы с отображением ордеров по айдишнику
func (ms *MyServer) orders(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("public", "orders.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		log.Printf("template error: %v", err)
	}
	data := r.FormValue("id")
	ordersData, err := ms.cache.GetData(data)
	if err != nil {
		errrr := tmpl.Execute(w, err)
		if errrr != nil {
			log.Printf("error handling error, %v", errrr)
		}
		return
	}
	fmt.Println(ordersData)
	err = tmpl.Execute(w, ordersData)
}
