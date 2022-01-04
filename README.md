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
数据库支持，支持metrics/tracing 
```go
// 选择某个数据库，可以支持多实例
conn := sqlx.Get(ctx, "db1")
ctx := context.TODO()

var u User
err = conn.GetContext(ctx, &u, "select * from users where id = ?", id)
if err != nil {
return
}
```
参考[README](./pkg/sqlx/README.md)

### `api`
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

### `rpc`
接口实现层
```go
type Server struct {}

func (s *Server) CreateArticle(ctx context.Context, req *pb.CreateArticleReq) (resp *pb.CreateArticleResp, err error) {
	// TODO 调用service层代码
}
```

### `service`
逻辑处理层
```go
// 1. 调用dao层代码，例如查询blog数据
func Foo(ctx context.Context, req pb.Req) (result interface{}, err error) {
    b, err := blog.QueryByID(ctx, req.ID)
    if err != nil {
    	return
    }   
    
    total, err := blog.CountOnline(ctx)
    if err != nil {
        return
    }
    
    // 聚合结果返回
    return
}

```

### `pkg`
工具包

* conf       配置
* db         数据库
* log        日志
* metrics    prometheus
* middleware 中间件
* trace      opentracing

## 开发流程
1. 定义接口服务
   在`api/`下定义接口，建议格式`api/demo/v${num}/${xx}.proto`，参考[pb描述文件](./rpc/demo/v0/demo.proto)
   然后执行命令
   ```shell
   $ make rpc
   ```
   
2. 实现接口
    在`rpc/`下实现定义的接口，建议格式`rpc/demov0/server.go`，参考[server](./server/demov0/server.go)
    
    server一般是参数校验，和归类响应数据等
    
    service层是逻辑处理层
   
    dao层是数据库访问层
    
    调用顺序 `rpc---->svc----->dao`
   

3. 注册接口
    将实现的接口注册路由，参考[注册路由](./app/demo/cmd/server/register.go)
   
