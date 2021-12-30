package middleware

import (
	"log"
	"net/http"
)

// func Logger(handler http.Handler) http.Handler {
// 	// middleware wraps a Handler, returning a new Handler
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		//before handler execution
// 		log.Printf("START: %s %s", r.Method, r.URL.Path)

// 		// handler execution
// 		handler.ServeHTTP(w, r)

// 		// after handler execution
// 		log.Printf("FINISH: %s %s", r.Method, r.URL.Path)
// 	})
// }

//Basic...
func Basic(auth func(login, pass string) bool) func(handler http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			username, password, ok := request.BasicAuth()
			if !ok {
				log.Print("Cant parse username and password")
				http.Error(writer, http.StatusText(401), 401)
				return
			}
			if !auth(username, password) {
				http.Error(writer, http.StatusText(401), 401)
				return
			}
			handler.ServeHTTP(writer, request)
		})
	}
}