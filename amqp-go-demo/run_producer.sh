#!/bin/bash

# 运行生产者脚本 - 连接阿里云 AMQP

go run amqp091/pub/producer.go \
  -uri="amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo" \
  -exchange="demo_exchange" \
  -exchange-type="direct" \
  -queue="demo_queue" \
  -key="abc" \
  -body="Hello GO AMQP" \
  -continuous=false
