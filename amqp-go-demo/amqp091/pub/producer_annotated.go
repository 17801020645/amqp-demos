// AMQP 消息发布者 - 使用 RabbitMQ AMQP 0.9.1 协议
// 该程序演示了如何使用发布确认模式可靠地发送消息到 RabbitMQ
package main

import (
	"crypto/tls"              // TLS/SSL 加密支持
	"flag"                    // 命令行参数解析
	amqp "github.com/rabbitmq/amqp091-go" // RabbitMQ AMQP 0.9.1 Go 客户端库
	"log"                     // 日志记录
	"os"                      // 操作系统功能
	"os/signal"               // 信号处理
	"strings"                 // 字符串操作
	"syscall"                 // 系统调用
	"time"                    // 时间相关功能
)

// 全局变量定义 - 命令行参数和日志器
var (
	// AMQP 连接地址，默认使用本地测试环境
	uri = flag.String("uri", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	
	// Exchange 名称（持久化）
	exchange = flag.String("exchange", "test-exchange", "Durable AMQP exchange name")
	
	// Exchange 类型：direct(直连)、fanout(广播)、topic(主题匹配)、x-custom(自定义)
	exchangeType = flag.String("exchange-type", "direct", "Exchange type - direct|fanout|topic|x-custom")
	
	// Queue 名称（临时队列）
	queue = flag.String("queue", "test-queue", "Ephemeral AMQP queue name")
	
	// 路由键 - 用于将消息路由到指定的队列
	routingKey = flag.String("key", "test-key", "AMQP routing key")
	
	// 消息体内容
	body = flag.String("body", "foobar", "Body of message")
	
	// 虚拟主机 - 用于隔离不同的应用或租户
	vhost = flag.String("vhost", "/", "vhost")
	
	// 是否持续发布模式 - 以每秒 1 条消息的速度持续发送
	continuous = flag.Bool("continuous", false, "Keep publishing messages at a 1msg/sec rate")
	
	// 警告日志 - 输出到标准错误流，带 [WARNING] 前缀
	WarnLog = log.New(os.Stderr, "[WARNING] ", log.LstdFlags|log.Lmsgprefix)
	
	// 错误日志 - 输出到标准错误流，带 [ERROR] 前缀
	ErrLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lmsgprefix)
	
	// 信息日志 - 输出到标准输出流，带 [INFO] 前缀
	Log = log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lmsgprefix)
)

// init 函数 - 在 main 之前执行，用于解析命令行参数
func init() {
	flag.Parse()
}

// main 函数 - 程序入口
func main() {
	// 创建退出信号通道 - 用于优雅关闭
	exitCh := make(chan struct{})
	
	// 创建延迟确认通道 - 接收发布确认
	confirmsCh := make(chan *amqp.DeferredConfirmation)
	
	// 创建确认完成通道 - 表示所有确认已完成
	confirmsDoneCh := make(chan struct{})
	
	// 注意：这是一个缓冲通道，这样发送 OK 信号不会阻塞确认处理器
	// 缓冲大小为 1，允许一个待处理的发布许可
	publishOkCh := make(chan struct{}, 1)

	// 设置关闭信号处理器 - 监听 Ctrl+C 等终止信号
	setupCloseHandler(exitCh)

	// 启动确认处理器 - 异步处理消息确认
	startConfirmHandler(publishOkCh, confirmsCh, confirmsDoneCh, exitCh)

	// 开始发布消息 - 主发布循环
	publish(publishOkCh, confirmsCh, confirmsDoneCh, exitCh)
}

// setupCloseHandler - 设置关闭信号处理器
// 用于捕获操作系统的中断信号（如 Ctrl+C），实现优雅关闭
func setupCloseHandler(exitCh chan struct{}) {
	// 创建信号通道，缓冲区大小为 2
	c := make(chan os.Signal, 2)
	
	// 注册需要监听的信号：Ctrl+C (SIGINT) 和 终止信号 (SIGTERM)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// 启动 goroutine 监听信号
	go func() {
		<-c // 等待信号到达
		Log.Printf("close handler: Ctrl+C pressed in Terminal") // 记录日志
		close(exitCh) // 关闭退出通道，通知其他 goroutine 停止
	}()
}

// publish - 主要的消息发布函数
// 负责建立连接、声明队列和交换机、发布消息
func publish(publishOkCh <-chan struct{}, confirmsCh chan<- *amqp.DeferredConfirmation, confirmsDoneCh <-chan struct{}, exitCh chan struct{}) {
	// 配置 AMQP 连接参数
	config := amqp.Config{
		Vhost:      *vhost,                        // 设置虚拟主机
		Properties: amqp.NewConnectionProperties(), // 创建连接属性
	}

	// 设置客户端连接名称，便于在管理界面识别
	config.Properties.SetClientConnectionName("producer-with-confirms")
	
	// 记录正在尝试连接的 URI
	Log.Printf("producer: dialing %s", *uri)

	var conn *amqp.Connection // AMQP 连接对象
	var err error             // 错误处理变量
	
	// 根据 URI 判断是否使用 SSL/TLS 加密连接
	if strings.Contains(*uri, "amqps://") {
		// SSL/TLS 连接配置
		tlsfg := &tls.Config{InsecureSkipVerify: true} // 跳过证书验证（生产环境应配置正确证书）
		conn, err = amqp.DialTLS(*uri, tlsfg)          // 建立 TLS 连接
		if err != nil {
			ErrLog.Fatalf("producer: error in TLS dial: %s", err) // TLS 连接失败，记录错误并退出
		}
	} else {
		// 普通非加密连接
		conn, err = amqp.DialConfig(*uri, config) // 使用配置建立连接
		if err != nil {
			ErrLog.Fatalf("producer: error in dial: %s", err) // 连接失败，记录错误并退出
		}
	}

	defer conn.Close() // 确保函数退出时关闭连接

	Log.Println("producer: got Connection, getting Channel") // 记录已获取连接，准备获取通道
	
	// 创建 AMQP 通道 - 所有操作都通过通道进行
	channel, err := conn.Channel()
	if err != nil {
		ErrLog.Fatalf("error getting a channel: %s", err) // 创建通道失败，记录错误并退出
	}
	defer channel.Close() // 确保函数退出时关闭通道

	Log.Printf("producer: declaring exchange") // 记录正在声明交换机
	
	// 声明交换机 - 如果不存在则创建
	if err := channel.ExchangeDeclare(
		*exchange,     // 交换机名称
		*exchangeType, // 交换机类型（direct/fanout/topic/x-custom）
		true,          // durable: 持久化，RabbitMQ 重启后仍然存在
		false,         // auto-delete: 不自动删除（最后一个队列解绑后删除）
		false,         // internal: 不是内部交换机（可直接发布消息）
		false,         // noWait: 等待服务器响应
		nil,           // arguments: 无额外参数
	); err != nil {
		ErrLog.Fatalf("producer: Exchange Declare: %s", err) // 声明交换机失败，记录错误并退出
	}

	Log.Printf("producer: declaring queue '%s'", *queue) // 记录正在声明队列
	
	// 声明队列
	queue, err := channel.QueueDeclare(
		*queue, // 队列名称
		true,   // durable: 持久化队列
		false,  // delete when unused: 不使用完不删除
		false,  // exclusive: 非独占队列（可被多个消费者共享）
		false,  // noWait: 等待服务器响应
		nil,    // arguments: 无额外参数
	)
	if err == nil {
		// 队列声明成功，记录队列信息
		Log.Printf("producer: declared queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
			queue.Name, queue.Messages, queue.Consumers, *routingKey)
	} else {
		ErrLog.Fatalf("producer: Queue Declare: %s", err) // 声明队列失败，记录错误并退出
	}

	Log.Printf("producer: declaring binding") // 记录正在声明绑定关系
	
	// 绑定队列到交换机 - 建立路由规则
	if err := channel.QueueBind(queue.Name, *routingKey, *exchange, false, nil); err != nil {
		ErrLog.Fatalf("producer: Queue Bind: %s", err) // 绑定失败，记录错误并退出
	}

	// 启用发布确认模式 - 确保消息可靠投递
	// 可靠的发布者确认需要服务器支持 confirm.select
	Log.Printf("producer: enabling publisher confirms.")
	if err := channel.Confirm(false); err != nil {
		ErrLog.Fatalf("producer: channel could not be put into confirm mode: %s", err) // 启用确认模式失败
	}

	// 主发布循环 - 持续发布消息
	for {
		canPublish := false // 发布许可标志
		
		Log.Println("producer: waiting on the OK to publish...") // 等待发布许可
		
		// 内层循环 - 等待发布许可或退出信号
		for {
			select {
			case <-confirmsDoneCh: // 收到确认完成信号
				Log.Println("producer: stopping, all confirms seen")
				return // 退出发布函数
			case <-publishOkCh: // 收到发布许可
				Log.Println("producer: got the OK to publish")
				canPublish = true
				break
			case <-time.After(time.Second): // 超时（1 秒）
				WarnLog.Println("producer: still waiting on the OK to publish...")
				continue // 继续等待
			}
			if canPublish {
				break // 获得许可，跳出等待循环
			}
		}

		// 记录即将发布的消息信息
		Log.Printf("producer: publishing %dB body (%q)", len(*body), *body)
		
		// 发布消息并使用延迟确认
		dConfirmation, err := channel.PublishWithDeferredConfirm(
			*exchange, // 交换机名称
			*routingKey, // 路由键
			true,       // mandatory: 必须投递到队列，否则返回错误
			false,      // immediate: 不立即投递（允许排队）
			amqp.Publishing{ // 消息属性
				Headers:         amqp.Table{},        // 消息头（空）
				ContentType:     "text/plain",        // 内容类型：纯文本
				ContentEncoding: "",                  // 内容编码：无
				DeliveryMode:    amqp.Persistent,     // 持久化消息（保存到磁盘）
				Priority:        0,                   // 优先级：0（最低）
				AppId:           "sequential-producer", // 应用 ID
				Body:            []byte(*body),       // 消息体
			},
		)
		if err != nil {
			ErrLog.Fatalf("producer: error in publish: %s", err) // 发布失败
		}

		// 将延迟确认对象发送给确认处理器
		select {
		case <-confirmsDoneCh: // 检查是否需要退出
			Log.Println("producer: stopping, all confirms seen")
			return
		case confirmsCh <- dConfirmation: // 发送到确认通道
			Log.Println("producer: delivered deferred confirm to handler")
			break
		}

		// 控制发布节奏
		select {
		case <-confirmsDoneCh: // 检查是否需要退出
			Log.Println("producer: stopping, all confirms seen")
			return
		case <-time.After(time.Millisecond * 250): // 等待 250 毫秒
			if *continuous {
				continue // 持续模式：继续发布下一条
			} else {
				// 非持续模式：准备停止
				Log.Println("producer: initiating stop")
				close(exitCh) // 发送退出信号
				
				// 等待确认完成，最多等待 10 秒
				select {
				case <-confirmsDoneCh:
					Log.Println("producer: stopping, all confirms seen")
					return
				case <-time.After(time.Second * 10):
					WarnLog.Println("producer: may be stopping with outstanding confirmations")
					return
				}
			}
		}
	}
}

// startConfirmHandler - 启动确认处理器（goroutine）
// 负责跟踪和管理所有未确认的消息
func startConfirmHandler(publishOkCh chan<- struct{}, confirmsCh <-chan *amqp.DeferredConfirmation, confirmsDoneCh chan struct{}, exitCh <-chan struct{}) {
	go func() {
		// 创建映射来跟踪未确认的消息 - key 是交付标签，value 是延迟确认对象
		confirms := make(map[uint64]*amqp.DeferredConfirmation)

		// 确认处理主循环
		for {
			select {
			case <-exitCh: // 收到退出信号
				exitConfirmHandler(confirms, confirmsDoneCh) // 清理并退出
				return
			default:
				break // 继续执行
			}

			// 计算未完成的确认数量
			outstandingConfirmationCount := len(confirms)

			// 注意：8 是任意值，你可以根据需要调整允许的未完成确认数量
			// 如果未完成确认数 <= 8，允许继续发布
			if outstandingConfirmationCount <= 8 {
				select {
				case publishOkCh <- struct{}{}: // 发送发布许可
					Log.Println("confirm handler: sent OK to publish")
				case <-time.After(time.Second * 5): // 5 秒超时
					WarnLog.Println("confirm handler: timeout indicating OK to publish (this should never happen!)")
				}
			} else {
				// 未完成确认太多，阻塞发布并记录警告
				WarnLog.Printf("confirm handler: waiting on %d outstanding confirmations, blocking publish", outstandingConfirmationCount)
			}

			// 等待新的确认请求或退出信号
			select {
			case confirmation := <-confirmsCh: // 收到新的延迟确认
				dtag := confirmation.DeliveryTag // 获取交付标签
				confirms[dtag] = confirmation    // 添加到跟踪映射
			case <-exitCh: // 收到退出信号
				exitConfirmHandler(confirms, confirmsDoneCh)
				return
			}

			// 检查已完成的确认
			checkConfirmations(confirms)
		}
	}()
}

// exitConfirmHandler - 退出确认处理器
// 等待所有未完成的确认完成后关闭
func exitConfirmHandler(confirms map[uint64]*amqp.DeferredConfirmation, confirmsDoneCh chan struct{}) {
	Log.Println("confirm handler: exit requested") // 记录退出请求
	waitConfirmations(confirms)                    // 等待所有确认完成
	close(confirmsDoneCh)                          // 关闭确认完成通道
	Log.Println("confirm handler: exiting")        // 记录正在退出
}

// checkConfirmations - 检查并处理已完成的确认
// 遍历所有未确认消息，移除已确认的
func checkConfirmations(confirms map[uint64]*amqp.DeferredConfirmation) {
	Log.Printf("confirm handler: checking %d outstanding confirmations", len(confirms))
	for k, v := range confirms {
		if v.Acked() { // 检查是否已确认
			Log.Printf("confirm handler: confirmed delivery with tag: %d", k)
			delete(confirms, k) // 从跟踪映射中移除
		}
	}
}

// waitConfirmations - 等待所有未完成的确认完成
// 用于优雅关闭时确保所有消息都被确认
func waitConfirmations(confirms map[uint64]*amqp.DeferredConfirmation) {
	Log.Printf("confirm handler: waiting on %d outstanding confirmations", len(confirms))

	// 先检查一遍当前状态
	checkConfirmations(confirms)

	// 等待每个未完成的确认
	for k, v := range confirms {
		select {
		case <-v.Done(): // 等待确认完成
			Log.Printf("confirm handler: confirmed delivery with tag: %d", k)
			delete(confirms, k)
		case <-time.After(time.Second): // 1 秒超时
			WarnLog.Printf("confirm handler: did not receive confirmation for tag %d", k)
		}
	}

	// 检查是否还有未完成的确认
	outstandingConfirmationCount := len(confirms)
	if outstandingConfirmationCount > 0 {
		ErrLog.Printf("confirm handler: exiting with %d outstanding confirmations", outstandingConfirmationCount)
	} else {
		Log.Println("confirm handler: done waiting on outstanding confirmations")
	}
}
