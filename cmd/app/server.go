package app

import (
	"encoding/json"
	"errors"
	//"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
	s.mux.HandleFunc("customers", s.handleSaveCustomer).Methods(POST)
	s.mux.HandleFunc("/customers/{id}", s.handleRemoveById).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnblockById).Methods(POST)
	s.mux.HandleFunc("/cutomers/{id}/block", s.handleBlockById).Methods(DELETE)
	s.mux.HandleFunc("/customers", s.handleGetAllCustomer).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
}



func (s *Server) handleBlockById(w http.ResponseWriter, r *http.Request){
	//idPharm := r.URL.Query().Get("id")

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil{
		log.Print(err)
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
	//idPharm := r.URL.Query().Get("id")

	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil{
		log.Print(err)
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
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil{
		log.Print(err)
		return
	}

	item, err := s.customersSvc.ByID(r.Context(), id)
	if errors.Is(err, customers.ErrNotFound){
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		log.Print(err)
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
	items, err := s.customersSvc.AllActive(r.Context())
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

func (s *Server) handleRemoveById(w http.ResponseWriter, r *http.Request)  {
	//idPharm := r.URL.Query().Get("id")
	idParam, ok := mux.Vars(r)["id"]
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil{
		log.Print(err)
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

func (s *Server) handleSaveCustomer(w http.ResponseWriter, r *http.Request){
	var item *customers.Customer
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}	
	
	// b, err := ioutil.ReadAll(r.Body)
	// defer r.Body.Close()
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	// var item *customers.Customer
	// err = json.Unmarshal(b, &item)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }
	item, err = s.customersSvc.Save(r.Context(), item.ID, item.Phone, item.Name)

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

