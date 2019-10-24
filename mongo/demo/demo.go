package module

import (
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/mongo"
)

// DemoCollectionName 在mongo表名，约定首字母小写
const DemoCollectionName = "demo"

// Demo 用户消息
type Demo struct {
	mongo.PublicFields`bson:",inline"`

	UserID     string      `bson:"userID" json:"userID"`         // 用户id
	FileID     string      `bson:"fileID" json:"fileID"`         // 文件id
	TaskID     string      `bson:"taskID" json:"taskID"`         // 策略运行id，需要建立索引
	PolicyName string      `bson:"policyName" json:"policyName"` // 用户命名的策略名称
	Message    interface{} `bson:"message" json:"message"`       // 信息内容
}

// Insert 插入一条新的记录
func (object *Demo) Insert() (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()
	object.SetFieldsValue()
	return mconn.Insert(DemoCollectionName, object)
}

// FindDemo 获取单条记录
func FindDemo(selector bson.M, field bson.M) (*Demo, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &Demo{}
	return object, mconn.FindOne(DemoCollectionName, object, selector, field)
}

// FindDemos 获取多条记录
func FindDemos(selector bson.M, field bson.M, page int, limit int, sort ...string) ([]*Demo, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	// 默认从第一页开始
	if page < 1 {
		page = 1
	}

	objects := []*Demo{}
	return objects, mconn.FindAll(DemoCollectionName, &objects, selector, field, page-1, limit, sort...)
}

// UpdateDemo 更新单条记录
func UpdateDemo(selector, update bson.M) (err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateOne(DemoCollectionName, selector, mongo.UpdatedTime(update))
}

// UpdateDemos 更新多条记录
func UpdateDemos(selector, update bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(DemoCollectionName, selector, mongo.UpdatedTime(update))
}

// FindAndModifyDemo 更新并返回最新记录
func FindAndModifyDemo(selector bson.M, update bson.M) (*Demo, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	object := &Demo{}
	return object, mconn.FindAndModify(DemoCollectionName, object, selector, mongo.UpdatedTime(update))
}

// CountDemos 统计数量，不包括删除记录
func CountDemos(selector bson.M) (n int, err error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.Count(DemoCollectionName, mongo.ExcludeDeleted(selector))
}

// DeleteDemo 删除记录
func DeleteDemo(selector bson.M) (int, error) {
	mconn := mongo.GetSession()
	defer mconn.Close()

	return mconn.UpdateAll(DemoCollectionName, selector, mongo.DeletedTime(bson.M{}))
}
