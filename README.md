# Nautilus

一个轻量的`gin`业务框架

安装使用
```shell
go get -u github.com/x-lambda/nautilus
```

依赖
* [protoc](https://github.com/protocolbuffers/protobuf)
* [protoc-gen-go](https://github.com/golang/protobuf/tree/master/protoc-gen-go)
* [protoc-gen-gin](https://github.com/x-lambda/protoc-gen-gin)

1. 安装`protoc`
```shell
$ brew install protobuf
```

2. 安装`protoc-gen-go`
```shell
$ git clone https://github.com/golang/protobuf.git
$ cd protobuf/protoc-gen-go
$ go install 
```

3. 安装`protoc-gen-gin`
```shell
$ go get -u github.com/x-lambda/protoc-gen-gin
```

## 项目结构
### `app`
支持多应用，可以在`app/`下定义多服务
```shell
app/
  |-demo/
        |-main.go  // demo应用
  |-example/
        |-main.go  // example应用
  |-shop/
        |-main.go  // shop应用
  |-login/
        |-main.go  // login应用
```

### `dao`
数据库支持，支持metrics/tracing，未使用`ORM`
```go
func QueryByID(ctx context.Context, id int32) (item Item, err error) {
	conn := db.Get(ctx, "default")  // 对应配置中的DSN，可以支持多DB
	sql := "select id, name, create_time, modify_time from t_demo where id = ?"
	q := db.SQLSelect("t_demo", sql)
	err = conn.QueryRowContext(ctx, q, id).Scan(&item.ID, &item.Name, &item.CreateTime, &item.ModifyTime)
	return
}
```

### `rpc`
服务定义
```proto
service BlogService {
	rpc CreateArticle(CreateArticleReq) returns (CreateArticleResp) {
		option (google.api.http) = {
			post: "/v1/author/{author_id}/articles"
		};
	}
}

message CreateArticleResp {
    int32 code = 1;
    string msg = 2;
}

message CreateArticleReq {
	string title  = 1;
	string content = 2;
	// @inject_tag: form:"author_id" uri:"author_id"
	int32 author_id = 3;
}

```

### `server`
接口实现层
```go
type Server struct {}

func (s *Server) CreateArticle(ctx context.Context, req *pb.CreateArticleReq) (resp *pb.CreateArticleResp, err error) {
	// TODO 调用service层代码
}
```

### `service`
逻辑处理层

### `util`
工具包

* conf       配置
* db         数据库
* log        日志
* metrics    prometheus
* middleware 中间件
* trace      opentracing

## 开发流程
1. 定义接口服务
   在`rpc/`下定义接口，建议格式`rpc/demo/v${num}/${xx}.proto`，参考[pb描述文件](./rpc/demo/v0/demo.proto)
   然后执行命令
   ```shell
   $ make rpc
   ```
   
2. 实现接口
    在`server/`下实现定义的接口，建议格式`server/demov0/server.go`，参考[server](./server/demov0/server.go)
    
    server一般是参数校验，和归类响应数据等
    
    service层是逻辑处理层
   
    dao层是数据库访问层
    
    调用顺序 `server---->service----->dao`
   

3. 注册接口
    将实现的接口注册路由，参考[注册路由](./app/demo/cmd/server/register.go)
   
