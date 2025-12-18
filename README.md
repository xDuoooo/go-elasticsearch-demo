# Hotel Demo - Go 版本

这是一个使用 Go 语言开发的酒店管理系统，从原 Spring Boot 项目转换而来。

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- **ORM**: [GORM](https://gorm.io/) - Go 语言 ORM 库
- **配置管理**: [Viper](https://github.com/spf13/viper) - 配置解决方案
- **数据库**: MySQL

## 项目结构

```
go-demo/
├── cmd/
│   └── main.go              # 应用入口
├── internal/
│   ├── model/               # 数据模型
│   │   ├── hotel.go         # 酒店实体
│   │   └── hotel_doc.go     # 酒店文档对象
│   ├── repository/          # 数据访问层
│   │   └── hotel_repository.go
│   └── service/             # 服务层
│       └── hotel_service.go
├── config/
│   └── config.yaml          # 配置文件
├── go.mod                   # Go 模块文件
└── README.md                # 本文件
```

## 数据库配置

确保 MySQL 数据库已安装并运行，数据库名为 `heima`，包含 `tb_hotel` 表。

修改 `config/config.yaml` 中的数据库配置：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: "123"
  dbname: heima
```

## 安装依赖

```bash
cd go-demo
go mod tidy
```

## 运行应用

```bash
go run cmd/main.go
```

服务器将在 `http://localhost:8089` 启动。

## API 接口

### 健康检查
- `GET /health` - 健康检查

### 酒店管理
- `GET /api/hotels` - 获取酒店列表（支持分页）
  - 参数: `page`（页码，默认1）, `pageSize`（每页数量，默认10）
- `GET /api/hotels/:id` - 根据ID获取酒店
- `POST /api/hotels` - 创建酒店
- `PUT /api/hotels/:id` - 更新酒店
- `DELETE /api/hotels/:id` - 删除酒店
- `GET /api/hotels/city/:city` - 根据城市查询酒店
- `GET /api/hotels/brand/:brand` - 根据品牌查询酒店

## 示例请求

### 获取酒店列表
```bash
curl http://localhost:8089/api/hotels?page=1&pageSize=10
```

### 获取单个酒店
```bash
curl http://localhost:8089/api/hotels/1
```

### 创建酒店
```bash
curl -X POST http://localhost:8089/api/hotels \
  -H "Content-Type: application/json" \
  -d '{
    "name": "如家酒店",
    "address": "北京市朝阳区",
    "price": 300,
    "score": 45,
    "brand": "如家",
    "city": "北京",
    "starName": "三钻",
    "business": "CBD",
    "longitude": "116.408",
    "latitude": "39.904",
    "pic": "https://example.com/pic.jpg"
  }'
```

### 根据城市查询
```bash
curl http://localhost:8089/api/hotels/city/上海
```

## 开发说明

### 添加新功能
1. 在 `internal/model` 添加或修改数据模型
2. 在 `internal/repository` 添加数据访问方法
3. 在 `internal/service` 添加业务逻辑
4. 在 `cmd/main.go` 添加路由和处理函数

### 数据库迁移
GORM 支持自动迁移，在 `main.go` 的 `initDatabase()` 函数中取消注释：
```go
db.AutoMigrate(&model.Hotel{})
```

## 与原 Java 项目的对应关系

| Java | Go |
|------|-----|
| `Hotel.java` | `internal/model/hotel.go` |
| `HotelDoc.java` | `internal/model/hotel_doc.go` |
| `HotelMapper.java` | `internal/repository/hotel_repository.go` |
| `HotelService.java` | `internal/service/hotel_service.go` |
| `HotelDemoApplication.java` | `cmd/main.go` |
| `application.yaml` | `config/config.yaml` |

## 编译

```bash
go build -o hotel-demo cmd/main.go
```

运行编译后的文件：
```bash
./hotel-demo
```

## 许可证

MIT
