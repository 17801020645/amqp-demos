# AMQP Go Demo 项目使用指南

## 📋 项目结构

```
amqp-go-demo/
├── demo/                          # 阿里云 AMQP 专用版本（硬编码配置）
│   ├── Publisher.go              # 消息发布者
│   └── Consumer.go               # 消息消费者
├── amqp091/                       # 通用版本（命令行参数配置）
│   └── pub/
│       ├── producer.go           # 消息生产者（带发布确认）
│       └── producer_annotated.go # 带详细注释的版本
├── Makefile                      # 构建和运行脚本
└── go.mod                        # Go 模块依赖
```

## 🚀 快速开始

### 方式 1：运行 demo/Publisher.go（阿里云 AMQP 专用版）

这是最简单的版本，配置已硬编码在代码中：

```bash
# 使用默认目标
make

# 或明确指定
make run-publisher

# 或直接运行
go run demo/Publisher.go
```

**特点：**
- ✅ 配置简单，无需传参
- ✅ 专为阿里云 AMQP 优化
- ❌ 需要修改代码才能更改配置

---

### 方式 2：运行 amqp091/pub/producer.go（推荐 ⭐）

这是通用版本，支持灵活的命令行参数配置：

#### 🎯 使用默认配置（连接阿里云 AMQP）

```bash
make run-producer
```

这会自动使用预设的阿里云 AMQP 配置。

#### 🔧 自定义配置参数

可以通过环境变量覆盖默认配置：

```bash
# 完整示例
make run-producer \
  URI="amqp://username:password@host:5672/vhost" \
  EXCHANGE="my_exchange" \
  QUEUE="my_queue" \
  KEY="my_routing_key" \
  BODY="Hello World" \
  CONTINUOUS=false
```

#### 📊 可用参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `URI` | AMQP 连接地址 | 阿里云 AMQP 地址 |
| `EXCHANGE` | 交换机名称 | `demo_exchange` |
| `EXCHANGE_TYPE` | 交换机类型 | `direct` |
| `QUEUE` | 队列名称 | `demo_queue` |
| `KEY` | 路由键 | `abc` |
| `BODY` | 消息体内容 | `Hello GO AMQP` |
| `CONTINUOUS` | 是否持续发布 | `false` |

#### 💡 部分覆盖示例

```bash
# 只修改消息内容和路由键
make run-producer \
  KEY="xyz" \
  BODY="Custom Message"

# 切换到其他 vhost
make run-producer \
  URI="amqp://user:pass@host:5672/other_vhost"

# 持续发布模式
make run-producer CONTINUOUS=true
```

---

### 方式 3：直接运行（最灵活）

```bash
# 完全手动指定所有参数
go run amqp091/pub/producer.go \
  -uri="amqp://guest:guest@localhost:5672/" \
  -exchange="test-exchange" \
  -exchange-type="direct" \
  -queue="test-queue" \
  -key="test-key" \
  -body="foobar" \
  -continuous=false
```

---

## 🔨 其他 Make 命令

```bash
# 运行消费者
make run-consumer

# 构建所有可执行文件
make build

# 清理构建产物
make clean
```

---

## 📝 两种版本的对比

| 特性 | demo/Publisher.go | amqp091/pub/producer.go |
|------|-------------------|-------------------------|
| **配置方式** | 硬编码 | 命令行参数 |
| **灵活性** | 低 | 高 |
| **适用场景** | 固定环境测试 | 多环境部署 |
| **发布确认** | ❌ 不支持 | ✅ 支持 |
| **优雅关闭** | ❌ 不支持 | ✅ 支持 |
| **SSL/TLS** | ❌ 不支持 | ✅ 支持 |
| **流量控制** | ❌ 不支持 | ✅ 支持（最多 8 个未完成确认） |
| **日志系统** | 基础 | 完善（INFO/WARN/ERROR） |

---

## 🎓 最佳实践建议

### ✅ 推荐做法

1. **开发测试**：使用 `demo/Publisher.go` 快速验证
2. **生产环境**：使用 `amqp091/pub/producer.go` + Makefile
3. **配置管理**：通过 Makefile 参数管理不同环境配置
4. **监控调试**：利用完善的日志系统跟踪消息状态

### ⚠️ 注意事项

1. **认证信息**：确保 URI 中的用户名和密码格式正确
2. **网络连接**：保证能访问 AMQP 服务器
3. **资源释放**：程序会自动清理连接和通道（defer）
4. **错误处理**：失败时会记录详细错误并退出

---

## 🔍 故障排查

### 连接失败

检查 URI 格式是否正确：
```
amqp://用户名：密码@服务器地址:5672/虚拟主机
```

### 消息未送达

1. 检查 Exchange 和 Queue 是否已声明
2. 检查 Routing Key 是否匹配绑定关系
3. 查看日志输出的错误信息

### 性能问题

调整 `-continuous` 参数或使用流量控制机制

---

## 📚 参考资料

- [RabbitMQ AMQP 0.9.1 协议](https://www.rabbitmq.com/amqp-0-9-1-reference.html)
- [amqp091-go 官方文档](https://pkg.go.dev/github.com/rabbitmq/amqp091-go)
- [阿里云 AMQP 文档](https://help.aliyun.com/product/29630.html)
