package app

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/NeverlandMJ/CRUD/pkg/customers"
	"github.com/gorilla/mux"
)

type Server struct {
	mux 			*mux.Router
	customersSvc 	*customers.Service
}

func NewServer (mux *mux.Router, customersSvc *customers.Service) *Server {
	return &Server{mux: mux, customersSvc: customersSvc}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	s.mux.ServeHTTP(writer, request)
}

const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

// Vars returns the route variable for the current request, if any
func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value("id"); rv != nil {
		return rv.(map[string]string)
	}

	return nil
}


//Init inisializes server (regetres all Handlers)
// func (s *Server) Init()  {
// 	s.mux.HandleFunc("/customers.getById", s.handleGetCustomerByID)
// 	s.mux.HandleFunc("/customers.getAll", s.handleGetAllCustomer)
// 	s.mux.HandleFunc("/customers.getAllActive", s.handleGetAllActiveCustomers)
// 	s.mux.HandleFunc("/customers.save", s.handleSave)
// 	s.mux.HandleFunc("/customers.removeById", s.handleRemoveById)
// 	s.mux.HandleFunc("/customers.blockById", s.handleBlockById)
// 	s.mux.HandleFunc("/customers.unblockById", s.handleUnblockById)
// }

func (s *Server) Init() {
	s.mux.HandleFunc("/customers", s.handleGetAllCustomer).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(POST)
	s.mux.HandleFunc("customers", s.handleSave).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveById).Methods(DELETE)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockById).Methods(POST)
	s.mux.HandleFunc("/cutomers/{id}/block", s.handleBlockById).Methods(DELETE)
}



func (s *Server) handleBlockById(w http.ResponseWriter, r *http.Request){
	idPharm := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idPharm, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.BlockById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
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

func (s *Server) handleUnblockById(w http.ResponseWriter, r *http.Request){
	idPharm := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idPharm, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.UnblockById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
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
	//idParam := r.URL.Query().Get("id")
	
	idParam, ok := mux.Vars(r)["id"]
	
	if !ok{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil{
		log.Print(err)
		return
	}

	item, err := s.customersSvc.ById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(item)
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

func (s *Server) handleGetAllCustomer(w http.ResponseWriter, r *http.Request)  {
		
	items, err := s.customersSvc.All(r.Context())
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

func (s *Server) handleGetAllActiveCustomers(w http.ResponseWriter, r *http.Request)  {
	items, err := s.customersSvc.GetAllActive(r.Context())
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

func (s *Server) handleSave(w http.ResponseWriter, r *http.Request){
	id := r.FormValue("id")
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	newId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if phone == "" && name == ""{
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}


	item := &customers.Customer{
		ID: newId,
		Name: name,
		Phone: phone,
	}

	NewItem, err := s.customersSvc.Save(r.Context(), item)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(NewItem)
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

func (s *Server) handleRemoveById(w http.ResponseWriter, r *http.Request)  {
	idPharm := r.URL.Query().Get("id")

	id, err := strconv.ParseInt(idPharm, 10, 64)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.customersSvc.RemoveById(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound){
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
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