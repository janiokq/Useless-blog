package cinit

import (
	"github.com/janiokq/Useless-blog/pkg/pprof"
	"github.com/jinzhu/configor"
	"log"
)

const (
	MySQL   = "MySQL"
	Trace   = "Trace"
	Redis   = "Redis"
	Metrics = "Metrics"
	Nats    = "Nats"
	Kafka   = "Kafka"
)

var Config = struct {
	Service struct {
		Name      string `default:""`
		Version   string `default:"v1.0"`
		RateTime  int    `default:"1024"`
		Appkey    string `default:"admin"`
		AppSecret string `default:"admin"`
	}
	//tracing
	Trace struct {
		Address       string  `default:"http://jaeger:14268/api/traces?format=jaeger.thrift"`
		SamplingRate  float64 `default:"1"`
		LogTraceSpans bool    `default:"false"`
	}
	//log config
	Log struct {
		Path         string `default:"tmp"`
		IsStdout     string `default:"yes"`
		MaxAge       int    `default:"7"`
		RotationTime int    `default:"1"`
		MaxSize      int    `default:"100"`
		LogLevel     string `default:"debug"`
	}
	//mysql config
	Mysql struct {
		Dbname   string `default:"useless"`
		Addr     string `default:"127.0.0.1"`
		User     string `default:"root"`
		Password string `default:"root"`
		Port     int    `default:"3306"`
		IDleConn int    `default:"6"`
		MaxConn  int    `default:"40"`
	}
	//redis  config
	Redis struct {
		Addr     string `default:"127.0.0.1:6379"`
		Password string `default:""`
		Db       int    `default:"0"`
	}
	//Nats config
	Nats struct {
		Addr string `default:"127.0.0.1:4443"`
	}
	//kafka
	Kafka struct {
		Addr string `default:"127.0.0.1:9092"`
	}
	//Metrics
	Metrics struct {
		Enable      string `default:"yes"`
		Duration    int    `default:"5"`
		URL         string `default:"http://influxdb:8086"`
		Database    string `default:"test01"`
		UserName    string `default:""`
		Password    string `default:""`
		Measurement string `default:""`
	}

	//User Service
	SrvUser struct {
		Port              string `default:"5001"`
		Address           string `default:"127.0.0.1:5001"`
		GateWayAddr       string `default:"9999"`
		GateWaySwaggerDir string `default:"/swagger"`
	}

	//API Service
	APIBackend struct {
		Port    string `default:":8888"`
		Address string `default:"127.0.0.1:8889"`
	}

	//api Frontend
	APIForntend struct {
		Port    string `default:":8889"`
		Address string `default:"127.0.0.1:8889"`
	}
}{}

func configInit(sn string) {
	err := configor.Load(&Config, "config.yml")
	if err != nil {
		log.Printf("load config error:%+v", err)
		return
	}
	if Config.Service.Name == "" {
		Config.Service.Name = sn
	}
	log.Printf("config: %+v\n", Config)
}

var closeArr []string

func InitOption(sn string, args ...string) {
	//启动pprof
	pprof.Run()
	closeArr = args
	//初始化配置参数
	configInit(sn)
	//初始化日志
	logInit()
	// 3.其他服务
	for _, o := range args {
		switch o {
		case Trace:
			tracerInit()
		case MySQL:
			initMysql()
		case Redis:
			redisInit()
		case Metrics:
			metricsInit(sn)
		case Nats:
			Natsinit()
		case Kafka:
			KafkaInit()

		}
	}
}

// 关闭打开的服务
func Close() {
	for _, o := range closeArr {
		switch o {
		case Trace:
			// 关闭链路跟踪
			tracerClose()
		case MySQL:
			// 关闭mysql
			closeMysql()
		case Redis:
			redisClose()
		case Metrics:
		case Nats:
			Natsclose()
		case Kafka:
			KafkaClose()
		}
	}
}
