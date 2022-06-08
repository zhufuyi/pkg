## render

gin请求参数校验。

<br>

## 安装

> go get -u github.com/zhufuyi/pkg/gin/validator

<br>

## 使用

```go
type userForm struct {
	Name  string `form:"name" binding:"required"`
	Age   int    `form:"age" binding:"gte=0,lte=150"`
	Email string `form:"email" binding:"email,required"`
}

func CreateUser(c *gin.Context) {
    form := &userForm{}
    err := c.ShouldBindJSON(form)
	// ...
}
func main() {
	r := gin.Default()
	binding.Validator = validator.Gin()
	
	r.POST("/createuser",CreateUser)
    // ...
}
```

