package mongo

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/k0kubun/pp"
	"github.com/zhufuyi/pkg/krand"
)

var (
	server = "mongodb://collectdata:123456@192.168.101.88:27018/collectdata"
	//server = "mongodb://vison:123456@192.168.8.200:27017/test1"
	srcURL = "mongodb://vison:123456@192.168.8.200:27017/test2"
	dstURL = "mongodb://vison:123456@192.168.8.200:27017/test3"

	srcDbName string
	dstDbName string

	srcSes *mgo.Session
	dstSes *mgo.Session

	err error
)

const testDataCollection = "testData"

type testData struct {
	PublicFields `bson:",inline"` // 公共字段，id和时间

	Name string `bson:"name" json:"name"`
	Age  int    `bson:"age" json:"age"`
}

// 初始化mongodb
func init() {
	err := InitializeMongodb(server)
	if err != nil {
		panic(err)
	}
}

func TestSingleInit(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	count, err := mconn.Count(testDataCollection, bson.M{})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("count", count)
}

func TestMultiInit(t *testing.T) {
	srcDbName, srcSes, err = InitMongo(srcURL)
	if err != nil {
		t.Error(err)
		return
	}

	dstDbName, dstSes, err = InitMongo(dstURL)
	if err != nil {
		t.Error(err)
		return
	}

	mconn := Clone(srcDbName, srcSes)
	defer mconn.Close()
	nSrc, err := mconn.Count(testDataCollection, bson.M{})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("src", nSrc)

	mconn = Clone(dstDbName, dstSes)
	defer mconn.Close()
	nDst, err := mconn.Count(testDataCollection, bson.M{})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("dst", nDst)
}

func TestClone(t *testing.T) {
	name, ses, err := InitMongo(server)
	if err != nil {
		t.Error(err)
		return
	}

	mconn := Clone(name, ses)
	defer mconn.Close()

	n, err := mconn.Count(testDataCollection, bson.M{"age": 12})
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(n)
}

// 测试插入和查找
func TestInsertAndFind(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	// 测试数据
	idString := bson.NewObjectId().Hex()
	name := "zhangsan_" + idString[len(idString)-5:]
	expected := &testData{Name: name, Age: 12}
	expected.SetFieldsValue() // 设置公共字段id和时间值

	// 测试插入数据
	err := mconn.WithLog().Insert(testDataCollection, expected)
	if err != nil {
		t.Error(err)
		return
	}

	// 测试查找一条数据
	actual := &testData{}
	err = mconn.WithLog().FindOne(testDataCollection, actual, bson.M{"name": name}, nil)
	if err != nil {
		t.Error(err)
	}
	if actual.Name != expected.Name {
		t.Errorf("got %s, expected %s\n", actual.Name, expected.Name)
	}

	// 查找所有数据
	tds := []testData{}
	err = mconn.WithLog().FindAll(testDataCollection, &tds, bson.M{"age": 12}, nil, 0, 0, "-_id")
	if err != nil {
		t.Error(err)
		return
	}
	if len(tds) > 0 && tds[0].Age != expected.Age {
		t.Errorf("got %d, expected %d\n", tds[0].Age, expected.Age)
	}
}

// 测试更新
func TestUpdate(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	selector := bson.M{"age": 12}
	update := bson.M{"$set": bson.M{"age": 22}}
	err := mconn.WithLog().UpdateOne(testDataCollection, selector, update)
	if err != nil && err != mgo.ErrNotFound {
		t.Error(err)
	}

	selector = bson.M{"age": 12}
	update = bson.M{"$set": bson.M{"age": 10}}
	_, err = mconn.WithLog().UpdateAll(testDataCollection, selector, update)
	if err != nil {
		t.Error(err)
	}
}

// 测试删除
func TestDeleteAndCount(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	// 测试标记性删除一条记录
	selector := bson.M{"age": 12}
	err := mconn.WithLog().DeleteOne(testDataCollection, selector)
	if err != nil && err != mgo.ErrNotFound {
		t.Error(err)
	}

	// 测试标记性删除所有记录
	selector = bson.M{"age": 12}
	_, err = mconn.WithLog().DeleteAll(testDataCollection, selector)
	if err != nil {
		t.Error(err)
	}

	// 测试统计不包括标记删除数量
	selector = bson.M{"age": 12}
	count, err := mconn.WithLog().Count(testDataCollection, selector)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Errorf("got %d, expected %d", count, 0)
	}

	// 测试统计包括标记删除数量
	selector = bson.M{"age": 12}
	count, err = mconn.WithLog().CountAll(testDataCollection, selector)
	if err != nil {
		t.Error(err)
	}
	if count == 0 {
		t.Error("got 0, expected > 0")
	}

	// 测试真实删除一条记录
	selector = bson.M{"age": 12}
	err = mconn.WithLog().DeleteOneReal(testDataCollection, selector)
	if err != nil && err != mgo.ErrNotFound {
		t.Error(err)
	}

	// 测试真实删除所有匹配记录
	selector = bson.M{"age": 12}
	_, err = mconn.WithLog().DeleteAllReal(testDataCollection, selector)
	if err != nil && err != mgo.ErrNotFound {
		t.Error(err)
	}
}

// 测试更新并返回最新记录
func TestFindAndModify(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	selector := bson.M{"age": 12}
	update := bson.M{"$set": bson.M{"name": "FindAndModify"}}
	result := &testData{}
	err := mconn.WithLog().FindAndModify(testDataCollection, result, selector, update)
	if err != nil && err != mgo.ErrNotFound {
		t.Error(err)
	}
}

// 测试索引设置
func TestIndex(t *testing.T) {
	mconn := GetSession()
	defer mconn.Close()

	err := mconn.EnsureIndexKey(testDataCollection, "age")
	if err != nil {
		t.Error(err)
	}

	err = mconn.EnsureIndex(testDataCollection, mgo.Index{Key: []string{"name"}, Unique: true})
	if err != nil {
		t.Error(err)
	}

	td1 := &testData{Name: "zhangsan", Age: 12}
	td1.SetFieldsValue() // 设置公共字段id和时间值
	mconn.Insert(testDataCollection, td1)

	td2 := &testData{Name: "zhangsan", Age: 13}
	td2.SetFieldsValue() // 设置公共字段id和时间值
	err = mconn.Insert(testDataCollection, td2)
	if err == nil {
		t.Error("index name unique failed")
		return
	}

	fmt.Println(err)
}

// 测试并发插入数据
func TestBenchInsert(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup
	var start = time.Now()

	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mconn := GetSession()
			defer mconn.Close()

			td := testData{Name: fmt.Sprintf("zhansan_%d", i), Age: randAge()}
			td.SetFieldsValue()
			err := mconn.Insert(testDataCollection, td)
			if err != nil {
				t.Error(err)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nwrite success count = %d, time = %s\n", successCount, time.Now().Sub(start))
}

type UserInfo struct {
	PublicFields `bson:",inline"` // 公共字段，id和时间

	Name   string `bson:"name"`
	Age    int    `bson:"age"`
	Weight int    `bson:"weight"`
}

// 测试并发插入数据
func TestBenchInsert2(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup
	var total = 5000
	var start = time.Now()

	for i := 0; i < total; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mconn := GetSession()
			defer mconn.Close()

			ui := &UserInfo{
				Name:   fmt.Sprintf("zhansan_%d", i),
				Age:    krand.Int(50) + 10,
				Weight: krand.Int(50) + 40,
			}
			ui.SetFieldsValue()

			err := mconn.Insert("user", ui)
			if err != nil {
				t.Error(err)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nwrite success count = %d, total = %d, time = %s\n", successCount, total, time.Now().Sub(start))
	// write success count = 5000, total = 5000, time = 606.3412ms
}

// 测试并发读取数据，建立索引前后并发查询速度相差10倍左右
func TestBenchRead(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup
	var total = 5000
	var start = time.Now()

	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			mconn := GetSession()
			defer mconn.Close()

			query := bson.M{"name": fmt.Sprintf("zhansan_%d", i)}
			ui := &UserInfo{}
			err := mconn.FindOne("user", ui, query, nil)
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
	// find success count = 5000, total = 5000, time = 747.966ms
}

func randAge() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(99) + 1
}

type result struct {
	//ID    bson.ObjectId `json:"_id" bson:"_id"`
	//Total float64       `json:"total" bson:"total"`
	Price    float64 `json:"price" bson:"price"`       // 成交价格
	Quantity float64 `json:"quantity" bson:"quantity"` // 买卖数量
	Volume   float64 `json:"volume" bson:"volume"`     // 成交额
}

func TestDefaultSession_Run(t *testing.T) {
	//	cmd := `
	//db.gridTradeRecord.aggregate([
	//    {
	//        $match:{"gsid":ObjectId("5f3a2816c7ef825a00e96578"),"side":"sell","orderState":"filled"}
	//    },
	//    {
	//        $group: {
	//            _id: "$gsid",
	//            total: { $sum: "$quantity"}
	//        },
	//    }
	//])
	//`

	cmd2 := bson.M{"aggregate": bson.D{
		{"$match", bson.M{
			"gsid":       bson.ObjectIdHex("5f3a2816c7ef825a00e96578"),
			"side":       "sell",
			"orderState": "filled",
		}},
		{"$group", bson.M{
			"_id":   "$gsid",
			"total": bson.M{"$sum": "$quantity"},
		}},
	},
	}

	res := &result{}

	mconn := GetSession()
	defer mconn.Close()

	err := mconn.Run(cmd2, res)
	if err != nil {
		t.Error(err)
		return
	}

	pp.Println(res)
}
