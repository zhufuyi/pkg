package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"

	"io/ioutil"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// DB gorm.DB 别名
type DB = gorm.DB

var (
	// ErrNotFound 空记录
	ErrNotFound = errors.New("record not found")
)

// Init 初始化mysql
func Init(addr string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	if o.caFile != "" {
		// 优先使用addr里的tls值
		if !strings.Contains(addr, "tls=") {
			addr += fmt.Sprintf("&tls=%s", o.tlsKey)
		} else {
			tlsKey := getTLSKeyFromAddr(addr)
			if tlsKey != "" {
				o.tlsKey = tlsKey
			}
		}
		registerTLS(o.tlsKey, o.caFile, o.clientKeyFile, o.clientCertFile)
	}

	db, err := gorm.Open("mysql", addr)
	if err != nil {
		return nil, err
	}

	db.DB().SetMaxIdleConns(o.maxIdleConns)       // 空闲连接数
	db.DB().SetMaxOpenConns(o.maxOpenConns)       // 最大连接数
	db.DB().SetConnMaxLifetime(o.connMaxLifetime) // 断开多余的空闲连接事件
	db.SingularTable(true)                        // 保持表名和对象名一致
	if o.IsLog {
		db.LogMode(true)                   // 开启日志
		db.SetLogger(newGormLogger(o.log)) // 自定义日志
	}

	if err = db.DB().Ping(); err != nil {
		return nil, err
	}

	if len(o.tables) > 0 {
		db.AutoMigrate(o.tables...) // 如果表不存在，自动创建，只支持自动添加新的列，对于存在的列不可以修改列属性
	}

	return db, nil
}

// registerTLS 注册tls配置，caFile是ca.pem文件路径， 规定可变参数keyCertFiles的第一个值为client-key.pem文件路径, 第二值为client-cert.pem文件路径
func registerTLS(tlsKey string, caFile string, clientKeyCertFiles ...string) {
	var clientKeyFile, clientCertFile string
	if len(clientKeyCertFiles) == 2 {
		clientKeyFile = clientKeyCertFiles[0]
		clientCertFile = clientKeyCertFiles[1]
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
		mysql.RegisterTLSConfig(tlsKey, &tls.Config{
			RootCAs:            rootCertPool,
			Certificates:       clientCert,
			InsecureSkipVerify: true,
		})
	} else { // 只使用ca.pem
		mysql.RegisterTLSConfig(tlsKey, &tls.Config{
			RootCAs:            rootCertPool,
			InsecureSkipVerify: true,
		})
	}
}

func getTLSKeyFromAddr(addr string) string {
	splits := strings.Split(addr, "&")
	for _, split := range splits {
		if strings.Contains(split, "tls=") {
			return strings.Replace(split, "tls=", "", -1)
		}
	}
	return ""
}
