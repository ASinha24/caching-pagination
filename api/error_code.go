package api

import "net/http"

type ErrorCode int

const (
	Unknown ErrorCode = iota
	ItemNotFound
	StoreNotFound
	ItemCreationFailed
	UnauthorisedStore
	ItemUpdateFailed
)

var statuCode = map[ErrorCode]int{
	ItemNotFound:      http.StatusNotFound,
	StoreNotFound:     http.StatusNotFound,
	UnauthorisedStore: http.StatusUnauthorized,
}

func (e ErrorCode) HTTPStatus() int {
	if code, ok := statuCode[e]; ok {
		return code
	}
	return http.StatusInternalServerError
}
