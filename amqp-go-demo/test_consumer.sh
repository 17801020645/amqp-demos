#!/bin/bash

# 测试 consumer 是否可以正常启动和消费消息

echo "=== 启动 Consumer (运行 5 秒后自动停止) ==="

# 在后台启动 consumer
go run amqp091/sub/consumer.go \
    -uri="amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo" \
    -exchange="demo_exchange" \
    -exchange-type="direct" \
    -queue="demo_queue" \
    -key="abc" \
    -consumer-tag="test-consumer" \
    -vhost="/demo" \
    -lifetime=5s \
    -verbose=true \
    -auto_ack=false &

CONSUMER_PID=$!

echo "Consumer 已启动 (PID: $CONSUMER_PID)，等待 2 秒..."
sleep 2

echo ""
echo "=== 发送测试消息 ==="
go run amqp091/pub/producer.go \
    -uri="amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo" \
    -exchange="demo_exchange" \
    -exchange-type="direct" \
    -queue="demo_queue" \
    -key="abc" \
    -body="Test message for consumer" \
    -continuous=false

echo ""
echo "=== 等待 Consumer 处理消息 ==="
sleep 3

echo ""
echo "=== 测试完成 ==="
