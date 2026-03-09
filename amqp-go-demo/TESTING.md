# AMQP Go Demo 快速测试指南

## 🚀 快速测试 - Producer 和 Consumer 配合

### 测试场景 1：基础消息收发

#### 步骤 1：启动消费者（终端 1）

```bash
# 使用 Makefile
make run-consumer

# 或直接运行
go run demo/Consumer.go
```

你会看到：
```
2026/03/09 15:04:34 [*] Waiting for messages. To exit press CTRL+C
```

#### 步骤 2：发送消息（终端 2）

```bash
# 使用默认配置发送 10 条消息
make run-publisher
```

你会看到：
```
Go Message [1] sent and confirmed.
Go Message [2] sent and confirmed.
...
Go Message [10] sent and confirmed.
```

#### 步骤 3：观察消费者接收

回到终端 1，你会看到：
```
2026/03/09 15:04:34 [*] Waiting for messages. To exit press CTRL+C
2026/03/09 15:04:35 Received a message: Hello GO AMQP[1]
2026/03/09 15:04:36 Received a message: Hello GO AMQP[2]
...每秒一条...
```

---

### 测试场景 2：自定义消息内容

#### 步骤 1：启动消费者（终端 1）

```bash
make run-consumer
```

#### 步骤 2：发送自定义消息（终端 2）

```bash
# 发送单条自定义消息
make run-producer \
  BODY="Hello from Custom Test" \
  KEY="test_key"
```

#### 步骤 3：观察结果

消费者会收到：
```
Received a message: Hello from Custom Test
```

---

### 测试场景 3：持续发布模式

#### 步骤 1：启动消费者（终端 1）

```bash
make run-consumer
```

#### 步骤 2：持续发布消息（终端 2）

```bash
# 启用持续发布模式，每秒发送一条
make run-producer CONTINUOUS=true
```

这会持续发送消息，直到你按 `Ctrl+C` 停止。

#### 步骤 3：停止发布

在终端 2 按 `Ctrl+C` 停止发布。

---

## 📊 完整命令参考

### Makefile 命令

| 命令 | 说明 | 示例 |
|------|------|------|
| `make` | 运行默认 Publisher | `make` |
| `make run-publisher` | 运行阿里云专用 Publisher | `make run-publisher` |
| `make run-consumer` | 运行阿里云专用 Consumer | `make run-consumer` |
| `make run-producer` | 运行通用 Producer（带确认） | `make run-producer` |
| `make build` | 构建所有可执行文件 | `make build` |
| `make clean` | 清理构建产物 | `make clean` |

### 自定义参数

```bash
# 完整参数示例
make run-producer \
  URI="amqp://user:pass@host:5672/vhost" \
  EXCHANGE="my_exchange" \
  EXCHANGE_TYPE="direct" \
  QUEUE="my_queue" \
  KEY="my_key" \
  BODY="My Message" \
  CONTINUOUS=false
```

### 常用组合

```bash
# 只修改消息内容
make run-producer BODY="Custom Message"

# 只修改路由键
make run-producer KEY="another_key"

# 切换到其他 vhost
make run-producer \
  URI="amqp://...@host:5672/other_vhost"

# 持续发布模式
make run-producer CONTINUOUS=true
```

---

## 🔍 故障排查

### Consumer 收不到消息？

1. **检查配置是否一致**
   ```bash
   # Publisher 和 Consumer 必须使用相同的：
   # - vhost
   # - exchange
   # - queue
   # - routing key
   ```

2. **检查队列中是否有消息**
   ```bash
   # 使用 make run-producer 查看队列状态
   # 日志中会显示："demo_queue" X messages
   ```

3. **重新启动 Consumer**
   ```bash
   # 停止当前的 Consumer (Ctrl+C)
   make run-consumer
   ```

### 消息堆积怎么办？

如果 Producer 显示队列中有消息堆积：
```
"demo_queue" 3 messages, 0 consumers
```

说明没有 Consumer 在监听，启动 Consumer 即可：
```bash
make run-consumer
```

### 如何优雅停止？

- **Producer**: 自动停止（发送完成后）
- **Consumer**: 手动按 `Ctrl+C` 停止

---

## 🎯 最佳实践

### ✅ 推荐做法

1. **开发测试流程**
   ```bash
   # 终端 1：先启动 Consumer
   make run-consumer &
   
   # 终端 2：运行 Producer 测试
   make run-publisher
   ```

2. **使用通用版本进行压力测试**
   ```bash
   # 持续发布模式
   make run-producer CONTINUOUS=true
   ```

3. **构建可执行文件用于部署**
   ```bash
   make build
   # 生成的文件在 bin/ 目录
   ```

### ⚠️ 注意事项

1. **资源清理**：Consumer 需要手动停止（Ctrl+C）
2. **连接复用**：两个版本都使用长连接，无需频繁创建
3. **消息确认**：Consumer 使用自动确认（auto-ack）
4. **持久化**：消息和队列都是持久化的，重启后仍存在

---

## 📈 性能测试示例

```bash
# 启动 Consumer
make run-consumer

# 在另一个终端，快速发送大量消息
for i in {1..100}; do
  make run-producer BODY="Message $i" > /dev/null 2>&1 &
done
wait

# 观察 Consumer 的处理速度
```

这可以测试系统的并发处理能力！
