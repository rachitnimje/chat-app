package server

import (
	"fmt"
	"net/http"
)

func StartHTTPServer(router http.Handler, port string) error {
	address := fmt.Sprintf(":%s", port)
	return http.ListenAndServe(address, router)
}
