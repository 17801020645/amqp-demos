.PHONY: run-publisher run-consumer run-producer clean build

# 运行发布者 (demo/Publisher.go - 阿里云 AMQP 专用版)
run-publisher:
	go run demo/Publisher.go

# 运行消费者 (demo/Consumer.go - 阿里云 AMQP 专用版)
run-consumer:
	go run demo/Consumer.go

# 运行生产者 (amqp091/pub/producer.go - 通用版本，带发布确认)
# 用法：make run-producer URI="amqp://..." EXCHANGE="demo_exchange" QUEUE="demo_queue" KEY="abc"
run-producer:
	go run amqp091/pub/producer.go \
		-uri="$(or $(URI),amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo)" \
		-exchange="$(or $(EXCHANGE),demo_exchange)" \
		-exchange-type="$(or $(EXCHANGE_TYPE),direct)" \
		-queue="$(or $(QUEUE),demo_queue)" \
		-key="$(or $(KEY),abc)" \
		-body="$(or $(BODY),Hello GO AMQP)" \
		-continuous=$(or $(CONTINUOUS),false)

# 运行消费者 (amqp091/sub/consumer.go - 通用版本)
# 用法：make run-consumer URI="amqp://..." EXCHANGE="demo_exchange" QUEUE="demo_queue" KEY="abc"
run-consumer-091:
	go run amqp091/sub/consumer.go \
		-uri="$(or $(URI),amqp://MjpyYWJiaXRtcS1jbi04bGw0b3NxbmYwMTpMVEFJNXQ2REY0RDZkVTgyaERYRlI0QUw=:MkMxN0Q5REE4RDQ3MzNCRDJCNTlDNTI0NTU3RUNDRDA4MkQ4NzE0NzoxNzczMDI1Nzc0Njcx@rabbitmq-cn-8ll4osqnf01.cn-beijing.amqp-82.net.mq.amqp.aliyuncs.com:5672/demo)" \
		-exchange="$(or $(EXCHANGE),demo_exchange)" \
		-exchange-type="$(or $(EXCHANGE_TYPE),direct)" \
		-queue="$(or $(QUEUE),demo_queue)" \
		-key="$(or $(KEY),abc)" \
		-consumer-tag="$(or $(CONSUMER_TAG),simple-consumer)" \
		-lifetime="$(or $(LIFETIME),0s)" \
		-verbose="$(or $(VERBOSE),true)" \
		-auto_ack="$(or $(AUTO_ACK),false)"

# 构建可执行文件
build:
	go build -o bin/publisher demo/Publisher.go
	go build -o bin/consumer demo/Consumer.go
	go build -o bin/producer amqp091/pub/producer.go
	go build -o bin/consumer-091 amqp091/sub/consumer.go

# 清理构建文件
clean:
	rm -rf bin/

# 默认目标
default: run-publisher
