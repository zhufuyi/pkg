package redis

import (
	"fmt"
	"time"

	"pkg/logger"

	"github.com/gomodule/redigo/redis"
)

// 当前应用允许最大连接数，最大不能超过redis极限连接数
// 如果多个项目共用同一个redis，要考虑每个项目限制连接数，防止过量连接造成redis卡死
var maxActiveCount = 2800
var pool *RedisPool

// ErrNil redis空值
var ErrNil = redis.ErrNil

// GetConn 对redis操作前，先要调用此函数从redis连接池获取一个连接，要对结果做错误判断，防止获取nil的连接，然后对空指针操作造成panic
func GetConn() (RedisConn, error) {
	if pool == nil {
		logger.Panic("redis pool is nil, go to connect redis first, eg: redis.NewRedisPool(server, password string)")
	}

	// 超出redis承受的极限最大连接数，直接拦截，并返回错误
	if pool.ActiveCount() > maxActiveCount {
		return nil, fmt.Errorf("redis connect clients exceeded the limit of %d", maxActiveCount)
	}

	return pool.Get(), nil
}

// RedisPool redis连接池
type RedisPool struct {
	redis.Pool
}

// Get 从连接池获取redis连接
func (r *RedisPool) Get() RedisConn {
	conn := r.Pool.Get()
	return &DefaultRedisConn{Conn: conn}
}

// RedisConn redis接口
type RedisConn interface {
	redis.Conn
	WithLog() RedisConn
}

// DefaultRedisConn 默认redis连接信息
type DefaultRedisConn struct {
	redis.Conn
	printLog bool
}

// WithLog 设置打印日志
func (d *DefaultRedisConn) WithLog() RedisConn {
	d.printLog = true
	return d
}

// Do The Do method combines the functionality of the Send, Flush and Receive methods.
func (d *DefaultRedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	result, err := d.Conn.Do(commandName, args...)
	if err != nil {
		if d.printLog {
			d.printLog = false

			logger.WithFields(
				logger.Err(err),
				logger.String("command", joinCommandAndArgs(commandName, args...)),
			).Error("redis do error")
		}
		return result, err
	}

	if d.printLog {
		d.printLog = false

		logger.WithFields(
			logger.String("command", joinCommandAndArgs(commandName, args...)),
			anyField("result", result),
		).Info("redis do success")
	}

	return result, err
}

// Send writes the command to the client's output buffer.
func (d *DefaultRedisConn) Send(commandName string, args ...interface{}) error {
	err := d.Conn.Send(commandName, args...)
	if err != nil {
		if d.printLog {
			d.printLog = false

			logger.WithFields(
				logger.Err(err),
				logger.String("command", joinCommandAndArgs(commandName, args...)),
			).Error("redis send error")
		}
		return err
	}

	if d.printLog {
		d.printLog = false

		logger.WithFields(
			logger.String("command", joinCommandAndArgs(commandName, args...)),
		).Info("redis send success")
	}

	return err
}

// Flush flushes the output buffer to the Redis server.
func (d *DefaultRedisConn) Flush() error {
	err := d.Conn.Flush()
	if err != nil {
		if d.printLog {
			d.printLog = false

			logger.WithFields(logger.Err(err)).Error("redis flush error")
		}
		return err
	}

	if d.printLog {
		d.printLog = false

		logger.WithFields().Info("redis flush success")
	}

	return nil
}

// Receive receives a single reply from the Redis server
func (d *DefaultRedisConn) Receive() (reply interface{}, err error) {
	result, err := d.Conn.Receive()
	if err != nil {
		if d.printLog {
			d.printLog = false

			logger.WithFields(logger.Err(err)).Error("redis receive error")
		}
		return result, err
	}

	if d.printLog {
		d.printLog = false

		logger.WithFields(anyField("result", result)).Info("redis receive success")
	}
	return result, err
}

// NewRedisPool connect redis，if test ping failed，return error
func NewRedisPool(server, password string) error {
	pool = &RedisPool{
		Pool: redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}

				if _, err = c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}

				c.Do("select", 0)
				return c, err
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	rconn, _ := GetConn()
	defer rconn.Close()

	_, err := rconn.Do("PING")
	return err
}

// NewRedisPoolWithNoAuth connect redis with no auth, if test ping failed，return error
func NewRedisPoolWithNoAuth(server string) error {
	pool = &RedisPool{
		Pool: redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}

				c.Do("select", 0)
				return c, err
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}

	rconn, _ := GetConn()
	defer rconn.Close()

	_, err := rconn.Do("PING")
	return err
}

// ------------------------------------------- 重新包装函数 ---------------------------------------

// Int 整型
func Int(reply interface{}, err error) (int, error) {
	return redis.Int(reply, err)
}

// Int64 64位整型
func Int64(reply interface{}, err error) (int64, error) {
	return redis.Int64(reply, err)
}

// Uint64 64位非负数整型
func Uint64(reply interface{}, err error) (uint64, error) {
	return redis.Uint64(reply, err)
}

// Float64 64位浮点数
func Float64(reply interface{}, err error) (float64, error) {
	return redis.Float64(reply, err)
}

// String 字符串类型
func String(reply interface{}, err error) (string, error) {
	return redis.String(reply, err)
}

// Bytes 字节流类型
func Bytes(reply interface{}, err error) ([]byte, error) {
	return redis.Bytes(reply, err)
}

// Bool 布尔类型
func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

// Values 任意类型
func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}

// Strings 字符串slice类型
func Strings(reply interface{}, err error) ([]string, error) {
	return redis.Strings(reply, err)
}

// ByteSlices 字节流slice类型
func ByteSlices(reply interface{}, err error) ([][]byte, error) {
	return redis.ByteSlices(reply, err)
}

// Ints 整型slice类型
func Ints(reply interface{}, err error) ([]int, error) {
	return redis.Ints(reply, err)
}

// StringMap kv都是字符串的map类型
func StringMap(result interface{}, err error) (map[string]string, error) {
	return redis.StringMap(result, err)
}

// IntMap key字符串，val为int的map类型
func IntMap(result interface{}, err error) (map[string]int, error) {
	return redis.IntMap(result, err)
}

// Int64Map key字符串，val为int64的map类型
func Int64Map(result interface{}, err error) (map[string]int64, error) {
	return redis.Int64Map(result, err)
}

// -------------------------------------------------------------------------------------------------

// 转换类型
func anyField(key string, a interface{}) logger.Field {
	switch a.(type) {
	case []byte:
		return logger.String(key, string(a.([]byte)))

	case string:
		return logger.String(key, a.(string))

	case []interface{}:
		value := []string{}
		for _, v := range a.([]interface{}) {
			if _, ok := v.([]byte); ok {
				value = append(value, fmt.Sprintf("%s", v))
			}
		}
		return logger.Any(key, value)

	default:
		return logger.String(key, fmt.Sprintf("%v", a))
	}
}

// 把连接命令和参数转换为字符串
func joinCommandAndArgs(commandName string, args ...interface{}) string {
	commandArgs := commandName
	for _, arg := range args {
		commandArgs += fmt.Sprintf(" %v", arg)
	}

	return commandArgs
}
