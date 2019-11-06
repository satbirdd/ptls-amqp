package ptls_amqp

import "fmt"

const (
	KeyArchive        = "archive"
	KeyUnarchive      = "unarchive"
	KeyOsd            = "osd"
	KeyUploadToLandon = "uploadToLandon"
	KeyLog            = "log"

	ExchangeKindDirect = "direct"
	ExchangeKindTopic  = "topic"
	ExchangeKindFanout = "fanout"
	ExchangeKindHeader = "header"

	AmqpUrlFmt = "amqp://%v:%v@%v%v"
)

type Config struct {
	User  string
	Pass  string
	Host  string
	Vhost string
}

func (cfg Config) String() string {
	return fmt.Sprintf(AmqpUrlFmt, cfg.User, cfg.Pass, cfg.Host, cfg.Vhost)
}

// // 管理RabbitMQ Connection，在需要的时候进行重连等操作
// type Connection struct {
// 	m              *sync.Mutex
// 	config         Config
// 	amqpConnection *amqp.Connection
// 	reconnectTimes int32
// 	channels       []*Channel
// }

// type Channel struct {
// 	connection  *Connection
// 	amqpChannel *amqp.Channel
// }

// func NewConnection(cfg Config) (*Connection, error) {
// 	con := Connection{
// 		config: &cfg,
// 		m:      &sync.Mutex{},
// 	}

// 	connection, err := amqp.Dial(cfg.String())
// 	if err != nil {
// 		return nil, fmt.Errorf("无法连接到RabbitMQ,", err)
// 	}

// 	con.amqpConnection = connection
// }

// func (con *Connection) Reconnect() error {
// 	con.m.Lock()
// 	defer con.m.Unlock()

// 	if !con.amqpConnection.IsClosed() {
// 		con.amqpConnection.Close()
// 	}

// 	connection, err := amqp.Dial(cfg.String())
// 	if err != nil {
// 		return fmt.Errorf("无法连接到RabbitMQ,%v", err)
// 	}

// 	// rbtchannel, err := connection.Channel()
// 	// if err != nil {
// 	// 	return fmt.Errorf("无法建立RabbitMQ频道,%v", err)
// 	// }

// 	con.amqpConnection = connection
// 	con.reconnectTimes += 1
// }
