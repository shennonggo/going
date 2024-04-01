# Going

## 介绍

* Going是一个Golang框架，采用三层代码结构进行系统解耦。

* http服务基于gin封装。rpc服务基于gRPC封装。http服务的认证方式默认是：JWT。

* 借鉴kratos服务注册进行解耦。默认使用consul进行服务的注册、发现、负载。

* 使用cobra、pflag、viper进行配置管理。

* 自定义log包 基于zap封装、errors包，适配官方的包。

* 链路追踪使用opentelemetry，默认使用jaeger。已集成到gin、gRPC。

* 监控指标使用prometheus，已集成到gin、gRPC。

## 安装教程


```
git clone https://github.com/shennonggo/going
```

在app/user/srv/app.go中默认读取配置文件
也就是得启动consul、jaeger、mysql、redis、nacos、elasticsearch

> 若不想启动请在app/user/srv/app.go的NewApp结构体中添加app.WithNoConfig()

启动时可使用--help进行查看。

运行：

```
go run .\cmd\user\user.go -c configs/user/srv.yaml
```

gRPC启动后台服务命令：
```
go run ./cmd/admin/admin.go -c=configs/admin/admin.yaml
```

## 详细介绍

### 三层目录结构

采用三层代码结构来解耦：controller(参数校验) ->service(具体的业务逻辑) -> data(数据库的接口)

以user的grpc服务实现为例：

1. 定义api proto，生成proto文件

```bash
# 根据定义的api/user/user.proto生成proto文件
protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto

# mac安装protoc:
brew search protobuf
brew install protobuf@3

brew install grpc
brew install protoc-gen-go
brew install protoc-gen-go-grpc

go version
protoc --version
protoc-gen-go --version
```

2. 在app/user/srv/internal下新建controller、service、data目录

3. 实现controller/user.go







