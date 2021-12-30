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
	customersSvc 	*customers.Service
	securitySvc		*security.Service
}
func NewServer (mux *mux.Router, customerSvc *customers.Service, security *security.Service) *Server {
	return &Server{mux: mux, customersSvc: customerSvc, securitySvc: security}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request){
	s.mux.ServeHTTP(writer, request)
}
// Token ...
type Token struct {
	Token string `json:"token"`
}

// Responce ...
type Responce struct {
	CustomerID int64  `json:"customerId"`
	Status     string `json:"status"`
	Reason     string `json:"reason"`
}

// ResponceOk ...
type ResponceOk struct {
	Status     string `json:"status"`
	CustomerID int64  `json:"customerId"`
}

// ResponceFail ...
type ResponceFail struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}
const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

//Init ...
func (s *Server) Init() {
	s.mux.Use(middleware.Basic(s.securitySvc.Auth))
	
	s.mux.HandleFunc("/api/customers", s.saveCustomers).Methods(POST)
	s.mux.HandleFunc("/api/customers/token", s.handleGetToken).Methods(POST)
	s.mux.HandleFunc("/api/customers/token/validate", s.handleValidateToken).Methods(POST)

	s.mux.HandleFunc("/customers", s.handleGetAllCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/active", s.handleGetAllActiveCustomers).Methods(GET)
	s.mux.HandleFunc("/customers/{id}", s.handleGetCustomerByID).Methods(GET)
	s.mux.HandleFunc("/customers/{id}/block", s.handleBlockByID).Methods(POST)
	s.mux.HandleFunc("/customers/{id}/block", s.handleUnBlockByID).Methods(DELETE)
	s.mux.HandleFunc("/customers/{id}", s.handleDelete).Methods(DELETE)
	s.mux.HandleFunc("/customers", s.handleSave).Methods(POST)

}

func (s *Server) handleGetToken(w http.ResponseWriter, r *http.Request) {
	var auth *security.Auth
	var tok Token
	err := json.NewDecoder(r.Body).Decode(&auth)
	log.Print(auth)
	if err != nil {
		log.Print("Can't Decode login and password")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	log.Print("Login: ", auth.Login, "Password: ", auth.Password)

	token, err := s.customersSvc.TokenForCustomer(r.Context(), auth.Login, auth.Password)
	if err != nil {
		log.Print(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tok.Token = token
	data, err := json.Marshal(tok)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) {
	var fail ResponceFail
	var ok ResponceOk
	var token Token
	var data []byte
	code := 200

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		log.Print("Can't Decode token")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	id, er := s.securitySvc.AuthenticateCusomer(r.Context(), token.Token)

	if er == security.ErrNoSuchUser {
		code = 404
		fail.Status = "fail"
		fail.Reason = "not found"
	} else if er == security.ErrExpired {
		code = 400
		fail.Status = "fail"
		fail.Reason = "expired"
	} else if er == nil {
		log.Print(id)
		ok.Status = "ok"
		ok.CustomerID = id
	} else {
		log.Print("err", er)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if code != 200 {
		w.WriteHeader(code)

		data, err = json.Marshal(fail)
		if err != nil {
			log.Print(err)
		}
	} else {
		data, err = json.Marshal(ok)
		if err != nil {
			log.Print(err)
		}
	}
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	return
}

func (s *Server) saveCustomers(writer http.ResponseWriter, request *http.Request) {
	var item *customers.Customer
	err := json.NewDecoder(request.Body).Decode(&item)
	log.Print(item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	cust, err := s.customersSvc.SaveCustomer(request.Context(), item)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print(cust)

	data, err := json.Marshal(cust)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Contetn-Type", "applicatrion/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
	}
	http.Error(writer, http.StatusText(200), 200)
	return

}


// хендлер метод для извлечения всех клиентов
func (s *Server) handleGetAllCustomers(w http.ResponseWriter, r *http.Request) {

	items, err := s.customersSvc.All(r.Context())
	if err != nil {
		//вызываем фукцию для ответа с ошибкой
		errorWriter(w, http.StatusInternalServerError, err)
		return
	}

	respondJSON(w, items)
}

// хендлер метод для извлечения всех активных клиентов
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
	item, err := s.customersSvc.ByID(r.Context(), id)

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
	item, err := s.customersSvc.ChangeActive(r.Context(), id, false)

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
	item, err := s.customersSvc.ChangeActive(r.Context(), id, true)

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
	item, err := s.customersSvc.Delete(r.Context(), id)

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
	customer, err := s.customersSvc.Save(r.Context(), item)

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

