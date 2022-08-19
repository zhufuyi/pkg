package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB gorm.DB 别名
type DB = gorm.DB

// Init 初始化mysql
func Init(dns string, opts ...Option) (*gorm.DB, error) {
	o := defaultOptions()
	o.apply(opts...)

	sqlDB, err := sql.Open("mysql", dns)
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(o.maxIdleConns)       // 空闲连接数
	sqlDB.SetMaxOpenConns(o.maxOpenConns)       // 最大连接数
	sqlDB.SetConnMaxLifetime(o.connMaxLifetime) // 断开多余的空闲连接事件

	db, err := gorm.Open(mysqlDriver.New(mysqlDriver.Config{Conn: sqlDB}), gormConfig(o))
	if err != nil {
		return nil, fmt.Errorf("gorm.Open error, err: %v", err)
	}
	db.Set("gorm:table_options", "CHARSET=utf8mb4") // 创建表时自动追加表后缀

	if o.enableTrace {
		err = db.Use(otelgorm.NewPlugin())
		if err != nil {
			return nil, fmt.Errorf("using gorm opentelemetry, err: %v", err)
		}
	}

	return db, nil
}

// gorm设置
func gormConfig(o *options) *gorm.Config {
	config := &gorm.Config{
		// 禁止外键约束, 生产环境不建议使用外键约束
		DisableForeignKeyConstraintWhenMigrating: o.disableForeignKey,
		// 去掉表名复数
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	}

	// 打印所有SQL
	if o.isLog {
		config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	// 只打印慢查询
	if o.slowThreshold > 0 {
		config.Logger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), //将标准输出作为Writer
			logger.Config{
				SlowThreshold: o.slowThreshold,
				Colorful:      true,
				LogLevel:      logger.Warn, //设置日志级别，只有指定级别以上会输出慢查询日志
			},
		)
	}

	return config
}
