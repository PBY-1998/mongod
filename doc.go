/**
    @author: potten
    @since: 2022/8/5
    @desc: //TODO 通用Doc
**/
package mongod

type Doc struct {
	ID         string `bson:"_id,omitempty" json:"_id,omitempty" form:"_id"`
	CreateTime int64  `bson:"createTime,omitempty" json:"createTime,omitempty" form:"createTime"`
	UpdateTime int64  `bson:"updateTime,omitempty" json:"updateTime,omitempty" form:"updateTime"`
	IsDelete   bool   `bson:"isDelete" json:"isDelete" form:"isDelete"`
	DeleteTime int64  `bson:"deleteTime,omitempty" json:"deleteTime,omitempty" form:"deleteTime"`
}
