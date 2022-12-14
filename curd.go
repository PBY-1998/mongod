/**
    @author: potten
    @since: 2022/8/5
    @desc: //TODO 封装go语言mongo curd操作
**/
package mongod

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Curd 调用自定义mongo curd接口方法
type Curd struct {
	Coll *mongo.Collection
}

//AbstractMongoCurd 自封装mongo curd接口方法
type AbstractMongoCurd interface {
	Insert(filter interface{}, multi bool, opts ...interface{}) (bool, error)
	DeleteHard(filter interface{}, multi bool, opts ...interface{}) (bool, error)
	DeleteSoft(filter interface{}, multi bool, opts ...interface{}) (bool, error)
	Query(filter interface{}, opts ...interface{}) (interface{}, error)
	Update(filter interface{}, update interface{}, multi bool, opts ...interface{}) (bool, error)
	Upsert(filter interface{}, update interface{}, multi bool, opts ...interface{}) (bool, error)
	Count(filter interface{}, opts ...interface{}) (int64, error)
}

// Insert 插入文档
//
// @param filter interface{} 过滤数据
//
// @param multi bool 是否批量
//
// @param opts ...interface{} 插入选项（可选）
//
// @return bool error
func (c *Curd) Insert(filter interface{}, multi bool, opts ...interface{}) (bool, error) {
	switch multi {
	case true:
		var docs []interface{}
		ctx, _opts := parseAnyOptions[*options.InsertManyOptions](opts...)
		insertDocs, _ := turnMap[[]map[string]interface{}, string, interface{}](filter)
		for _, insertDoc := range insertDocs {
			commonDoc, _ := newInsertCommonDoc[map[string]interface{}, string, interface{}]()
			doc := mergeMap[map[string]interface{}, string, interface{}](commonDoc, insertDoc)
			docs = append(docs, doc)
		}
		many, err := c.Coll.InsertMany(ctx, docs, _opts...)
		if many.InsertedIDs != nil {
			return true, nil
		}

		return false, err
	case false:
		ctx, _opts := parseAnyOptions[*options.InsertOneOptions](opts...)
		commonDoc, _ := newInsertCommonDoc[map[string]interface{}, string, interface{}]()
		insertDoc, _ := turnMap[map[string]interface{}, string, interface{}](filter)
		doc := mergeMap[map[string]interface{}, string, interface{}](commonDoc, insertDoc)
		one, err := c.Coll.InsertOne(ctx, doc, _opts...)

		if one.InsertedID != nil {
			return true, nil
		}

		return false, err
	}

	return false, errors.New("insert interface unknown error")
}

// DeleteHard 硬删除文档
//
// @param filter interface{} 过滤数据
//
// @param multi bool 是否批量
//
// @param opts ...interface{} 删除选项（可选）
//
// @return bool error
func (c *Curd) DeleteHard(filter interface{}, multi bool, opts ...interface{}) (bool, error) {
	ctx, _opts := parseAnyOptions[*options.DeleteOptions](opts...)
	switch multi {
	case true:
		deleteHardDocs, _ := turnMap[map[string]interface{}, string, interface{}](filter)
		one, err := c.Coll.DeleteMany(ctx, deleteHardDocs, _opts...)
		if one.DeletedCount > 0 {
			return true, nil
		}
		return false, err
	case false:
		deleteHardDoc, _ := turnMap[map[string]interface{}, string, interface{}](filter)
		one, err := c.Coll.DeleteOne(ctx, deleteHardDoc, _opts...)
		if one.DeletedCount > 0 {
			return true, nil
		}
		return false, err
	}

	return false, errors.New("delete hard interface unknown error")
}

// DeleteSoft 软删除文档
//
// @param filter interface{} 过滤数据
//
// @param multi bool 是否批量
//
// @param opts ...interface{} 删除选项（可选）
//
// @return bool error
func (c *Curd) DeleteSoft(filter interface{}, multi bool, opts ...interface{}) (bool, error) {
	ctx, _opts := parseAnyOptions[*options.UpdateOptions](opts...)
	deleteSoftCommonDoc, _ := newDeleteSoftCommonDoc[map[string]interface{}, string, interface{}]()
	switch multi {
	case true:
		deleteSoftDocs, _ := turnMap[[]map[string]interface{}, string, interface{}](filter)
		many, err := c.Coll.UpdateMany(ctx, deleteSoftDocs, bson.D{{"$set", deleteSoftCommonDoc}}, _opts...)
		if many.MatchedCount > 0 {
			return true, nil
		}
		return false, err
	case false:
		deleteSoftDoc, _ := turnMap[map[string]interface{}, string, interface{}](filter)
		one, err := c.Coll.UpdateOne(ctx, deleteSoftDoc, bson.D{{"$set", deleteSoftCommonDoc}}, _opts...)
		if one.MatchedCount > 0 {
			return true, nil
		}
		return false, err
	}
	return false, errors.New("delete soft interface unknown error")
}

// Query 查询
//
// @param filter interface{} 过滤数据
//
// @param opts ...interface{} 查找选项（可选）
//
// @return interface{} error
func (c *Curd) Query(filter interface{}, opts ...interface{}) (interface{}, error) {
	var res []bson.M
	ctx, _opts := parseAnyOptions[*options.FindOptions](opts...)
	find, _ := c.Coll.Find(ctx, filter, _opts...)
	if err := find.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// Update 修改文档
//
// @param filter interface{} 过滤数据
//
// @param update interface{} 更改文档
//
// @param multi bool 是否批量
//
// @param opts ...interface{} 更改选项（可选）
//
// @return bool error
func (c *Curd) Update(filter interface{}, update interface{}, multi bool, opts ...interface{}) (bool, error) {
	ctx, _opts := parseAnyOptions[*options.UpdateOptions](opts...)
	updateCommonDoc, _ := newUpdateCommonDoc[map[string]interface{}, string, interface{}]()
	updateMap, _ := turnMap[map[string]interface{}, string, interface{}](update)
	updateDoc := mergeMap[map[string]interface{}, string, interface{}](updateCommonDoc, updateMap)
	switch multi {
	case true:
		many, err := c.Coll.UpdateMany(ctx, filter, bson.D{{"$set", updateDoc}}, _opts...)
		if many.MatchedCount+many.ModifiedCount > 0 {
			return true, nil
		}
		return false, err
	case false:
		one, err := c.Coll.UpdateOne(ctx, filter, bson.D{{"$set", updateDoc}}, _opts...)
		if one.MatchedCount+one.ModifiedCount > 0 {
			return true, nil
		}
		return false, err
	}
	return false, errors.New("update interface unknown error")
}

// Upsert 更新插入
//
// @param filter interface{} 过滤数据
//
// @param update interface{} 更改文档
//
// @param multi bool 是否批量
//
// @param opts ...interface{} 更改选项（可选）
//
// @return bool error
func (c *Curd) Upsert(filter interface{}, update interface{}, multi bool, opts ...interface{}) (bool, error) {
	upsert := options.Update().SetUpsert(true)
	ctx, _opts := parseAnyOptions[*options.UpdateOptions](opts, upsert)
	upsertCommonDoc, _ := newUpsertCommonDoc[map[string]interface{}, string, interface{}]()
	updateCommonDoc, _ := newUpdateCommonDoc[map[string]interface{}, string, interface{}]()
	updateMap, _ := turnMap[map[string]interface{}, string, interface{}](update)
	updateDoc := mergeMap[map[string]interface{}, string, interface{}](updateMap, updateCommonDoc)

	switch multi {
	case true:
		many, err := c.Coll.UpdateMany(ctx, filter, bson.D{{"$set", update}}, _opts...)
		if many.MatchedCount+many.UpsertedCount > 0 {
			return true, nil
		}
		return false, err
	case false:
		one, err := c.Coll.UpdateOne(ctx, filter, bson.D{{"$set", updateDoc}, {"$setOnInsert", upsertCommonDoc}}, _opts...)
		if one.ModifiedCount+one.UpsertedCount > 0 {
			return true, nil
		}
		return false, err
	}
	return false, errors.New("upsert interface unknown error")
}

func (c *Curd) Count(filter interface{}, opts ...interface{}) (int64, error) {
	ctx, _opts := parseAnyOptions[*options.CountOptions](opts)
	if documents, err := c.Coll.CountDocuments(ctx, filter, _opts...); err != nil {
		return -1, err
	} else {
		return documents, nil
	}
}
