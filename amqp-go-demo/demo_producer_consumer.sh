#!/bin/bash

# 完整的生产者 - 消费者演示脚本

echo "=========================================="
echo "AMQP 消息队列演示 - 生产者与消费者"
echo "=========================================="
echo ""

# 配置参数
URI="${URI:-amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo}"
EXCHANGE="${EXCHANGE:-demo_exchange}"
EXCHANGE_TYPE="${EXCHANGE_TYPE:-direct}"
QUEUE="${QUEUE:-demo_queue}"
KEY="${KEY:-abc}"

echo "配置信息:"
echo "  URI: ${URI}"
echo "  Exchange: ${EXCHANGE} (${EXCHANGE_TYPE})"
echo "  Queue: ${QUEUE}"
echo "  Routing Key: ${KEY}"
echo ""

# 启动消费者（后台运行，持续 20 秒）
echo "[1/3] 启动消费者..."
go run amqp091/sub/consumer.go \
    -uri="$URI" \
    -exchange="$EXCHANGE" \
    -exchange-type="$EXCHANGE_TYPE" \
    -queue="$QUEUE" \
    -key="$KEY" \
    -consumer-tag="demo-consumer" \
    -vhost="/demo" \
    -lifetime=20s \
    -verbose=true \
    -auto_ack=false &

CONSUMER_PID=$!
echo "      消费者已启动 (PID: $CONSUMER_PID)"
echo ""

# 等待消费者完全启动
sleep 2

# 发送多条消息
echo "[2/3] 发送测试消息..."
echo ""

MESSAGES=(
    "第一条测试消息 - $(date +%T)"
    "第二条测试消息 - Hello World"
    "第三条测试消息 - AMQP 演示"
    "第四条测试消息 - Go 语言实现"
    "第五条测试消息 - 消息队列测试完成"
)

for msg in "${MESSAGES[@]}"; do
    echo "  发送：$msg"
    go run amqp091/pub/producer.go \
        -uri="$URI" \
        -exchange="$EXCHANGE" \
        -exchange-type="$EXCHANGE_TYPE" \
        -queue="$QUEUE" \
        -key="$KEY" \
        -body="$msg" \
        -continuous=false 2>&1 | grep -E "(publishing|confirmed)" || true
    sleep 1
done

echo ""
echo "[3/3] 等待消费者处理所有消息..."
echo ""

# 等待消费者完成
wait $CONSUMER_PID

echo ""
echo "=========================================="
echo "演示完成！"
echo "=========================================="
