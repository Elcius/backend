package server

import (
	"github.com/gorilla/mux"
	"github.com/minha-cidade/backend/db"
	"log"
	"net/http"
	"errors"
	"strconv"
)

func apiNotFound(w http.ResponseWriter, r *http.Request) {
	writeJsonError(w, http.StatusNotFound, errors.New("Not Found"))
}

func Start() {
	log.Println("Iniciando aplicação backend...")
	router := mux.NewRouter()

	// Conecta ao banco de dados
	log.Println("Conectando ao banco de dados...")
	db.Connect("HOST", 5432, "USUARIO",
		 "SENHA", "BANCO")

	// Api
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/gastometro",
		func(w http.ResponseWriter, r *http.Request) {
			gastometro, err := db.GetGastometro()
			if err != nil {
				writeJsonError(w, http.StatusInternalServerError, err)
				return
			}

			writeJson(w, http.StatusOK, gastometro)
		}).Methods("GET")
	api.HandleFunc("/area/{area}/{ano}",
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			area := vars["area"]
			ano, err := strconv.Atoi(vars["ano"])

			if err != nil {
				writeJsonError(w, http.StatusBadRequest, errors.New("Ano inválido"))
				return
			}

			info, err := db.GetInformacoesArea(area, ano)
			if err != nil {
				writeJsonError(w, http.StatusInternalServerError, err)
				return
			}

			writeJson(w, http.StatusOK, info)
		}).Methods("GET")

	// Processa a página de 404
	api.NotFoundHandler = http.HandlerFunc(apiNotFound)

	log.Println("Listenning... :1337")
	chk(http.ListenAndServe(":1337", router))
}

func chk(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}