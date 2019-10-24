package mysql

import (
	"testing"

	"time"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/zhufuyi/logger"
)

var addr = "root:123456@(192.168.101.88:3306)/account?charset=utf8mb4&parseTime=True&loc=Local"

type User struct {
	Model

	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
}


func init() {
	AddTables(&User{}) // 把所有表对应的对象添加过来

	//RegisterTLS("./ca.pem") // 只使用ca.pem认证，配置文件设置了require_secure_transport = ON情况下使用
	//RegisterTLS("./ca.pem", "./client-key.pem", "./client-cert.pem") // mysql设置用户强制要求x509认证时使用

	err := Init(addr, true)
	if err != nil {
		logger.Fatal("connect to mysql failed", logger.Err(err), logger.Any("addr", strings.Split(addr, "@")[1:]))
	}
	logger.Info("connect mysql success")
}

// ----------------------------------------------------- 插入 --------------------------------------------

func TestInsert(t *testing.T) {
	user := &User{Name: "小乔3", Age: 15, Gender: "女"}
	if err := GetDB().Create(user).Error; err != nil {
		logger.Error("insert error", logger.Err(err), logger.Any("user", user))
	}
}

// ---------------------------------------------------- 删除 ---------------------------------------------

func TestDelete(t *testing.T) {
	// 软删除，查询时会被忽略，如果想查询被软删除记录，在where前使用Unscoped()
	if err := GetDB().Where("name = ?", "小乔2").Delete(&User{}).Error; err != nil {
		logger.Error("delete error", logger.Err(err))
	}

	// 物理删除
	if err := GetDB().Unscoped().Where("name = ?", "小乔3").Delete(&User{}).Error; err != nil {
		logger.Error("delete error", logger.Err(err))
	}
}

// -------------------------------------------------- 修改 -----------------------------------------------

func TestUpdate(t *testing.T) {
	// Save会更新所有字段，即使你没有赋值也会替换，不建议使用
	//user := &User{}
	//GetDB().First(user, "id = ?", 21)
	//user.Name = "小乔5"
	//user.Age = 15
	//db.Save(&user)

	// 使用map更新指定多个字段(Updates)
	update := KV{"name": "小乔7", "age": 18}
	if err := GetDB().Model(&User{}).Where("id = ?", 21).Updates(update).Error; err != nil {
		t.Error(err)
	}

	// 使用struct更新指定多个字段(Updates)，只会更新其中有变化且为非零值的字段
	updateFields := User{Name: "小乔8", Age: 19}
	if err := GetDB().Model(&User{}).Where("id = ?", 21).Updates(updateFields).Error; err != nil {
		t.Error(err)
	}

	// 使用表达式更新
	update = KV{"age": gorm.Expr("age  + ?", 10)}
	if err := GetDB().Model(&User{}).Where("id = ?", 21).Updates(update).Error; err != nil {
		t.Error(err)
	}
}

// ------------------------------------------------- 查询 ------------------------------------------------

// first和find区别是：first获取一条记录，未找到时会返回错误，find获取多条记录，未找到记录返回空，不会报错
func TestQueryNormal(t *testing.T) {
	// 查找第一条记录
	user := &User{}
	if err := GetDB().First(user).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("查找第一条记录", logger.Any("user", user))

	// 通过主键查询最新一条记录
	user = &User{}
	if err := GetDB().Last(user).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("通过主键查询最新一条记录", logger.Any("user", user))

	// 根据主键id查找指定一条记录(只可在主键为整数型时使用)
	user = &User{}
	if err := GetDB().First(user, 10).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("根据主键id查找指定一条记录", logger.Any("user", user))

	// 获取所有记录
	users := []User{}
	if err := GetDB().Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("获取所有记录", logger.Any("users", users))
}

// where 查询
func TestQueryWithWhere(t *testing.T) {
	// 获取第一条匹配的记录
	user := &User{}
	if err := GetDB().Where("name = ?", "刘备").First(user).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("获取第一条匹配的记录", logger.Any("user", user))

	// 获取全部匹配的记录
	users := []User{}
	if err := GetDB().Where("age = ?", 23).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("获取全部匹配的记录", logger.Any("users", users))

	// <>反向条件查询
	users = []User{}
	if err := GetDB().Where("gender <> ?", "男").Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("<>反向件查询", logger.Any("users", users))

	// IN集合查询
	users = []User{}
	if err := GetDB().Where("name IN (?)", []string{"刘备", "关羽", "张三"}).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("IN集合查询", logger.Any("users", users))

	// LIKE 模糊查询
	users = []User{}
	if err := GetDB().Where("name LIKE ?", "%乔%").Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("LIKE模糊查询", logger.Any("users", users))

	// AND 多条件查询
	users = []User{}
	if err := GetDB().Where("age > ? AND gender = ?", 25, "男").Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("AND多条件查询", logger.Any("users", users))

	// 时间条件查询
	users = []User{}
	layaout := "2006-01-02 15:04:05"
	pt, _ := time.ParseInLocation(layaout, "2019-10-14 18:00:00", time.Local)
	if err := GetDB().Where("created_at > ?", pt).Find(&users).Error; err != nil {
		//if err = GetDB().Where("created_at > ?", time.Unix(1571047200, 0)).Find(&users).Error;err != nil {
		t.Error(err)
	}
	logger.Info("时间条件查询", logger.Any("users", users))

	// BETWEEN范围查询，包括边沿值
	users = []User{}
	if err := GetDB().Where("age BETWEEN ? AND ?", 20, 30).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("BETWEEN范围查询", logger.Any("users", users))
}

// Struct & Map & Slice查询
func TestQueryWithType(t *testing.T) {
	// struct查询，当通过结构体进行查询时，GORM将会只通过非零值字段查询，这意味着如果你的字段值为0，''， false 或者其他 零值时，将不会被用于构建查询条件
	user := &User{}
	if err := GetDB().Where(&User{Name: "刘备"}).First(user).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("struct查询", logger.Any("user", user))

	// map查询
	users := []User{}
	if err := GetDB().Where(map[string]interface{}{"name": "曹操"}).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("map查询", logger.Any("users", users))

	// Slice查询，主键id
	users = []User{}
	if err := GetDB().Where([]int{2, 4, 6}).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("slice查询", logger.Any("users", users))
}

// Not条件查询
func TestQueryWithNot(t *testing.T) {
	// not查询
	users := []User{}
	if err := GetDB().Not("gender", "男").Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("not查询", logger.Any("users", users))

	// Not In查询
	users = []User{}
	if err := GetDB().Not("name", []string{"刘备", "关羽", "张飞"}).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("not in查询", logger.Any("users", users))

	// Not In slice of primary keys
	users = []User{}
	if err := GetDB().Not([]int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("not in主键", logger.Any("users", users))
}

// Or条件查询
func TestQueryWithOr(t *testing.T) {
	// or查询
	users := []User{}
	if err := GetDB().Where("gender = ?", "女").Or("name = ?", "赵云").Find(&users).Error; err != nil {
		t.Error(err)
	}
	logger.Info("or查询", logger.Any("users", users))
}

// Inline Condition 内联条件(推荐使用)，和where条件类似
// 当与多个立即执行方法 一起使用时, 内联条件不会传递给后面的立即执行方法。
func TestQueryWithInline(t *testing.T) {
	// 根据主键获取记录 (只适用于整形主键)
	user := &User{}
	if err := GetDB().First(user, 5).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("int主键内联条件查询查询", logger.Any("user", user))

	// 根据主键获取记录, 如果它是一个非整形主键
	user = &User{}
	if err := GetDB().First(user, "id = ?", "5").Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("非int类型主键内联条件查询", logger.Any("user", user))

	// 内联条件查询一条记录
	user = &User{}
	if err := GetDB().First(user, "name = ?", "貂蝉").Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("内联条件查询一条记录", logger.Any("user", user))

	// 内联条件查询多条记录
	users := []User{}
	if err := GetDB().Find(&users, "gender = ?", "女").Error; err != nil {
		t.Error(err)
	}
	logger.Info("or查询", logger.Any("users", users))
}

// Attrs 查询和插入记录，如果记录未找到，将使用参数创建struct和记录，如果找到记录，从数据库获取数据
// 只支持struct条件查询
func TestQueryWithAttrs(t *testing.T) {
	// 不存在
	user := &User{}
	if err := GetDB().Where(User{Name: "孟获"}).Attrs(User{Age: 30, Gender: "男"}).FirstOrCreate(&user).Error; err != nil && err != ErrNotFound {
		t.Error(err)
	}
	logger.Info("不存在时使用Attrs参数创建新记录", logger.Any("user", user))
}

// Select选择字段
func TestQueryWithSelect(t *testing.T) {
	users := []User{}
	fields := []string{"name, age"}
	if err := GetDB().Select(fields).Find(&users, "gender = ?", "女").Error; err != nil {
		t.Error(err)
	}
	logger.Info("select查询", logger.Any("users", users))
}

// 排序、数量、偏移
func TestQueryWithOrderLimitOffset(t *testing.T) {
	// order 可以多个列排序
	users := []User{}
	if err := GetDB().Order("age desc").Find(&users, "gender = ?", "女").Error; err != nil {
		t.Error(err)
	}
	logger.Info("order排序", logger.Any("users", users))

	// limit 限制
	users = []User{}
	if err := GetDB().Limit(3).Find(&users, "gender = ?", "男").Error; err != nil {
		t.Error(err)
	}
	logger.Info("limit限制", logger.Any("users", users))

	// offset 偏移
	users = []User{}
	if err := GetDB().Limit(3).Offset(10).Find(&users, "gender = ?", "男").Error; err != nil {
		t.Error(err)
	}
	logger.Info("offset偏移", logger.Any("users", users))
}

func TestQueryWithPage(t *testing.T) {
	// 获取全部字段，不需要另外传入表名
	users := []User{}
	where := "gender = '男'"
	if err := FindPage(&users, where, 0, 5); err != nil {
		t.Error(err)
	}
	logger.Info("按页获取数据", logger.Any("users", users))

	// 获取部分字段，需要指定表名
	type Result struct {
		Name   string
		Age    int
		Gender string
	}

	results := []Result{}
	where = "gender = '男'"
	sort := "age desc"
	if err := FindPage2(&results, &User{}, where, 0, 5, sort); err != nil {
		t.Error(err)
	}
	logger.Info("按页获取数据2", logger.Any("results", results))
}

func FindPage(out interface{}, where string, page int, limit int, sort ...string) error {
	order := "id desc"
	if len(sort) > 0 {
		order = sort[0]
	}
	if page < 0 {
		page = 0
	}

	return GetDB().Limit(limit).Order(order).Offset(page*limit).Find(out, where).Error
}

func FindPage2(out interface{}, tableModel interface{}, where string, page int, limit int, sort ...string) error {
	order := "id desc"
	if len(sort) > 0 {
		order = sort[0]
	}
	if page < 0 {
		page = 0
	}

	return GetDB().Model(tableModel).Limit(limit).Order(order).Offset(page * limit).Where(where).Scan(out).Error
}

// count 统计数量
func TestQueryWithCount(t *testing.T) {
	count := 0

	// 只获取数量，不获取结果
	if err := GetDB().Model(&User{}).Where("gender = ?", "男").Count(&count).Error; err != nil {
		//if err := GetDB().Table("user").Where("gender = ?", "男").Count(&count).Error;err != nil {
		t.Error(err)
	}
	logger.Info("count 统计数量，不获取结果", logger.Int("count", count))

	// 获取数量同时获取结果数据
	users := []User{}
	if err := GetDB().Where("gender = ?", "女").Find(&users).Count(&count).Error; err != nil {
		t.Error(err)
	}
	logger.Info("count 统计数量，并获取结果", logger.Int("count", count), logger.Any("users", users))
}

// ------------------------------------------------- 事务 ------------------------------------------------

func TestTransaction(t *testing.T) {
	if err := CreatePeople(); err != nil {
		t.Error(err)
	}
}

func CreatePeople() error {
	// 注意，当你在一个事务中应使用 tx 作为数据库句柄
	tx := GetDB().Begin()
	defer TxRecover(tx) // 发生panic时回滚

	var err error

	if err = tx.Error; err != nil {
		return err
	}

	if err = tx.Create(&User{Name: "阿斗", Age: 5, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}
	//panic("发生了异常")
	if err = tx.Create(&User{Name: "曹冲", Age: 4, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
