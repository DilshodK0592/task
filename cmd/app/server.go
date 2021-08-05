package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/DilshodK0592/task/pkg/customers"
)

// Методы запросов
const (
	GET    = "GET"    // Метод получения GET
	POST   = "POST"   // Метод отправ/обновлении POST
	DELETE = "DELETE" // Метод удаления DELETE
)

//Server ...
type Server struct {
	mux         *mux.Router
	customerSvc *customers.Service
}

//NewServer ...
func NewServer(m *mux.Router, customerSvc *customers.Service) *Server {
	return &Server{mux: m, customerSvc: customerSvc}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

//Init инициализирует сервер (регистрирует все Handler'ы)
func (s *Server) Init() {
	// s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomers)
	s.mux.HandleFunc("/quote", s.handleGetAllQuote).Methods(GET)
	// s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
	s.mux.HandleFunc("/quote/category/{category}", s.handleGetQuoteByID).Methods(GET)
	// s.mux.HandleFunc("/customers.save", s.handleSave)
	s.mux.HandleFunc("/quote", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/quote/delete/{id:[0-9]+}", s.handleDelete).Methods(DELETE)
}

// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllQuote(w http.ResponseWriter, r *http.Request) {

	items, err := s.customerSvc.All(r.Context())
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

func (s *Server) handleGetQuoteByID(w http.ResponseWriter, r *http.Request) {
	category, ok := mux.Vars(r)["category"]
	// fmt.Println(category)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customerSvc.Bycategory(r.Context(), category)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		fmt.Println(1, idParam)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		fmt.Println(1, idParam)
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	item, err := s.customerSvc.Delete(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Quotes
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	customer, err := s.customerSvc.Save(r.Context(), item)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	respondJSON(w, customer)
}

func errorWriter(w http.ResponseWriter, httpSts int, err error) {
	log.Print(err)
	http.Error(w, http.StatusText(httpSts), httpSts)
}
func respondJSON(w http.ResponseWriter, iData interface{}) {

	data, err := json.Marshal(iData)

	if err != nil {
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
	}
}
