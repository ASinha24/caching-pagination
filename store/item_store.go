package store

import (
	"context"
	"fmt"

	"github.com/alka/supermart/cache"
	"github.com/alka/supermart/db"
	"github.com/alka/supermart/store/model"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	collection = db.Dbconnect()
)

func CreateItem(ctx context.Context, item *model.Item) (*mongo.InsertOneResult, error) {
	res, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func UpdateItem(ctx context.Context, filter interface{}, update interface{}, item *model.Item) error {
	if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&item); err != nil {
		return err
	}

	return nil
}

func DeleteItem(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

func GetAllMartItems(ctx context.Context, filter interface{}) ([]*model.Item, error) {
	var items []*model.Item
	items = cache.GetCache("mart-data")
	if len(items) > 0 {
		fmt.Println(len(items))
		return items, nil
	}

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var item model.Item
		err := cur.Decode(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func PaginateItems(begin, limit int, items []*model.Item) []*model.Item {
	pItems := []*model.Item{}
	count := 0
	for index, val := range items {
		if index+1 >= begin {
			pItems = append(pItems, val)
			count++
			if count == limit {
				break
			}
		}
	}
	fmt.Println(pItems)
	return pItems
}
