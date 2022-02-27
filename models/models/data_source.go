package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DataSource struct {
	Id   primitive.ObjectID `json:"_id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
}

func (ds *DataSource) GetId() (id primitive.ObjectID) {
	return ds.Id
}

func (ds *DataSource) SetId(id primitive.ObjectID) {
	ds.Id = id
}
