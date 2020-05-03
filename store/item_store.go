package store

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alka/supermart/cache"
	"github.com/alka/supermart/db"
	"github.com/alka/supermart/kafka"
	"github.com/alka/supermart/store/model"
)

var (
	Collection = db.Dbconnect()
)

//CreateItem save the item to mongo db
func CreateItem(ctx context.Context, item *model.Item) (*mongo.InsertOneResult, error) {
	res, err := Collection.InsertOne(context.Background(), item)
	if err != nil {
		return nil, err
	}
	return res, nil
}

//UpdateItem update the item to mongo db
func UpdateItem(ctx context.Context, filter interface{}, update interface{}, item *model.Item) error {
	if err := Collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&item); err != nil {
		return err
	}

	return nil
}

//DeleteItem delete the item from mongo db
func DeleteItem(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	deleteResult, err := Collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

//GetAllMartItemms, when there is no item in cahche the data will fetch from db and if there is item in cahche
//it will get return from cache
func GetAllMartItems(ctx context.Context, filter interface{}) ([]*model.Item, error) {
	var items []*model.Item
	items = cache.GetCache("mart-data")
	if len(items) > 0 {
		fmt.Println(len(items))
		return items, nil
	}

	cur, err := Collection.Find(context.TODO(), filter)
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

//PaginateItems paginate the Get items
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
	return pItems
}

//SendKafkaNotification send kafka notification
func SendKafkaNotification() ([]*model.Item, error) {
	var items []*model.Item
	cur, err := Collection.Find(context.TODO(), bson.M{})
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
	if ok := SendKafkaMsg(items); !ok {
		log.Println("failed to publish data to kafka")
	}
	return items, nil
}

//SendKafkaNotification send kafka notification
func SendKafkaMsg(items []*model.Item) bool {
	prd, err := kafka.NewProducer()
	if err != nil {
		log.Println("failed to create producer")
		return false
	}
	kafka.PublishMsg(items, prd)
	return true
}
