package app

import (
	"encoding/json"
	"errors"

	//"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/NeverlandMJ/CRUD/cmd/app/middleware"
	"github.com/NeverlandMJ/CRUD/pkg/customers"
	"github.com/NeverlandMJ/CRUD/pkg/security"
	"github.com/gorilla/mux"
)

type Server struct {
	mux 			*mux.Router
	customerSvc 	*customers.Service
	securitySvc		*security.Service
}
func NewServer (mux *mux.Router, customerSvc *customers.Service) *Server {
	return &Server{mux: mux, customerSvc: customerSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	s.mux.ServeHTTP(writer, request)
}

const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

//Init ...
func (s *Server) Init() {
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))
	s.mux.Use(middleware.Logger)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)
}




// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	items, err := s.customerSvc.All(r.Context())
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

// хендлер метод для извлечения всех активных клиентов
func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request)  {
	items, err := s.customerSvc.AllActive(r.Context())
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Print(err)
		return
	}
	
}

func (s *Server) handleGetCustomerByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
		
	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ByID(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, item)
}

func (s *Server) handleBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ChangeActive(r.Context(), id, false)

	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, item)
}

func (s *Server) handleUnBlockByID(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.ChangeActive(r.Context(), id, true)

	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, item)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	//получаем ID из параметра запроса
	idP, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// переобразуем его в число
	id, err := strconv.ParseInt(idP, 10, 64)
	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusBadRequest, err)
		return
	}

	//получаем баннер из сервиса
	item, err := s.customerSvc.Delete(r.Context(), id)

	if errors.Is(err, customers.ErrNotFound) {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusNotFound, err)
		return
	}

	//если получили ошибку то отвечаем с ошибкой
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}
	//вызываем функцию для ответа в формате JSON
	respondJSON(w, item)
}

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
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

