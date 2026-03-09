# AMQP Go Demo 项目

基于 Go 语言的 RabbitMQ/阿里云 AMQP 消息队列示例项目。

## 📋 项目结构

```
amqp-go-demo/
├── demo/                          # 阿里云 AMQP 专用版本（硬编码配置）
│   ├── Publisher.go              # 消息发布者 - 发送消息
│   └── Consumer.go               # 消息消费者 - 接收消息
├── amqp091/                       # 通用版本（命令行参数配置）
│   └── pub/
│       ├── producer.go           # 消息生产者（带发布确认机制）
│       └── producer_annotated.go # 带详细注释的版本
├── Makefile                      # 构建和运行脚本
├── README.md                     # 项目说明（本文档）
├── README_USAGE.md               # 详细使用指南
├── TESTING.md                    # 快速测试指南
└── run_producer.sh               # 快速启动脚本
```

## 🎯 特性对比

| 特性 | demo/ 版本 | amqp091/ 版本 |
|------|-----------|--------------|
| **配置方式** | 硬编码 | 命令行参数 |
| **适用场景** | 固定环境、快速测试 | 多环境、生产部署 |
| **发布确认** | ❌ | ✅ |
| **优雅关闭** | ❌ | ✅ |
| **SSL/TLS** | ❌ | ✅ |
| **流量控制** | ❌ | ✅ |
| **日志系统** | 基础 | 完善（INFO/WARN/ERROR） |

## 🚀 快速开始

### 1. 环境准备

确保已安装：
- Go 1.19+
- RabbitMQ 或阿里云 AMQP 服务

### 2. 安装依赖

```bash
go mod download
```

### 3. 运行消息发布者

```bash
# 方式 1：使用阿里云专用版本（推荐用于快速测试）
make run-publisher

# 方式 2：使用通用版本（推荐用于生产环境）
make run-producer
```

### 4. 运行消息消费者

在另一个终端窗口：

```bash
make run-consumer
```

## 📖 使用说明

### 基础用法

#### 发送消息（终端 1）
```bash
make run-publisher
```

#### 接收消息（终端 2）
```bash
make run-consumer
```

### 高级用法

#### 自定义消息内容
```bash
make run-producer BODY="Custom Message" KEY="my_key"
```

#### 持续发布模式
```bash
make run-producer CONTINUOUS=true
```

#### 完全自定义配置
```bash
make run-producer \
  URI="amqp://user:pass@host:5672/vhost" \
  EXCHANGE="my_exchange" \
  QUEUE="my_queue" \
  KEY="routing_key" \
  BODY="My Message"
```

## 🔧 Makefile 命令

| 命令 | 说明 |
|------|------|
| `make` | 运行默认 Publisher |
| `make run-publisher` | 运行阿里云专用 Publisher |
| `make run-consumer` | 运行阿里云专用 Consumer |
| `make run-producer` | 运行通用 Producer（带确认） |
| `make build` | 构建所有可执行文件 |
| `make clean` | 清理构建产物 |

## 📚 文档索引

- **[README_USAGE.md](README_USAGE.md)** - 详细使用指南
  - 完整的参数说明
  - 两种版本的对比
  - 最佳实践建议
  - 故障排查指南

- **[TESTING.md](TESTING.md)** - 快速测试指南
  - 测试场景示例
  - Producer + Consumer 配合测试
  - 性能测试方法

## 🎓 核心概念

### AMQP 架构

```
Publisher → Exchange → Queue → Consumer
            (RouteKey)
```

### 关键组件

1. **Exchange（交换机）**
   - 接收生产者发送的消息
   - 根据路由规则分发到队列
   - 类型：direct, fanout, topic, headers

2. **Queue（队列）**
   - 存储消息的缓冲区
   - 消费者从队列中获取消息
   - 支持持久化

3. **RouteKey（路由键）**
   - 决定消息如何从 Exchange 路由到 Queue
   - 不同类型的 Exchange 有不同的匹配规则

## 💡 最佳实践

### ✅ 推荐做法

1. **开发测试**
   - 使用 `demo/Publisher.go` 快速验证
   - 配置简单，无需传参

2. **生产环境**
   - 使用 `amqp091/pub/producer.go`
   - 启用发布确认机制
   - 通过环境变量管理配置

3. **监控调试**
   - 利用完善的日志系统
   - 关注发布确认状态
   - 定期检查队列积压

### ⚠️ 注意事项

1. **连接管理**
   - 使用长连接，避免频繁创建/销毁
   - 正确关闭连接和资源（defer）

2. **消息可靠性**
   - 启用发布确认（Publisher Confirms）
   - 设置消息持久化（DeliveryMode = Persistent）
   - 队列持久化（durable = true）

3. **错误处理**
   - 检查所有操作的返回值
   - 记录详细的错误日志
   - 实现重试机制

## 🔍 故障排查

### 常见问题

#### 1. 连接失败
```bash
# 检查 URI 格式
amqp://用户名：密码@服务器地址:5672/虚拟主机

# 检查网络连通性
telnet rabbitmq-cn-xxx.mq.amqp.aliyuncs.com 5672
```

#### 2. 消息未送达
- 检查 Exchange 和 Queue 是否已声明
- 检查 Routing Key 是否匹配
- 查看绑定关系是否正确

#### 3. 消息堆积
```bash
# 查看队列中的消息数量
# 日志中会显示："queue_name" X messages

# 启动更多消费者
make run-consumer
```

## 🛠️ 开发指南

### 添加新的生产者

参考 `demo/Publisher.go` 或 `amqp091/pub/producer.go`

### 添加新的消费者

参考 `demo/Consumer.go`

### 自定义配置

修改对应文件的配置部分，或通过命令行参数传递

## 📊 性能优化

### 批量发送

```bash
# 使用持续发布模式测试吞吐量
make run-producer CONTINUOUS=true
```

### 并发消费

启动多个 Consumer 实例：
```bash
# 终端 1
make run-consumer

# 终端 2
make run-consumer

# 终端 3
make run-consumer
```

## 🔗 相关链接

- [RabbitMQ 官方文档](https://www.rabbitmq.com/documentation.html)
- [AMQP 0.9.1 协议规范](https://www.rabbitmq.com/amqp-0-9-1-reference.html)
- [amqp091-go GitHub](https://github.com/rabbitmq/amqp091-go)
- [阿里云 AMQP 文档](https://help.aliyun.com/product/29630.html)

## 📄 许可证

本项目仅用于学习和演示目的。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**最后更新时间**: 2026-03-09
