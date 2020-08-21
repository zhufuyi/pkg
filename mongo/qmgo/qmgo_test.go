package qmgo

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"context"

	"github.com/k0kubun/pp"
	"github.com/zhufuyi/pkg/krand"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	uri            = "mongodb://collectdata:123456@192.168.101.88:27018/collectdata"
	collectionName = "user"
)

type UserInfo struct {
	PublicFields `bson:",inline"` // 公共字段，id和时间

	Name   string `bson:"name"`
	Gender int    `bson:"gender"`
	Age    int    `bson:"age"`
	Weight int    `bson:"weight"` // 性别 0:未知，1:男，2:女
}

func getGender() int {
	if krand.Int()%2 == 1 {
		return 1
	}

	return 2
}

func init() {
	err := InitToGlobal(uri)
	if err != nil {
		panic(err)
	}

	fmt.Println("init mongodb success")
}

func TestDefaultClient_Insert(t *testing.T) {
	ui := &UserInfo{
		Name:   string(krand.String(krand.R_All)),
		Age:    krand.Int(50) + 10,
		Weight: krand.Int(50) + 40,
		Gender: getGender(),
	}
	ui.SetFieldsValue()

	err := GetCli().WithLog().Insert(collectionName, ui)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDefaultClient_InsertMany(t *testing.T) {
	var uis []interface{}
	for i := 0; i < 10; i++ {
		ui := &UserInfo{
			Name:   string(krand.String(krand.R_All)),
			Age:    krand.Int(50) + 10,
			Weight: krand.Int(50) + 40,
			Gender: getGender(),
		}
		ui.SetFieldsValue()
		uis = append(uis, ui)
	}

	err := GetCli().WithLog().InsertMany(collectionName, uis)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDefaultClient_FindOne(t *testing.T) {
	query := bson.M{"_id": ObjectIDHex("5f3f2957ae5f3a33e26fa959")}
	//query = bson.M{"name": "iK11Mp"}
	fields := bson.M{}
	//fields = bson.M{"weight": true} // 指定输出的字段
	ui := &UserInfo{}

	err := GetCli().WithLog().FindOne(collectionName, ui, query, fields)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(ui)
}

func TestDefaultClient_FindAll(t *testing.T) {
	query := bson.M{}
	//query = bson.M{"age": bson.M{"$gt": 40}}
	fields := bson.M{}
	//fields = bson.M{"weight": true} // 指定输出的字段
	uis := []*UserInfo{}
	page := int64(0)
	limit := int64(20)
	sort := "-name"

	err := GetCli().WithLog().FindAll(collectionName, &uis, query, fields, page, limit, sort)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(uis)
}

func TestDefaultClient_UpdateOne(t *testing.T) {
	query := bson.M{"name": "iK11Mp"}
	update := bson.M{"$set": bson.M{"weight": 80}}

	err := GetCli().WithLog().UpdateOne(collectionName, query, update)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDefaultClient_UpdateAll(t *testing.T) {
	query := bson.M{"age": bson.M{"$gt": 40}}
	update := bson.M{"$set": bson.M{"weight": 75}}

	n, err := GetCli().WithLog().UpdateAll(collectionName, query, update)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("updated success count:", n)
}

func TestDefaultClient_Count(t *testing.T) {
	query := bson.M{"age": bson.M{"$gt": 40}}

	count, err := GetCli().Count(collectionName, query) // 不包括标记性删除
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("count =", count)
}

func TestDefaultClient_CountAll(t *testing.T) {
	query := bson.M{"age": bson.M{"$gt": 40}}
	count, err := GetCli().CountAll(collectionName, query) // 包括标记性删除
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("count =", count)
}

func TestDefaultClient_DeleteOne(t *testing.T) {
	query := bson.M{"_id": ObjectIDHex("5f3f2957ae5f3a33e26fa959")}
	err := GetCli().WithLog().DeleteOne(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	count, err := GetCli().Count(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	if count != 0 {
		t.Error("mark delete failed")
		return
	}

	fmt.Println("mark delete success")
}

func TestDefaultClient_DeleteAll(t *testing.T) {
	query := bson.M{"age": bson.M{"$gt": 40}}
	delCount, err := GetCli().WithLog().DeleteAll(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("delete count =", delCount)

	count, err := GetCli().Count(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	if count != 0 {
		t.Error("mark delete all failed")
		return
	}

	fmt.Println("mark delete success")
}

func TestDefaultClient_DeleteOneCompletely(t *testing.T) {
	query := bson.M{"_id": ObjectIDHex("5f3f2957ae5f3a33e26fa959")}
	err := GetCli().WithLog().DeleteOneCompletely(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	count, err := GetCli().CountAll(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	if count != 0 {
		t.Error("completely delete failed")
		return
	}

	fmt.Println("completely delete success")
}

func TestDefaultClient_DeleteAllCompletely(t *testing.T) {
	query := bson.M{"age": bson.M{"$gt": 40}}
	delCount, err := GetCli().WithLog().DeleteAllCompletely(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("delete count =", delCount)

	count, err := GetCli().CountAll(collectionName, query)
	if err != nil {
		t.Error(err)
		return
	}

	if count != 0 {
		t.Error("completely delete all failed")
		return
	}

	fmt.Println("completely delete success")
}

func TestDefaultClient_Index(t *testing.T) {
	GetCli().EnsureIndexes(collectionName, []string{"age"})
	GetCli().EnsureUniques(collectionName, []string{"name"})
}

func TestDefaultClient_Aggregate(t *testing.T) {
	// 筛选条件：age<20，聚合：以gender字段分组，统计各组的体重总和
	/*
		db.user.aggregate([
			{
				$match: {"age": {$lt: 20}}
			},
			{
				$group: {_id: "$gender", total: { $sum: "$weight"}},
			}
		])
	*/
	matchStage := bson.D{{"$match", []bson.E{{"age", bson.D{{"$lt", 20}}}}}}
	groupStage := bson.D{{"$group", bson.D{{"_id", "$gender"}, {"total", bson.D{{"$sum", "$weight"}}}}}}

	result, err := GetCli().WithLog().Aggregate(collectionName, matchStage, groupStage)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(result)
}

func TestDefaultClient_Transaction(t *testing.T) {
	cli := GetCli().GetQClient(collectionName)

	callback := func(sessCtx context.Context) (interface{}, error) {
		ui := &UserInfo{
			Name:   string(krand.String(krand.R_All)),
			Age:    krand.Int(50) + 10,
			Weight: krand.Int(50) + 40,
			Gender: getGender(),
		}
		ui.SetFieldsValue()
		if _, err := cli.InsertOne(sessCtx, ui); err != nil {
			return nil, err
		}

		query := bson.M{"name": "oaxl59"}
		update := bson.M{"$set": bson.M{"weight": 80}}
		if err := cli.UpdateOne(sessCtx, query, update); err != nil {
			return nil, err
		}

		return nil, nil
	}

	val, err := GetCli().Transaction(collectionName, callback)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(val)
}

// 测试并发插入数据
func TestBenchInsert(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup
	var total = 5000
	var start = time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			ui := &UserInfo{
				Name:   fmt.Sprintf("wangwu_%d", i),
				Age:    krand.Int(50) + 10,
				Weight: krand.Int(50) + 40,
			}
			ui.SetFieldsValue()

			err := GetCli().Insert(collectionName, ui)
			if err != nil {
				t.Error(err)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nwrite success count = %d, total =%d, time = %s\n", successCount, total, time.Now().Sub(start))
	// write success count = 5000, total =5000, time = 1.1688409s
}

// 测试并发读取数据，建立索引前后并发查询速度相差10倍左右
func TestBenchRead(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup
	var total = 5000
	var start = time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			query := bson.M{"name": fmt.Sprintf("wangwu_%d", i)}
			ui := &UserInfo{}
			err := defaultClient.FindOne(collectionName, ui, query, bson.M{})
			if err != nil {
				t.Error(err)
				return
			}

			if ui.Age < 1 {
				t.Errorf("got %d, expected >0", ui.Age)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nfind success count = %d, total = %d, time = %s\n", successCount, total, time.Now().Sub(start))
	// find success count = 5000, total = 5000, time = 1.6625511s
}
