package server

import (
	"fmt"
	"net/http"
)

func StartHTTPServer(router http.Handler, port int) error {
	address := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(address, router)
}
