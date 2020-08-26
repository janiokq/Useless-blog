package configlog

//日志级别
type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

func (level Level) String() string {
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "Fata"
	}
	return "unknown"
}

var AllLevels = []Level{
	PanicLevel,
	FatalLevel,
	ErrorLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

//默认参数
const (
	LogPath     string = "/var/log" //日志保存路径
	LogName     string = "output"   //日志保存的名称，不些随机生成
	LogLevel    string = "debug"    //日志记录级别
	MaxSize     int    = 100        //日志分割的尺寸 MB
	MaxAge      int    = 7          //分割日志保存的时间 day
	Stacktrace  string = "error"    //记录堆栈的级别
	IsStdOut    string = "yes"      //是否标准输出console输出
	ProjectName string = "test"     //项目名称

)
