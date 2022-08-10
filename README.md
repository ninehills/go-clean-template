# Go 网站应用模版

[![License](https://img.shields.io/github/license/ninehills/go-webapp-template.svg)](https://github.com/ninehills/go-webapp-template/blob/master/LICENSE)
[![Release](https://img.shields.io/github/v/release/evrone/go-clean-template.svg)](https://github.com/ninehills/go-webapp-template/releases/)
[![codecov](https://codecov.io/gh/evrone/go-clean-template/branch/master/graph/badge.svg?token=XE3E0X3EVQ)](https://codecov.io/gh/evrone/go-clean-template)

本工程参考了 <https://github.dev/evrone/go-clean-template> 的模版。其主要依赖的库如下：

- HTTP Framework: [Gin](https://github.com/gin-gonic/gin)
- SQL: [sqlc](https://sqlc.dev/)
- Database: [MySQL](https://github.com/go-sql-driver/mysql)

提供两套接口：

- HTTP API
- GRPC API

## 依赖工具安装

首先安装 go >= 1.17 版本，然后安装以下依赖工具：

```bash
# 如下工具也可以使用 homebrew 或者直接下载对应的 binary 安装。
# swagger 命令行工具：
go install github.com/swaggo/swag/cmd/swag@latest
# sqlc 工具
brew install sqlc
# mockery 工具
brew install mockery
# goreleaser 工具
brew install goreleaser
# golangci-lint
brew install golangci-lint
```

## 快速开始

本地开发

```sh
# Postgres, RabbitMQ
$ make compose-up
# Run app with migrations
$ make run
```

集成测试（可以在 CI 中运行）

```sh
# DB, app + migrations, integration tests
$ make compose-up-integration-test
```

## 代码结构

### `main.go`

程序入口，主要的功能在 `internal/app/app.go` 中。

本地开发环境启动命令：`make run`

### `sql`、`internal/dao`、`sqlc.yaml`

- `sql` 是 sqlc 依赖的原始 SQL 语句。
  - `schema`： 存放所有的建表语句
  - `query`: 存放所有的查询语句，最好和 schema 相对应
- `sqlc.yaml` 是 sqlc 的配置文件。
- `internal/dao` 是 sqlc 生成的代码，请不要修改。

生成方法：`make sqlc`

### `config`

配置，首先读取 `config/config.yml`中的默认内容，然后读取环境变量里面有符合的变量，将其覆盖 yml 中的配置

配置的结构在 `config.go`中
`env-required:true` 标签强制你指定值（在 yml 文件或者环境变量中）

配置使用的[cleanenv](https://github.com/ilyakaznacheev/cleanenv) 库

`make run`会从 `.env.example` 中读取测试环境的变量。

### `docs`

Swagger 文档。由 [swag](https://github.com/swaggo/swag) 库自动生成
你不需要自己修改任何内容。

生成命令：`make swag`
测试环境访问：`http://127.0.0.1:8080/swagger/index.html`

### `integration-test`

功能测试目录，它会在应用容器旁启动独立的容器。具体的测试逻辑在 integration_test.go 文件中，主要对 Restful 接口进行测试。

使用了[go-hit](https://github.com/Eun/go-hit) 库。

- `main_test.go` 为测试入口。
- `xxxx_test.go` 等为各个功能的 Restful 测试用例。

启动功能测试命令：`make integration-test` （启动之前请确保服务启动在本地并且相关依赖 Ready）

同时该测试可以测试非本地环境，配置在`main_test.go`中，使用方法为传入环境变量`ENV`，目前有如下环境

```bash
# 远程测试环境
ENV=testing make integration-test

# 本地环境，不加ENV
```

## `internal/app`

APP 主逻辑入口，其通过依赖注入的方式生成主要的业务逻辑对象，配置路由。
然后启动启动服务器并阻塞等待。

### `internal/controller`

MVC 中的控制层，服务的路由用同样的风格进行编写

- handler 按照应用领域进行分组（有共同的基础）
- 对于每一个分组，创建自己的路由结构体和请求 path 路径
- 业务逻辑的结构被注入到路由器结构中，它将被处理程序调用

### `internal/entity`

业务逻辑实体（模型），将所有实体统一定义在此处。注意实体的定义和`internal/dao`中的定义是不同的，前者更多使用在接口中，后者是数据库的原始结构。

所以此处增加转换逻辑（其实可以给 sqlc 配置生成的 model 带 json tags，但是为了隔离两层，所以宁愿人工转换）

### `internal/service`

业务逻辑的核心部分，以 Workspace 为例：

- `interfaces.go`: 将所有业务接口放到一起
- `user.go`: 实现的业务逻辑（相当于 Service）
- `user_test.go`: 对应的单元测试

此处可以自动生成单测所依赖的 mock，具体使用方法：

- 生成 mock 代码： `make mock`
- 进行单元测试： `make test`

### `pkg`

和业务逻辑无关的库。

## 开发相关

```bash
# 本地测试数据库的搭建

docker run -d --name go_webapp -p 3306:3306 -e MYSQL_ROOT_PASSWORD=go_webapp mysql:5.7
mysql -h127.0.0.1 -P3306 -uroot -pgo_webapp

CREATE DATABASE go_webapp CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
CREATE USER 'go_webapp'@'%' IDENTIFIED BY 'go_webapp';
GRANT ALL PRIVILEGES ON go_webapp.* TO 'go_webapp'@'%';
use go_webapp;

# 执行sql/schema/ 下的建表语句

```
