/**
    @author: potten
    @since: 2022/8/8
    @desc: //TODO 通用类
**/
package mongo_go

import (
	"bytes"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

//parseAnyOptions 解析传入的options
func parseAnyOptions[T interface{}](anyOptions ...interface{}) (ctx context.Context, opts []T) {
	ctx = context.TODO()
	for _, option := range anyOptions {
		switch value := option.(type) {
		case context.Context:
			ctx = value
		case T:
			opts = append(opts, value)
		}
	}
	return ctx, opts
}

//mergeMap 合并多个map
func mergeMap[M ~map[K]V, K comparable, V any](ms ...M) M {
	newMap := make(M)
	for _, m := range ms {
		for k, v := range m {
			newMap[k] = v
		}
	}

	return newMap
}

//turnMap struct转为map slice
func turnMap[M ~map[string]interface{} | ~[]map[string]interface{}, K comparable, V any](v V) (M, error) {
	data := new(M)
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	// 防止unit64 和 float64 类型都转为 float64 类型
	decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
	decoder.UseNumber()
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return *data, nil
}

func newStringID() string {
	objectId := primitive.NewObjectID()
	return primitive.ObjectID.Hex(objectId)
}

func newInsertCommonDoc[M ~map[K]V, K comparable, V any]() (M, error) {
	newMap := make(M)
	doc := &Doc{
		ID:         newStringID(),
		CreateTime: time.Now().Unix(),
		UpdateTime: time.Now().Unix(),
		IsDelete:   false,
	}
	marshal, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	_err := bson.Unmarshal(marshal, &newMap)
	if err != nil {
		return nil, _err
	}
	return newMap, nil
}

func newDeleteSoftCommonDoc[M ~map[K]V, K comparable, V any]() (M, error) {
	newMap := make(M)
	doc := &Doc{
		IsDelete:   true,
		DeleteTime: time.Now().Unix(),
	}
	marshal, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	_err := bson.Unmarshal(marshal, &newMap)
	if err != nil {
		return nil, _err
	}
	return newMap, nil
}

func newUpdateCommonDoc[M ~map[K]V, K comparable, V any]() (M, error) {
	newMap := make(M)
	doc := &Doc{
		UpdateTime: time.Now().Unix(),
	}
	marshal, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	_err := bson.Unmarshal(marshal, &newMap)
	if err != nil {
		return nil, _err
	}
	return newMap, nil
}

func newUpsertCommonDoc[M ~map[K]V, K comparable, V any]() (M, error) {
	newMap := make(M)
	doc := &Doc{
		ID:         newStringID(),
		CreateTime: time.Now().Unix(),
	}
	marshal, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	_err := bson.Unmarshal(marshal, &newMap)
	if err != nil {
		return nil, _err
	}
	return newMap, nil
}
