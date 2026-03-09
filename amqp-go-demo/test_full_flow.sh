#!/bin/bash

echo "========================================"
echo "AMQP 生产 - 消费完整测试"
echo "========================================"
echo ""

# 清理之前的进程
echo "[步骤 0] 确保没有残留进程..."
pkill -f "producer.go" 2>/dev/null || true
pkill -f "consumer.go" 2>/dev/null || true
sleep 1

echo ""
echo "========================================"
echo "[步骤 1] 启动消费者（后台运行）"
echo "========================================"
make run-consumer-091 &
CONSUMER_PID=$!
echo "消费者进程 PID: $CONSUMER_PID"
sleep 3

echo ""
echo "========================================"
echo "[步骤 2] 发送消息（生产者）"
echo "========================================"
make run-producer BODY="测试消息 1"
sleep 1

make run-producer BODY="测试消息 2"
sleep 1

make run-producer BODY="测试消息 3"
sleep 1

echo ""
echo "========================================"
echo "[步骤 3] 等待消费者处理消息"
echo "========================================"
sleep 2

echo ""
echo "========================================"
echo "[步骤 4] 停止消费者"
echo "========================================"
kill $CONSUMER_PID 2>/dev/null || true
wait $CONSUMER_PID 2>/dev/null || true

echo ""
echo "========================================"
echo "测试完成！"
echo "========================================"
