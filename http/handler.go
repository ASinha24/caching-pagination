package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/alka/supermart"
	"github.com/alka/supermart/api"
	"github.com/alka/supermart/http/utils"
	"github.com/alka/supermart/store"
)

var (
	Page  int
	Begin int
)

func InstallRoutes(mux *mux.Router) {
	// Create a new item of super mart
	mux.Methods(http.MethodPost).Path("/api/v1/supermarts/items").HandlerFunc(createItem)
	// Update an existing item of a mart
	mux.Methods(http.MethodPut).Path("/api/v1/supermarts/items/{itemID}").HandlerFunc(updateItem)
	// delete any item of a mart
	mux.Methods(http.MethodDelete).Path("/api/v1/supermarts/items/{itemID}").HandlerFunc(deleteItem)
	// List all items of a supermart
	mux.Methods(http.MethodGet).Path("/api/v1/supermarts/items").HandlerFunc(getItems)
	mux.Methods(http.MethodGet).Path("/api/v1/supermarts/items/kafka").HandlerFunc(sendKafkaNotification)
}

//createItem create item
func createItem(w http.ResponseWriter, r *http.Request) {
	itemCreateReq := &api.ItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(itemCreateReq); err != nil {
		fmt.Println(err)
		utils.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}
	resp, err := supermart.CreateItem(r.Context(), itemCreateReq)
	if err != nil {
		utils.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	utils.WriteResponse(http.StatusCreated, resp, w)
}

//updateItem update item
func updateItem(w http.ResponseWriter, r *http.Request) {
	itemID := mux.Vars(r)["itemID"]

	itemCreateReq := &api.ItemRequest{}
	if err := json.NewDecoder(r.Body).Decode(itemCreateReq); err != nil {
		utils.WriteErrorResponse(http.StatusBadRequest, err, w)
		return
	}

	update := bson.D{
		{"$set", bson.D{
			{"name", itemCreateReq.Name},
			{"price", itemCreateReq.Price},
			{"quantity", itemCreateReq.Quantity},
		}},
	}
	resp, err := supermart.UpdateItem(r.Context(), itemID, update, itemCreateReq)
	if err != nil {
		utils.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	resp.ID = itemID
	utils.WriteResponse(http.StatusCreated, resp, w)
}

//deleteItem delete item
func deleteItem(w http.ResponseWriter, r *http.Request) {
	itemID := mux.Vars(r)["itemID"]
	res, err := supermart.DeleteItem(r.Context(), itemID)
	if err != nil {
		utils.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	utils.WriteResponse(http.StatusNoContent, res, w)
	return
}

//getItems get all items
func getItems(w http.ResponseWriter, r *http.Request) {
	items, err := supermart.GetItems(r.Context(), r)
	if err != nil {
		utils.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	utils.WriteResponse(http.StatusOK, items, w)
}

//sendKafkaNotification kafka notification
func sendKafkaNotification(w http.ResponseWriter, r *http.Request) {
	items, err := store.SendKafkaNotification()
	if err != nil {
		utils.WriteErrorResponse(http.StatusInternalServerError, err, w)
		return
	}
	utils.WriteResponse(http.StatusOK, items, w)
}
