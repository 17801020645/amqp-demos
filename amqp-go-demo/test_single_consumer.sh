#!/bin/bash

echo "========================================"
echo "AMQP 单消费者完整测试"
echo "========================================"
echo ""

# 使用唯一队列名和路由键，避免与共享环境中的其他消费者竞争消息
# RabbitMQ 会轮询分发消息给所有消费者，若使用 demo_queue 会有其他消费者抢走消息
UNIQUE_SUFFIX="$$"
UNIQUE_QUEUE="demo_queue_test_${UNIQUE_SUFFIX}"
UNIQUE_KEY="abc_test_${UNIQUE_SUFFIX}"
echo "[配置] 使用隔离队列: $UNIQUE_QUEUE, 路由键: $UNIQUE_KEY"
echo ""

# 清理之前的进程
echo "[步骤 0] 清理所有残留进程..."
pkill -f "producer.go" 2>/dev/null || true
pkill -f "consumer.go" 2>/dev/null || true
sleep 2

# 确认没有消费者在运行
CONSUMER_COUNT=$(ps aux | grep -E "consumer\.go" | grep -v grep | wc -l)
echo "当前运行的消费者数量：$CONSUMER_COUNT"
if [ $CONSUMER_COUNT -gt 0 ]; then
    echo "⚠️  还有消费者在运行，等待清理..."
    sleep 3
fi

echo ""
echo "========================================"
echo "[步骤 1] 启动唯一的消费者（后台运行）"
echo "========================================"

# 将消费者的输出重定向到临时文件
CONSUMER_LOG=$(mktemp)
echo "消费者日志文件：$CONSUMER_LOG"
make run-consumer-091 QUEUE="$UNIQUE_QUEUE" KEY="$UNIQUE_KEY" > "$CONSUMER_LOG" 2>&1 &
CONSUMER_PID=$!
echo "消费者进程 PID: $CONSUMER_PID"
sleep 3

# 显示消费者启动日志
echo ""
echo "--- 消费者启动日志 ---"
cat "$CONSUMER_LOG"
echo "----------------------"

# 检查队列中的消费者数量
CONSUMER_IN_LOG=$(grep "consumers" "$CONSUMER_LOG" | tail -1 | grep -o '[0-9] consumers')
echo "检测到的消费者数量：$CONSUMER_IN_LOG"

echo ""
echo "========================================"
echo "[步骤 2] 发送消息（生产者）"
echo "========================================"
echo ""
echo ">>> 发送第 1 条消息..."
make run-producer QUEUE="$UNIQUE_QUEUE" KEY="$UNIQUE_KEY" BODY="独占消息 1"
sleep 1

echo ""
echo ">>> 发送第 2 条消息..."
make run-producer QUEUE="$UNIQUE_QUEUE" KEY="$UNIQUE_KEY" BODY="独占消息 2"
sleep 1

echo ""
echo ">>> 发送第 3 条消息..."
make run-producer QUEUE="$UNIQUE_QUEUE" KEY="$UNIQUE_KEY" BODY="独占消息 3"
sleep 1

echo ""
echo "========================================"
echo "[步骤 3] 等待消费者处理消息"
echo "========================================"
sleep 2

# 显示消费者处理日志
echo ""
echo "--- 消费者完整日志 ---"
cat "$CONSUMER_LOG"
echo "----------------------"

# 统计收到的消息数量
MSG_COUNT=$(grep "got.*delivery:" "$CONSUMER_LOG" | wc -l)
echo ""
echo "📊 统计：消费者共收到 $MSG_COUNT 条消息"

echo ""
echo "========================================"
echo "[步骤 4] 停止消费者"
echo "========================================"
kill $CONSUMER_PID 2>/dev/null || true
wait $CONSUMER_PID 2>/dev/null || true

# 显示最终日志
echo ""
echo "--- 消费者最终日志 ---"
cat "$CONSUMER_LOG"
echo "----------------------"

# 清理临时文件
rm -f "$CONSUMER_LOG"

echo ""
echo "========================================"
echo "✅ 测试完成！"
echo "========================================"
