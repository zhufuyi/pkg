package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/zhufuyi/pkg/utils"
)

func init() {
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		err := NewRedisPool("192.168.101.88:6379", "123456")
		if err != nil {
			panic(err)
		}
	})
}

// 测试数据
var cmds = []struct {
	Cmd    string
	Args   []interface{}
	Result []string
}{
	// 测试key
	{"SET", []interface{}{"testKey", "hello redis"}, []string{"OK"}},
	{"GET", []interface{}{"testKey"}, []string{"hello redis"}},

	// 测试集合
	{"SADD", []interface{}{"testSet", "zhangsan", "lisi", 100}, []string{"0"}},
	{"SMEMBERS", []interface{}{"testSet"}, []string{"100", "zhangsan", "lisi"}},

	// 测试有序集合
	{"ZADD", []interface{}{"testZset", 28, "zhangsan", 24, "lisi", 26, "wangwu"}, []string{"0"}},
	{"ZRANGE", []interface{}{"testZset", 0, -1, "withscores"}, []string{"lisi", "24", "wangwu", "26", "zhangsan", "28"}},
	{"ZREVRANGEBYSCORE", []interface{}{"testZset", "+inf", "-inf", "withscores", "limit", 0, 100}, []string{"zhangsan", "28", "wangwu", "26", "lisi", "24"}},

	// 测试哈希
	{"HMSET", []interface{}{"testHSet", "name", "lisi", "age", 11}, []string{"OK"}},
	{"HGETALL", []interface{}{"testHSet"}, []string{"name", "lisi", "age", "11"}},
}

// 测试单个连接的Do函数
func TestDo(t *testing.T) {
	defer func() { recover() }()

	var rconn RedisConn
	var err error
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		rconn, err = GetConn()
		if err != nil {
			println(err)
			return
		}
	})
	defer rconn.Close()

	for _, v := range cmds {
		result, err := rconn.WithLog().Do(v.Cmd, v.Args...)
		if err != nil {
			t.Logf("redis do failed, %v, %+v\n", err, v)
			continue
		}

		actual, _ := json.Marshal(typeAsserts(result))
		expect, _ := json.Marshal(v.Result)
		if string(actual) != string(expect) {
			t.Errorf("got %v, expected %v\n", string(actual), string(expect))
		}
	}
}

// 测试单个连接的Send_Flush_Receive函数，等效于do，当超过redis最大连接数时，do命令不是并发安全的
func TestSend_Flush_Receive(t *testing.T) {
	defer func() { recover() }()

	var rconn RedisConn
	var err error
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		rconn, err = GetConn()
		if err != nil {
			println(err)
			return
		}
	})
	defer rconn.Close()

	for _, v := range cmds {
		err := rconn.WithLog().Send(v.Cmd, v.Args...)
		if err != nil {
			t.Logf("redis send error, %v, %+v\n", err, v)
			continue
		}

		err = rconn.WithLog().Flush()
		if err != nil {
			t.Logf("redis flush error, %v, %+v\n", err, v)
			continue
		}

		result, err := rconn.WithLog().Receive()
		if err != nil {
			t.Logf("redis receive error, %v, %v\n", err, v)
			continue
		}

		actual, _ := json.Marshal(typeAsserts(result))
		expect, _ := json.Marshal(v.Result)
		if string(actual) != string(expect) {
			t.Logf("got %v, expected %v", string(actual), string(expect))
		}
	}
}

func TestTransaction(t *testing.T) {
	defer func() { recover() }()

	var rconn RedisConn
	var err error
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		rconn, err = GetConn()
		if err != nil {
			println(err)
			return
		}
	})
	defer rconn.Close()

	// 使用send和do方法实现事务
	rconn.Send("MULTI")
	rconn.Send("INCR", "test-foo")
	rconn.Send("INCR", "test-bar")
	rconn.Send("HSET", "test-hset", "name", "李四", "age", 11, "gender", "male")
	resp, err := rconn.Do("EXEC")
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
}

// 并发写压测，调整maxActiveCount值大小，统计出redis能够承受并发写的最大值(也就是maxActiveCount的大小)
func TestRedisWriteLimit(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup

	for i := 0; i < maxActiveCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			err := redisWrite(i)
			if err != nil {
				t.Log(err)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nwrite success count = %d\n", successCount)
}

// 并发读压测，调整maxActiveCount值大小，统计出redis能够承受并发读的最大值(也就是maxActiveCount的大小)
func TestRedisReadLimit(t *testing.T) {
	var successCount int32
	var wg sync.WaitGroup

	for i := 0; i < maxActiveCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			err := redisRead(i)
			if err != nil {
				t.Error(err)
				return
			}

			atomic.AddInt32(&successCount, 1)
		}(i)
	}

	wg.Wait()

	fmt.Printf("\nread success count = %d\n", successCount)
}

// 读写压测，调整maxActiveCount值大小，统计出redis能够承受并发读写的最大值(也就是maxActiveCount的大小)
func TestRedisReadWriteLimit(t *testing.T) {
	var writeCount int32
	var readCount int32
	var wg sync.WaitGroup

	for i := 0; i < maxActiveCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			if i%2 == 0 {
				err := redisWrite(i)
				if err != nil {
					t.Error(err)
					return
				}
				atomic.AddInt32(&writeCount, 1)
			} else {
				err := redisRead(i)
				if err != nil {
					t.Error(err)
					return
				}
				atomic.AddInt32(&readCount, 1)
			}

		}(i)
	}

	wg.Wait()

	fmt.Printf("\nsuccess read count = %d, write  count = %d, total = %d\n", readCount, writeCount, readCount+writeCount)
}

func redisWrite(i int) error {
	defer func() { recover() }()

	var rconn RedisConn
	var err error
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		rconn, err = GetConn()
		if err != nil {
			fmt.Printf("#%d  get redis conn error, %v\n", i, err)
		}
	})
	defer rconn.Close()

	key := fmt.Sprintf("NO_%d", i)
	value := fmt.Sprintf("value_%d", i)
	result, err := String(rconn.Do("SET", key, value))
	if err != nil {
		return fmt.Errorf("#%d redis do error,%v\n", i, err)
	}

	if result != "OK" {
		return fmt.Errorf("#%d  got %v, expected %v", i, result, "OK")
	}

	return nil
}

func redisRead(i int) error {
	defer func() { recover() }()

	var rconn RedisConn
	var err error
	utils.SafeRunWithTimeout(time.Second*2, func(cancel context.CancelFunc) {
		rconn, err = GetConn()
		if err != nil {
			fmt.Printf("#%d  get redis conn error, %v\n", i, err)
		}
	})
	defer rconn.Close()

	key := fmt.Sprintf("NO_%d", i)
	value := fmt.Sprintf("value_%d", i)
	result, err := String(rconn.Do("GET", key))
	if err != nil {
		return fmt.Errorf("#%d redis do error,%v\n", i, err)
	}

	if result != value {
		return fmt.Errorf("#%d got %v, expected %v", i, result, value)
	}

	return nil
}

// 单个连接对string类型压测
func BenchmarkSingleConnString(b *testing.B) {
	rconn, err := GetConn()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("NO_%d", i)
		value := generateRandString(8)
		_, err := rconn.Do("SET", key, value)
		if err != nil {
			b.Errorf("redis execute cmd failed, %v", err)
			continue
		}
	}
}

// 单个连接对哈希类型写压测
func BenchmarkSingleConnHash(b *testing.B) {
	cmd := "HSET"
	hsetKey := "benchHsetKey"

	rconn, err := GetConn()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		memberKey := fmt.Sprintf("NO_%d", i)
		value := generateRandString(8)
		_, err := rconn.Do(cmd, hsetKey, memberKey, value)
		if err != nil {
			b.Errorf("redis execute cmd failed, %v", err)
			continue
		}
	}
}

// 类型断言
func typeAsserts(a interface{}) []string {
	var rerults []string

	switch a.(type) {
	case []byte:
		return append(rerults, string(a.([]byte)))

	case string:
		return append(rerults, a.(string))

	case []interface{}:
		for _, v := range a.([]interface{}) {
			if _, ok := v.([]byte); ok {
				rerults = append(rerults, fmt.Sprintf("%s", v))
			}
		}
		return rerults

	default:
		return append(rerults, fmt.Sprintf("%v", a))
	}
}

var stringRef = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 生成指定长度的随机字符串
func generateRandString(length int) string {
	maxIndex := len(stringRef)
	if length <= 0 {
		length = 6 // 默认长度
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	vl := make([]byte, length)
	for i := 0; i < length; i++ {
		vl[i] = stringRef[r.Intn(maxIndex)]
	}

	return string(vl)
}
