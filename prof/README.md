## prof

封装官方`net/http/pprof`路由，添加profile io wait time路由。

<br>

### 使用示例

```go
	mux := http.NewServeMux()
    prof.Register(r, WithPrefix("/myServer"), WithIOWaitTime())

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	
    if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        panic("listen and serve error: " + err.Error())
    }
```
