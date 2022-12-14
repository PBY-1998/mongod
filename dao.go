/**
    @author: potten
    @since: 2022/8/5
    @desc: //TODO
**/
package mongo_go

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Dao struct {
	client     *mongo.Client
	database   *string
	collection *string
}

type AbstractMongoDao interface {
	SetDatabase(database string) *Dao
	SetCollection(collection string) *Dao
	SetClient(client *mongo.Client) *Dao

	GetDatabase() *string
	GetCollection() *string
	GetClient() *mongo.Client

	Conn() *Curd
}

func (d *Dao) SetDatabase(database string) *Dao {
	d.database = &database
	return d
}

func (d *Dao) SetCollection(collection string) *Dao {
	d.collection = &collection
	return d
}

func (d *Dao) SetClient(client *mongo.Client) *Dao {
	d.client = client
	return d
}

func (d *Dao) GetDatabase() *string {
	return d.database
}

func (d *Dao) GetCollection() *string {
	return d.collection
}

func (d *Dao) GetClient() *mongo.Client {
	return d.client
}

func (d *Dao) Conn() *Curd {
	coll := d.client.Database(*d.database).Collection(*d.collection)
	return &Curd{
		Coll: coll,
	}
}
