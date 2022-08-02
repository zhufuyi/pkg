## mysql客户端

在[gorm](https://github.com/jinzhu/gorm)基础上封装库，例如日志，分页查询等。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/mysql

<br>

## 使用示例

### 初始化连接示例

```go
    var addr = "vison:123456@(192.168.1.2:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

    // (1) 使用默认设置连接数据库
    db, err := mysql.Init(addr)

    // (2) 自定义设置连接数据库
    db, err := mysql.Init(
        addr,
        mysql.WithLog(logger.Get()),                         // 打印日志，默认不打印
        mysql.WithMaxIdleConns(5),                          // 空闲连接数，默认3
        mysql.WithMaxOpenConns(50),                      // 最大连接数，默认30
        mysql.WithConnMaxLifetime(time.Minute*3),  // 多久断开多余的空闲连接，默认5分钟
    )

    // (3) TLS连接数据库
    db, err := Init(
        addr,
        mysql.WithLog(logger.Get()),
        mysql.WithMaxIdleConns(5),
        mysql.WithMaxOpenConns(50),
        mysql.WithConnMaxLifetime(time.Minute*3),

        //mysql.WithTLSKey("custom"),
        // 只使用ca.pem认证，配置文件的[mysqld]设置了require_secure_transport = ON，
        // 并且设置用户是要求ssl连接，
        // ALTER USER 'vison'@'%' REQUIRE SSL;
        // grant all privileges on *.* to 'vison'@'%';
        // FLUSH PRIVILEGES;
        mysql.WithCAFile("certs/ca.pem"),
    )

    // (4) TLSx509连接数据库
	db, err := Init(
		addr,
		mysql.WithLog(logger.Get()),
		mysql.WithMaxIdleConns(5),
		mysql.WithMaxOpenConns(50),
		mysql.WithConnMaxLifetime(time.Minute*3),

		//mysql.WithTLSKey("custom"),
		// 只使用ca.pem认证，配置文件的[mysqld]设置了require_secure_transport = ON，
		// 并且设置用户是要求ssl连接，
		// Create user 'vison'@'%' identified by '123456' REQUIRE SSL; 或 ALTER USER 'vison'@'%' REQUIRE SSL;
		mysql.WithCAFile("certs/ca.pem"),

		// 在开启ssl基础上，如果mysql设置用户要求x509认证时使用，
		// Create user 'vison'@'%' identified by '123456' REQUIRE X509;  或 ALTER USER 'vison'@'%' REQUIRE X509;
		// grant all privileges on *.* to 'vison'@'%';
		// FLUSH PRIVILEGES;
		// 未使用证书报错信息 Error 1045: Access denied for user 'vison'@'192.168.3.27' (using password: YES)
		mysql.WithClientKeyFile("certs/client-key.pem"),
		mysql.WithClientCertFile("certs/client-cert.pem"),
	)
```

<br>

### model示例

```go
package model

import (
	"github.com/zhufuyi/pkg/mysql"
)

// UserExample object fields mapping table
type UserExample struct {
	mysql.Model

	/* todo
	Name   string `gorm:"type:varchar(40);unique_index;not null" json:"name"`
	Age    int    `gorm:"not null" json:"age"`
	Gender string `gorm:"type:varchar(10);not null" json:"gender"`
	*/
}

// TableName get table name
func (table *UserExample) TableName() string {
	return mysql.GetTableName(table)
}

// Create a new record
func (table *UserExample) Create(db *mysql.DB) error {
	return db.Create(table).Error
}

// Delete record
func (table *UserExample) Delete(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).Delete(table).Error
}

// DeleteByID delete record by id
func (table *UserExample) DeleteByID(db *mysql.DB) error {
	return db.Where("id = ?", table.ID).Delete(table).Error
}

// Updates record
func (table *UserExample) Updates(db *mysql.DB, update mysql.KV, query interface{}, args ...interface{}) error {
	return db.Model(table).Where(query, args...).Updates(update).Error
}

// Get one record
func (table *UserExample) Get(db *mysql.DB, query interface{}, args ...interface{}) error {
	return db.Where(query, args...).First(table).Error
}

// GetByID get record by id
func (table *UserExample) GetByID(db *mysql.DB, id uint64) error {
	return db.Where("id = ?", id).First(table).Error
}

// Gets multiple records, starting from page 0
func (table *UserExample) Gets(db *mysql.DB, page *mysql.Page, query interface{}, args ...interface{}) ([]*UserExample, error) {
	out := []*UserExample{}
	err := db.Order(page.Sort()).Limit(page.Size()).Offset(page.Offset()).Where(query, args...).Find(&out).Error
	return out, err
}

// Count number of statistics
func (table *UserExample) Count(db *mysql.DB, query interface{}, args ...interface{}) (int, error) {
	count := 0
	err := db.Model(table).Where(query, args...).Count(&count).Error
	return count, err
}
```

<br>

### 事务示例

```go
func CreateUser() error {
	// 注意，当你在一个事务中应使用 tx 作为数据库句柄
	tx := db.Begin()
	if r := recover(); r != nil {   // 在事务执行过程发生panic后回滚
		fmt.Printf("transaction failed, err = %v\n", r)
		tx.Rollback()
	}

	var err error
	if err = tx.Error; err != nil {
		return err
	}

	if err = tx.Create(&User{Name: "zhangsan", Age: 5, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	//panic("发生了异常")

	if err = tx.Create(&User{Name: "lisi", Age: 4, Gender: "男"}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
```


更多使用查看gorm的使用指南 https://gorm.io/zh_CN/docs/index.html
