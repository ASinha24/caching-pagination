package supermart

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/alka/supermart/api"
	"github.com/alka/supermart/cache"
	"github.com/alka/supermart/http/utils"
	"github.com/alka/supermart/store"
	"github.com/alka/supermart/store/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var cacheInmem []*model.Item

//CreateItem create item and parallely save into cache
func CreateItem(ctx context.Context, item *api.ItemRequest) (*api.CreateItemRespose, error) {
	res, err := store.CreateItem(ctx, &model.Item{
		ID:       uuid.New(),
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	})
	if err != nil {
		return nil, &api.MartError{Code: api.ItemCreationFailed, Message: "can not create new item", Description: err.Error()}
	}
	//cache
	cacheInmem = append(cacheInmem, &model.Item{
		ID:       res.InsertedID.(string),
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	})

	if ok := cache.SetCache("mart-data", cacheInmem); !ok {
		fmt.Println("failed to add data into cache")
	}
	return &api.CreateItemRespose{ID: res.InsertedID.(string), ItemRequest: item}, nil
}

//DeleteItem delete item and parallely save the update into cache
func DeleteItem(ctx context.Context, itemID string) (*mongo.DeleteResult, error) {
	res, err := store.DeleteItem(ctx, bson.M{"_id": itemID})
	if err != nil {
		return nil, &api.MartError{Code: api.ItemCreationFailed, Message: "failed to delete item", Description: err.Error()}
	}
	//cache
	temp := []*model.Item{}
	for _, c := range cacheInmem {
		if c.ID == itemID {
			continue
		}
		temp = append(temp, c)
	}
	cacheInmem = temp
	fmt.Println(&cacheInmem)
	if ok := cache.SetCache("mart-data", cacheInmem); !ok {
		fmt.Println("failed to add data into cache")
	}
	return res, nil
}

//UpdateItem update item and parallely save the update into cache
func UpdateItem(ctx context.Context, itemID string, update interface{}, item *api.ItemRequest) (*api.CreateItemRespose, error) {

	err := store.UpdateItem(ctx, bson.M{"_id": itemID}, update, &model.Item{
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	})
	if err != nil {
		return nil, &api.MartError{Code: api.ItemUpdateFailed, Message: "can't update item", Description: err.Error()}
	}
	index := 0
	for i, c := range cacheInmem {
		if c.ID == itemID {
			index = i
			break
		}
	}
	cacheInmem[index] = &model.Item{
		ID:       itemID,
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	}
	fmt.Println(&cacheInmem)
	if ok := cache.SetCache("mart-data", cacheInmem); !ok {
		fmt.Println("failed to add data into cache")
	}
	data := cache.GetCache("mart-data")
	fmt.Println(&data)
	return &api.CreateItemRespose{ItemRequest: item}, nil
}

//GetItems get item from cache
func GetItems(ctx context.Context, r *http.Request) ([]api.CreateItemRespose, error) {

	items, err := store.GetAllMartItems(ctx, bson.M{})
	if err != nil {
		return nil, &api.MartError{Code: api.ItemNotFound, Message: fmt.Sprintf("failed to get data"), Description: err.Error()}
	}
	limit := 3
	page, begin := utils.Pagination(r, limit)
	log.Printf("current page %d Begin %d", page, begin)
	items = store.PaginateItems(begin, limit, items)

	martItems := []api.CreateItemRespose{}
	for _, item := range items {
		martItems = append(martItems, api.CreateItemRespose{
			ItemRequest: &api.ItemRequest{
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			},
			ID: item.ID,
		})
	}
	return martItems, nil
}
