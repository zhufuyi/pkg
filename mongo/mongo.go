package mongo

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/zhufuyi/logger"
)

var session *mgo.Session
var dbName string

// ErrNotFound 错误
var ErrNotFound = mgo.ErrNotFound

// GetSession 获取新的session，使用前先调用函数InitializeMongodb初始化session，否则抛出异常
func GetSession() Sessioner {
	if session == nil {
		panic("mongodb session is nil, go to connect mongodb first, eg: mongo.InitializeMongodb(url)")
	}

	return &DefaultSession{newSession: session.Clone()} // 并发时共用同一个socket，可能会造成goroutine的等待
	//return &DefaultSession{newSession: session.Copy()} // 复制一个新的session，必须手动close，经过并发压测，性能不如session.Clone()好
}

// InitializeMongodb 初始化mongodb，如果不指定数据库，默认数据库名为test
// 		形式一：localhost 或 localhost:27017
// 		形式二：mongodb://localhost:27017/database_name 或 mongodb://localhost1:port,localhost2:port/database_name
// 		形式三：mongodb://user:password@localhost:27017/database_name 或 mongodb://user:password@localhost1:port,localhost2:port/database_name
func InitializeMongodb(url string) error {
	var err error
	session, err = mgo.Dial(url)
	if err != nil {
		return err
	}

	dbName = getDatabaseName(url)

	session.SetMode(mgo.Monotonic, true)
	session.SetPoolLimit(1024)

	return session.Ping()
}

func getDatabaseName(url string) string {
	url = strings.Trim(url, "/")

	matches := regexp.MustCompile(`mongodb://[^/]+/(.*)`).FindAllSubmatch([]byte(url), -1)

	for _, match := range matches {
		return string(match[1])
	}

	return "test"
}

// Sessioner 提供外部调用的接口
type Sessioner interface {
	WithLog() Sessioner
	Close()
	Insert(collection string, a interface{}) error
	FindOne(collection string, result interface{}, selector bson.M, fields bson.M) error
	FindAll(collection string, result interface{}, selector bson.M, fields bson.M, skip int, limit int, sort ...string) error
	UpdateOne(collection string, selector bson.M, update bson.M) error
	UpdateAll(collection string, selector bson.M, update bson.M) (int, error)
	DeleteOne(collection string, selector bson.M) error
	DeleteAll(collection string, selector bson.M) (int, error)
	DeleteOneReal(collection string, selector bson.M) error
	DeleteAllReal(collection string, selector bson.M) (int, error)
	Count(collection string, selector bson.M) (int, error)
	CountAll(collection string, selector bson.M) (int, error)
	FindAndModify(collection string, result interface{}, selector bson.M, update bson.M) error
	EnsureIndexKey(collection string, indexKeys ...string) error
	EnsureIndex(collection string, index mgo.Index) error
}

// DefaultSession 默认mgo的会话
type DefaultSession struct {
	newSession *mgo.Session
	printLog   bool
}

// WithLog 带mongo执行命令log输出
func (d *DefaultSession) WithLog() Sessioner {
	d.printLog = true
	return d
}

// Close 关闭连接
func (d *DefaultSession) Close() {
	d.newSession.Close()
}

// Insert 插入一条新数据
func (d *DefaultSession) Insert(collection string, a interface{}) error {
	err := d.newSession.DB(dbName).C(collection).Insert(a)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.insert", collection)),
				logger.Any("content", a),
			).Error("mongodb insert failed!")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.insert", collection)),
			logger.Any("content", a),
		).Info("mongodb insert done.")
	}

	return nil
}

// FindOne 查找一条记录
func (d *DefaultSession) FindOne(collection string, result interface{}, selector bson.M, fields bson.M) error {
	err := d.newSession.DB(dbName).C(collection).Find(ExcludeDeleted(selector)).Select(fields).One(result)
	if err != nil {
		if d.printLog {
			d.printLog = false
			if err == ErrNotFound { // 没有找到不属于错误
				logger.WithFields(
					logger.String("command", fmt.Sprintf("db.%s.findOne", collection)),
					logger.Any("selector", selector),
					logger.Any("fields", fields),
				).Info("mongodb not found")
			} else {
				logger.WithFields(
					logger.Err(err),
					logger.String("command", fmt.Sprintf("db.%s.findOne", collection)),
					logger.Any("selector", selector),
					logger.Any("fields", fields),
				).Error("mongodb find one failed!")
			}
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.findOne", collection)),
			logger.Any("selector", selector),
			logger.Any("fields", fields),
			logger.Any("result", result),
		).Info("mongodb find one done.")
	}

	return nil
}

// FindAll 查找符合条件的多条记录，这里参数skip表示页码，limit表示每页多少行数据
func (d *DefaultSession) FindAll(collection string, result interface{}, selector bson.M, fields bson.M, skip int, limit int, sort ...string) error {
	err := d.newSession.DB(dbName).C(collection).Find(ExcludeDeleted(selector)).Select(fields).Sort(sort...).Skip(skip * limit).Limit(limit).All(result)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.find", collection)),
				logger.Any("selector", selector),
				logger.Any("fields", fields),
				logger.Int("skip", skip),
				logger.Int("limit", limit),
				logger.Any("sort", sort),
			).Error("mongodb find all failed!")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.find", collection)),
			logger.Any("selector", selector),
			logger.Any("fields", fields),
			logger.Int("skip", skip),
			logger.Int("limit", limit),
			logger.Any("sort", sort),
			logger.Any("result", result),
		).Info("mongodb find all done.")
	}

	return nil
}

// UpdateOne 更新一条记录
func (d *DefaultSession) UpdateOne(collection string, selector bson.M, update bson.M) error {
	err := CheckUpdateContent(update)
	if err != nil {
		return err
	}

	err = d.newSession.DB(dbName).C(collection).Update(ExcludeDeleted(selector), UpdatedTime(update))
	if err != nil {
		if d.printLog {
			d.printLog = false
			if err == ErrNotFound {
				logger.WithFields(
					logger.String("command", fmt.Sprintf("db.%s.update", collection)),
					logger.Any("selector", selector),
					logger.Any("update", update),
				).Info("mongodb nothing to update, because not found")
			} else {
				logger.WithFields(
					logger.Err(err),
					logger.String("command", fmt.Sprintf("db.%s.update", collection)),
					logger.Any("selector", selector),
					logger.Any("update", update),
				).Error("mongodb update one failed!")
			}
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.update", collection)),
			logger.Any("selector", selector),
			logger.Any("update", update),
		).Info("mongodb update one done.")
	}

	return nil
}

// UpdateAll 更新所有匹配的记录
func (d *DefaultSession) UpdateAll(collection string, selector bson.M, update bson.M) (int, error) {
	err := CheckUpdateContent(update)
	if err != nil {
		return 0, err
	}

	changeInfo, err := d.newSession.DB(dbName).C(collection).UpdateAll(ExcludeDeleted(selector), UpdatedTime(update))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.updateAll", collection)),
				logger.Any("selector", selector),
				logger.Any("update", update),
			).Error("mongodb update all failed!")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.updateAll", collection)),
			logger.Any("selector", selector),
			logger.Any("update", update),
			logger.Int("updatedNum", changeInfo.Updated),
		).Info("mongodb update all done.")
	}

	return changeInfo.Updated, nil
}

// DeleteOne 标记性删除一条记录
func (d *DefaultSession) DeleteOne(collection string, selector bson.M) error {
	err := d.newSession.DB(dbName).C(collection).Update(ExcludeDeleted(selector), DeletedTime(bson.M{}))
	if err != nil {
		if d.printLog {
			d.printLog = false
			if err == ErrNotFound {
				logger.WithFields(
					logger.String("command", fmt.Sprintf("db.%s.update", collection)),
					logger.Any("selector", selector),
				).Info("mongodb nothing to markup delete, because not found")
			} else {
				logger.WithFields(
					logger.Err(err),
					logger.String("command", fmt.Sprintf("db.%s.update", collection)),
					logger.Any("selector", selector),
				).Error("mongodb markup delete one failed!")
			}
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
			logger.Any("selector", selector),
		).Info("mongodb markup delete one done.")
	}

	return nil
}

// DeleteAll 标记性删除所有匹配的记录
func (d *DefaultSession) DeleteAll(collection string, selector bson.M) (int, error) {
	changeInfo, err := d.newSession.DB(dbName).C(collection).UpdateAll(ExcludeDeleted(selector), DeletedTime(bson.M{}))
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
				logger.Any("selector", selector),
			).Error("mongodb markup delete all failed!")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
			logger.Any("selector", selector),
			logger.Int("removedNum", changeInfo.Updated),
		).Info("mongodb markup delete all done.")
	}

	return changeInfo.Updated, nil
}

// DeleteOneReal 删除一条记录，包括标记性删除的记录
func (d *DefaultSession) DeleteOneReal(collection string, selector bson.M) error {
	err := d.newSession.DB(dbName).C(collection).Remove(selector)
	if err != nil {
		if d.printLog {
			d.printLog = false
			if err == ErrNotFound {
				logger.WithFields(
					logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
					logger.Any("selector", selector),
				).Info("mongodb nothing to real delete, because not found")
			} else {
				logger.WithFields(
					logger.Err(err),
					logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
					logger.Any("selector", selector),
				).Error("mongodb delete one failed!")
			}
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
			logger.Any("selector", selector),
		).Info("mongodb delete one done.")
	}

	return nil
}

// DeleteAllReal 删除所有匹配的记录，包括标记性删除的记录
func (d *DefaultSession) DeleteAllReal(collection string, selector bson.M) (int, error) {
	if len(selector) == 0 {
		err := errors.New("mongodb selector is nil, do not any operate")
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.count", collection)),
				logger.Any("selector", selector),
			).Error("mongodb real delete all failed!")
		}

		return 0, err
	}

	changeInfo, err := d.newSession.DB(dbName).C(collection).RemoveAll(selector)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
				logger.Any("selector", selector),
			).Error("mongodb real delete all failed!")
		}
		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.remove", collection)),
			logger.Any("selector", selector),
			logger.Int("removedNum", changeInfo.Removed),
		).Info("mongodb real delete all done.")
	}

	return changeInfo.Removed, nil
}

// Count 统计匹配的数量，不包括标记性删除的记录
func (d *DefaultSession) Count(collection string, selector bson.M) (int, error) {
	count, err := d.newSession.DB(dbName).C(collection).Find(ExcludeDeleted(selector)).Count()
	if err != nil && err != ErrNotFound {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.count", collection)),
				logger.Any("selector", selector),
			).Info("mongodb count failed!")
		}

		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.count", collection)),
			logger.Any("selector", selector),
			logger.Int("count", count),
		).Info("mongodb count done.")
	}

	return count, nil
}

// CountAll 统计匹配的数量，包括标记性删除的记录
func (d *DefaultSession) CountAll(collection string, selector bson.M) (int, error) {
	count, err := d.newSession.DB(dbName).C(collection).Find(selector).Count()
	if err != nil && err != ErrNotFound {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.count", collection)),
				logger.Any("selector", selector),
			).Info("mongodb count all failed!")
		}

		return 0, err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.count", collection)),
			logger.Any("selector", selector),
			logger.Int("count", count),
		).Info("mongodb count all done.")
	}

	return count, nil
}

// FindAndModify 更新并返回最新记录
func (d *DefaultSession) FindAndModify(collection string, result interface{}, selector bson.M, update bson.M) error {
	err := CheckUpdateContent(update)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.find.apply", collection)),
				logger.Any("selector", selector),
				logger.Any("update", update),
			).Error("mongodb find and modify failed!")
		}
		return err
	}

	change := mgo.Change{ReturnNew: true, Update: update}
	changeInfo, err := d.newSession.DB(dbName).C(collection).Find(ExcludeDeleted(selector)).Apply(change, result)
	if err != nil {
		if d.printLog {
			d.printLog = false
			if err == ErrNotFound {
				logger.WithFields(
					logger.String("command", fmt.Sprintf("db.%s.find.apply", collection)),
					logger.Any("selector", selector),
					logger.Any("update", update),
				).Info("mongodb nothing to modify, because not found")
			} else {
				logger.WithFields(
					logger.Err(err),
					logger.String("command", fmt.Sprintf("db.%s.find.apply", collection)),
					logger.Any("selector", selector),
					logger.Any("update", update),
				).Error("mongodb find and modify failed!")
			}
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.find.apply", collection)),
			logger.Any("selector", selector),
			logger.Any("update", update),
			logger.Any("result", result),
			logger.Int("updatedNum", changeInfo.Updated),
		).Info("mongodb find and modify done.")
	}

	return nil
}

// EnsureIndexKey 设置索引key
func (d *DefaultSession) EnsureIndexKey(collection string, indexKeys ...string) error {
	err := d.newSession.DB(dbName).C(collection).EnsureIndexKey(indexKeys...)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.ensureIndexKey", collection)),
				logger.Any("indexKeys", indexKeys),
			).Error("mongodb set index key failed!")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.ensureIndexKey", collection)),
			logger.Any("indexKeys", indexKeys),
		).Error("mongodb set index key done.")
	}

	return nil
}

// EnsureIndex 设置索引
func (d *DefaultSession) EnsureIndex(collection string, index mgo.Index) error {
	err := d.newSession.DB(dbName).C(collection).EnsureIndex(index)
	if err != nil {
		if d.printLog {
			d.printLog = false
			logger.WithFields(
				logger.Err(err),
				logger.String("command", fmt.Sprintf("db.%s.ensureIndex", collection)),
				logger.Any("index", index),
			).Error("mongodb set index failed!")
		}
		return err
	}

	if d.printLog {
		d.printLog = false
		logger.WithFields(
			logger.String("command", fmt.Sprintf("db.%s.ensureIndex", collection)),
			logger.Any("index", index),
		).Error("mongodb set index done.")
	}

	return nil
}
