## redis 客户端

在[redigo](https://github.com/gomodule/redigo)基础上封装redis客户端。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/redis

<br>

## 使用示例

```go
    err := redis.NewRedisPool("192.168.101.88:6379", "123456")

    // 数据
    var cmds = []struct {
        Cmd    string
        Args   []interface{}
        Result []string
    }{
        // string
        {"SET", []interface{}{"testKey", "hello redis"}, []string{"OK"}},
        {"GET", []interface{}{"testKey"}, []string{"hello redis"}},

        // 集合
        {"SADD", []interface{}{"testSet", "zhangsan", "lisi", 100}, []string{"0"}},
        {"SMEMBERS", []interface{}{"testSet"}, []string{"100", "zhangsan", "lisi"}},

        // 有序集合
        {"ZADD", []interface{}{"testZset", 28, "zhangsan", 24, "lisi", 26, "wangwu"}, []string{"0"}},
        {"ZRANGE", []interface{}{"testZset", 0, -1, "withscores"}, []string{"lisi", "24", "wangwu", "26", "zhangsan", "28"}},
        {"ZREVRANGEBYSCORE", []interface{}{"testZset", "+inf", "-inf", "withscores", "limit", 0, 100}, []string{"zhangsan", "28", "wangwu", "26", "lisi", "24"}},

        // 哈希
        {"HMSET", []interface{}{"testHSet", "name", "lisi", "age", 11}, []string{"OK"}},
        {"HGETALL", []interface{}{"testHSet"}, []string{"name", "lisi", "age", "11"}},
    }


	rconn, err := redis.GetConn()
	if err != nil {
		return
	}
	defer rconn.Close()

	for _, v := range cmds {
	    result, err := rconn.Do(v.Cmd, v.Args...)
		// result, err := rconn.WithLog().Do(v.Cmd, v.Args...) // 打印执行语句
		if err != nil {
			continue
		}
        fmt.Println(result)
	}
```
