package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db     *gorm.DB
	tables []interface{} // 各个对象指针地址集合
	tlsKey = ""          // 默认不使用tls传输，如果不为空，使用tls传输数据

	// ErrNotFound 空记录
	ErrNotFound = errors.New("record not found")
)

// Init 初始化mysql
func Init(addr string, isEnableLog ...bool) error {
	var err error

	if tlsKey != "" {
		addr += fmt.Sprintf("&tls=%s", tlsKey)
	}

	db, err = gorm.Open("mysql", addr)
	if err != nil {
		return err
	}

	db.DB().SetMaxIdleConns(3)                  // 空闲连接数
	db.DB().SetMaxOpenConns(100)                // 最大连接数
	db.DB().SetConnMaxLifetime(3 * time.Minute) // 3分钟后断开多余的空闲连接
	db.SingularTable(true)                      // 保持表名和对象名一致

	if len(isEnableLog) == 1 && isEnableLog[0] {
		db.LogMode(true)              // 开启日志
		db.SetLogger(newGormLogger()) // 自定义日志
	}

	if err = db.DB().Ping(); err != nil {
		return err
	}

	SyncTable()

	return nil
}

// GetDB 获取连接
func GetDB() *gorm.DB {
	if db == nil {
		panic("db is nil, please reconnect mysql.")
	}
	return db
}

// AddTables 添加表
func AddTables(object ...interface{}) {
	tables = append(tables, object...)
}

// SyncTable 同步表
func SyncTable() {
	GetDB().AutoMigrate(tables...) // 确保对象和mysql表一致，只支持自动添加新的列，对于存在的列不可以修改列属性
}

// Model 表内嵌字段
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

// KV map类型
type KV map[string]interface{}

// TxRecover 在事务执行过程发生panic后回滚，使用时在前面添加defer关键字，例如：defer TxRecover(tx)
func TxRecover(tx *gorm.DB) {
	if r := recover(); r != nil {
		fmt.Printf("transaction failed, err = %v\n", r)
		tx.Rollback()
	}
}

// RegisterTLS 注册tls配置，caFile是cat.pem文件路径， 规定可变参数keyCertFiles的第一个值为client-key.pem文件路径, 第二值为client-cert.pem文件路径
func RegisterTLS(caFile string, keyCertFiles ...string) {
	tlsKey = ""

	var clientKeyFile, clientCertFile string
	if len(keyCertFiles) == 2 {
		clientKeyFile = keyCertFiles[0]
		clientCertFile = keyCertFiles[1]
	}

	if caFile == "" {
		panic("caFile is empty.")
	}

	pem, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(fmt.Sprintf("ioutil.ReadFile %s error, err = %v", caFile, err))
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic(fmt.Sprintf(" rootCertPool.AppendCertsFromPEM error, failed to append root CA cert at %s", caFile))
	}

	// 如果mysql里使用命令(ALTER USER 'vison'@'%' REQUIRE x509;)强制要求使用x509，必须使用client-key.pem client-cert.pem
	if clientKeyFile != "" && clientCertFile != "" {
		certs, err := tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			panic(fmt.Sprintf("tls.LoadX509KeyPair error, err = %v", err))
		}

		clientCert := make([]tls.Certificate, 0, 1)
		clientCert = append(clientCert, certs)
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
		tlsKey = "custom"
	} else { // 只使用ca.pem
		mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		})
		tlsKey = "custom"
	}
}
