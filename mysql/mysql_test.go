package mysql

import (
	"fmt"
	"testing"
	"time"
)

var addr = "vison:123456@(192.168.3.37:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

func TestInit(t *testing.T) {
	db, err := Init(addr)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, addr=%s", err, addr))
		return
	}

	t.Logf("%+v", db.DB().Stats())
}

func TestInitNoTLS(t *testing.T) {
	db, err := Init(
		addr,
		WithLog(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),
		//WithTable(&User{}), // 自动创建表
	)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, addr=%s", err, addr))
		return
	}

	t.Logf("%+v", db.DB().Stats())
}

func TestInitTLS(t *testing.T) {
	db, err := Init(
		addr,
		WithLog(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),

		//WithTLSKey("custom"),
		// 只使用ca.pem认证，配置文件的[mysqld]设置了require_secure_transport = ON，
		// 并且设置用户是要求ssl连接，
		// ALTER USER 'vison'@'%' REQUIRE SSL;
		// grant all privileges on *.* to 'vison'@'%';
		// FLUSH PRIVILEGES;
		WithCAFile("certs/ca.pem"),
	)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, addr=%s", err, addr))
		return
	}

	t.Logf("%+v", db.DB().Stats())
}

func TestInitTLSx509(t *testing.T) {
	db, err := Init(
		addr,
		WithLog(),
		WithMaxIdleConns(5),
		WithMaxOpenConns(50),
		WithConnMaxLifetime(time.Minute*3),

		//WithTLSKey("custom"),
		// 只使用ca.pem认证，配置文件的[mysqld]设置了require_secure_transport = ON，
		// 并且设置用户是要求ssl连接，
		// Create user 'vison'@'%' identified by '123456' REQUIRE SSL; 或 ALTER USER 'vison'@'%' REQUIRE SSL;
		WithCAFile("certs/ca.pem"),

		// 在开启ssl基础上，如果mysql设置用户要求x509认证时使用，
		// Create user 'vison'@'%' identified by '123456' REQUIRE X509;  或 ALTER USER 'vison'@'%' REQUIRE X509;
		// grant all privileges on *.* to 'vison'@'%';
		// FLUSH PRIVILEGES;
		// 未使用证书报错信息 Error 1045: Access denied for user 'vison'@'192.168.3.27' (using password: YES)
		WithClientKeyFile("certs/client-key.pem"),
		WithClientCertFile("certs/client-cert.pem"),
	)
	if err != nil {
		t.Error(fmt.Sprintf("connect to mysql failed, err=%v, addr=%s", err, addr))
		return
	}

	t.Logf("%+v", db.DB().Stats())
}
