package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//WriteResponse response writer to client
func WriteResponse(status int, response interface{}, rw http.ResponseWriter) {
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	if response != nil {
		json.NewEncoder(rw).Encode(response)
	}
}

//WriteErrorResponse eror writer to client
func WriteErrorResponse(status int, err error, rw http.ResponseWriter) {
	// martErr, ok := err.(*api.MartError)
	// if !ok {
	// 	martErr = &api.MartError{
	// 		Code:        0,
	// 		Message:     "failed in serving request",
	// 		Description: martErr.Error(),
	// 	}
	// }
	// status = martErr.Code.HTTPStatus()
	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(status)
	json.NewEncoder(rw).Encode(err.Error())
}

func Pagination(r *http.Request, limit int) (int, int) {
	keys := r.URL.Query()
	fmt.Println("keys = ", keys)
	if keys.Get("page") == "" {
		return 1, 0
	}
	page, _ := strconv.Atoi(keys.Get("page"))
	if page < 1 {
		return 1, 0
	}
	begin := (limit * page) - limit
	return page, begin
}
