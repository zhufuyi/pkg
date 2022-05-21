package qmgo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"pkg/logger"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	defaultClient *DefaultClient

	// ErrNotFoundStr 错误
	ErrNotFoundStr = qmgo.ErrNoSuchDocuments.Error()
)

// Clienter 提供外部调用的接口
type Clienter interface {
	WithLog() Clienter
	Close() error
	Insert(collection string, doc interface{}) error
	InsertMany(collection string, docs []interface{}) error
	FindOne(collection string, result interface{}, query bson.M, fields bson.M) error
	FindAll(collection string, result interface{}, query bson.M, fields bson.M, skip int64, limit int64, sort ...string) error
	UpdateOne(collection string, query bson.M, update bson.M) error
	UpdateAll(collection string, query bson.M, update bson.M) (int64, error)
	DeleteOne(collection string, query bson.M) error
	DeleteAll(collection string, query bson.M) (int64, error)
	DeleteOneCompletely(collection string, query bson.M) error
	DeleteAllCompletely(collection string, query bson.M) (int64, error)
	Count(collection string, query bson.M) (int64, error)
	CountAll(collection string, query bson.M) (int64, error)
	EnsureIndexes(collection string, indexes []string) error
	EnsureUniques(collection string, uniques []string) error
	Aggregate(collection string, matchStage primitive.D, groupStage primitive.D) ([]bson.M, error)
	Transaction(collection string, callback func(sessCtx context.Context) (interface{}, error)) (interface{}, error)
	GetQmgoClient(collection string) *qmgo.QmgoClient
}

// GetCli 获取新的session，使用前先调用函数InitToGlobal，否则抛出异常
func GetCli() Clienter {
	if defaultClient == nil {
		panic("mongodb session is nil, go to connect mongodb first, eg: mongo.InitToGlobal(uri)")
	}
	return &DefaultClient{
		ctx:      defaultClient.ctx,
		mdbName:  defaultClient.mdbName,
		cli:      defaultClient.cli,
		printLog: defaultClient.printLog,
	}
	//return defaultClient
}

// InitToGlobal 初始化mongodb，全局使用，只适用对单个mongodb操作，如果不指定数据库，默认数据库名为test
// 		形式一：localhost 或 localhost:27017
// 		形式二：mongodb://localhost:27017/database_name 或 mongodb://localhost1:port,localhost2:port/database_name
// 		形式三：mongodb://user:password@localhost:27017/database_name 或 mongodb://user:password@localhost1:port,localhost2:port/database_name
func InitToGlobal(url string) error {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: url})
	if err != nil {
		return err
	}

	dbName := getDBName(url)

	defaultClient = &DefaultClient{
		ctx:     ctx,
		mdbName: dbName,
		cli:     client,
	}

	return client.Ping(10)
}

// Init 初始化mongodb，使用对多个mongodb独立操作，如果不指定数据库，默认数据库名为test
func Init(url string) (string, Clienter, error) {
	ctx := context.Background()
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: url})
	if err != nil {
		return "", nil, err
	}

	dbName := getDBName(url)

	dc := &DefaultClient{
		ctx:     ctx,
		mdbName: dbName,
		cli:     client,
	}

	return dbName, dc, client.Ping(10)
}

// -------------------------------------------------------------------------------------------------

// DefaultClient mgo的客户端对象
type DefaultClient struct {
	ctx      context.Context // 上下文
	mdbName  string          // 数据库名称
	cli      *qmgo.Client    // 客户端
	printLog bool            // 是否打印执行命令信息
}

// WithLog 打印执行命令
func (d *DefaultClient) WithLog() Clienter {
	d.printLog = true
	return d
}

// Close 关闭连接
func (d *DefaultClient) Close() error {
	return d.cli.Close(d.ctx)
}

// Insert 插入一条新数据
func (d *DefaultClient) Insert(collection string, doc interface{}) error {
	_, err := d.cli.Database(d.mdbName).Collection(collection).InsertOne(d.ctx, doc)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.insert", d.mdbName, collection)),
				logger.Any("content", doc),
			).Warn("mongodb insert failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.insert", d.mdbName, collection)),
			logger.Any("content", doc),
		).Info("mongodb insert success.")
	}

	return nil
}

// InsertMany 插入多条新数据
func (d *DefaultClient) InsertMany(collection string, docs []interface{}) error {
	_, err := d.cli.Database(d.mdbName).Collection(collection).InsertMany(d.ctx, docs)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.insert", d.mdbName, collection)),
				logger.Any("content", docs),
			).Warn("mongodb insert failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.insert", d.mdbName, collection)),
			logger.Any("content", docs),
		).Info("mongodb insert success.")
	}

	return nil
}

// FindOne 查找一条记录
func (d *DefaultClient) FindOne(collection string, result interface{}, query bson.M, fields bson.M) error {
	err := d.cli.Database(d.mdbName).Collection(collection).Find(d.ctx, ExcludeDeleted(query)).Select(fields).One(result)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.findOne", d.mdbName, collection)),
				logger.Any("query", query),
				logger.Any("fields", fields),
			).Warn("mongodb find one failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.findOne", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Any("fields", fields),
			logger.Any("result", result),
		).Info("mongodb find one success.")
	}

	return nil
}

// FindAll 查找符合条件的多条记录，这里参数skip表示页码，limit表示每页多少行数据
func (d *DefaultClient) FindAll(collection string, result interface{}, query bson.M, fields bson.M, page int64, limit int64, sort ...string) error {
	err := d.cli.Database(d.mdbName).Collection(collection).Find(d.ctx, ExcludeDeleted(query)).Select(fields).Sort(sort...).Skip(page * limit).Limit(limit).All(result)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.find", d.mdbName, collection)),
				logger.Any("query", query),
				logger.Any("fields", fields),
				logger.Int64("page", page),
				logger.Int64("limit", limit),
				logger.Any("sort", sort),
			).Warn("mongodb find all failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.find", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Any("fields", fields),
			logger.Int64("page", page),
			logger.Int64("limit", limit),
			logger.Any("sort", sort),
			logger.Any("result", result),
		).Info("mongodb find all success.")
	}

	return nil
}

// UpdateOne 更新一条记录
func (d *DefaultClient) UpdateOne(collection string, query bson.M, update bson.M) error {
	err := CheckUpdateContent(update)
	if err != nil {
		return err
	}

	err = d.cli.Database(d.mdbName).Collection(collection).UpdateOne(d.ctx, ExcludeDeleted(query), UpdatedTime(update))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.update", d.mdbName, collection)),
				logger.Any("query", query),
				logger.Any("update", update),
			).Warn("mongodb update one failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.update", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Any("update", update),
		).Info("mongodb update one success.")
	}

	return nil
}

// UpdateAll 更新所有匹配的记录
func (d *DefaultClient) UpdateAll(collection string, query bson.M, update bson.M) (int64, error) {
	err := CheckUpdateContent(update)
	if err != nil {
		return 0, err
	}

	updateResult, err := d.cli.Database(d.mdbName).Collection(collection).UpdateAll(d.ctx, ExcludeDeleted(query), UpdatedTime(update))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.updateAll", d.mdbName, collection)),
				logger.Any("query", query),
				logger.Any("update", update),
			).Warn("mongodb update all failed.")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.updateAll", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Any("update", update),
			logger.Int64("updatedNum", updateResult.ModifiedCount),
		).Info("mongodb update all success.")
	}

	return updateResult.ModifiedCount, nil
}

// DeleteOne 标记性删除一条记录
func (d *DefaultClient) DeleteOne(collection string, query bson.M) error {
	err := d.cli.Database(d.mdbName).Collection(collection).UpdateOne(d.ctx, ExcludeDeleted(query), DeletedTime(bson.M{}))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.update", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb markup delete one failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
			logger.Any("query", query),
		).Info("mongodb markup delete one success.")
	}

	return nil
}

// DeleteAll 标记性删除所有匹配的记录
func (d *DefaultClient) DeleteAll(collection string, query bson.M) (int64, error) {
	updateResult, err := d.cli.Database(d.mdbName).Collection(collection).UpdateAll(d.ctx, ExcludeDeleted(query), DeletedTime(bson.M{}))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb markup delete all failed.")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Int64("removedNum", updateResult.ModifiedCount),
		).Info("mongodb markup delete all success.")
	}

	return updateResult.ModifiedCount, nil
}

// DeleteOneReal 删除一条记录，包括标记性删除的记录
func (d *DefaultClient) DeleteOneCompletely(collection string, query bson.M) error {
	err := d.cli.Database(d.mdbName).Collection(collection).Remove(d.ctx, query)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb delete one failed.")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
			logger.Any("query", query),
		).Info("mongodb delete one success.")
	}

	return nil
}

// DeleteAllReal 删除所有匹配的记录，包括标记性删除的记录
func (d *DefaultClient) DeleteAllCompletely(collection string, query bson.M) (int64, error) {
	deleteResult, err := d.cli.Database(d.mdbName).Collection(collection).DeleteAll(d.ctx, query)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb real delete all failed.")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.remove", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Int64("removedNum", deleteResult.DeletedCount),
		).Info("mongodb real delete all success.")
	}

	return deleteResult.DeletedCount, nil
}

// Count 统计匹配的数量，不包括标记性删除的记录
func (d *DefaultClient) Count(collection string, query bson.M) (int64, error) {
	count, err := d.cli.Database(d.mdbName).Collection(collection).Find(d.ctx, ExcludeDeleted(query)).Count()
	if err != nil && err.Error() != ErrNotFoundStr {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb count failed.")
		}

		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Int64("count", count),
		).Info("mongodb count success.")
	}

	return count, nil
}

// CountAll 统计匹配的数量，包括标记性删除的记录
func (d *DefaultClient) CountAll(collection string, query bson.M) (int64, error) {
	count, err := d.cli.Database(d.mdbName).Collection(collection).Find(d.ctx, query).Count()
	if err != nil && err.Error() != ErrNotFoundStr {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
				logger.Any("query", query),
			).Warn("mongodb count all failed.")
		}

		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
			logger.Any("query", query),
			logger.Int64("count", count),
		).Info("mongodb count all success.")
	}

	return count, nil
}

// EnsureIndexes 设置普通索引
func (d *DefaultClient) EnsureIndexes(collection string, indexes []string) error {
	d.cli.Database(d.mdbName).Collection(collection).EnsureIndexes(d.ctx, []string{}, indexes)

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.ensureIndexes", d.mdbName, collection)),
			logger.Any("indexes", indexes),
		).Info("mongodb set index key success.")
	}

	return nil
}

// EnsureUniques 设置唯一索引
func (d *DefaultClient) EnsureUniques(collection string, uniques []string) error {
	d.cli.Database(d.mdbName).Collection(collection).EnsureIndexes(d.ctx, uniques, []string{})

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.ensureUniques", d.mdbName, collection)),
			logger.Any("uniques", uniques),
		).Info("mongodb set uniques key success.")
	}

	return nil
}

// Aggregate 聚合查询
func (d *DefaultClient) Aggregate(collection string, matchStage primitive.D, groupStage primitive.D) ([]bson.M, error) {
	var showsWithInfo []bson.M
	err := d.cli.Database(d.mdbName).Collection(collection).Aggregate(d.ctx, qmgo.Pipeline{matchStage, groupStage}).All(&showsWithInfo)
	if err != nil && err.Error() != ErrNotFoundStr {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
				logger.Any("matchStage", matchStage),
				logger.Any("groupStage", groupStage),
			).Warn("mongodb aggregate failed.")
		}

		return nil, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("%s.%s.count", d.mdbName, collection)),
			logger.Any("matchStage", matchStage),
			logger.Any("groupStage", groupStage),
			logger.Any("showsWithInfo", showsWithInfo),
		).Info("mongodb aggregate success.")
	}

	return showsWithInfo, nil
}

func (d *DefaultClient) GetQmgoClient(collection string) *qmgo.QmgoClient {
	db := d.cli.Database(d.mdbName)
	coll := db.Collection(collection)
	cli := &qmgo.QmgoClient{
		coll, db, d.cli,
	}

	return cli
}

// Transaction 事务
func (d *DefaultClient) Transaction(collection string, callback func(sessCtx context.Context) (interface{}, error)) (interface{}, error) {
	db := d.cli.Database(d.mdbName)
	coll := db.Collection(collection)
	cli := &qmgo.QmgoClient{
		coll, db, d.cli,
	}

	return cli.DoTransaction(d.ctx, callback)
}

func getDBName(url string) string {
	url = strings.Trim(url, "/")

	matches := regexp.MustCompile(`mongodb://[^/]+/(.*)`).FindAllSubmatch([]byte(url), -1)

	for _, match := range matches {
		return string(match[1])
	}

	return "test"
}
