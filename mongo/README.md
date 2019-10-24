## mongo

在[mgo](https://github.com/globalsign/mgo)基础上封装mongodbd的golang客户端。

<br>

### 安装

> go get -u github.com/zhufuyi/pkg/mongo

<br>

### 使用

使用方式请看[示例](./demo)。

如果想新建对象，只需对模板批量修改名称，操作命令：

```bash
git clone https://github.com/zhufuyi/mongo.git
cd demo

# 修改名称(文件名和内容同时修改)
./rename.sh ./ yourObject demo

# 例如：./rename.sh ./ user demo
# 注：(1)不要使用其它名称替代demo；(2)package名默认为module，如果不是放在module目录下，需要手动更改。
```

<br>
