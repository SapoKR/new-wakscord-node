package config

var Default Config

type Config struct {
	Host string `env:"HOST" envDefault:"0.0.0.0" json:"-"`
	Port int    `env:"PORT" envDefault:"3000" json:"-"`

	MaxConcurrent    int `env:"MAX_CONCURRENT" envDefault:"500" json:"max_concurrent"`
	WaitConcurrent   int `env:"WAIT_CONCURRENT" envDefault:"1" json:"wait_concurrent"`
	MessageQueueSize int `env:"MESSAGE_QUEUE_SIZE" envDefault:"100" json:"message_queue_size"`

	Key   string `env:"KEY" envDefault:"wakscord" json:"-"`
	ID    int    `env:"ID" envDefault:"0" json:"-"`
	Owner string `env:"OWNER" envDefault:"Unknown" json:"-"`
}
