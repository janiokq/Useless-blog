package configlog

type Option func(*Options)

type Options struct {
	LogPath     string
	LogName     string
	LogLevel    string
	MaxSize     int
	MaxAge      int
	Stacktrace  string
	IsStdOut    string
	ProjectName string
}

func WithLogPath(logpath string) Option {
	return func(options *Options) {
		options.LogPath = logpath
	}
}

func WithLogName(logname string) Option {
	return func(options *Options) {
		options.LogName = logname
	}
}

func WithLogLevel(logLevel string) Option {
	return func(options *Options) {
		options.LogLevel = logLevel
	}
}

func WithMaxSize(maxSize int) Option {
	return func(options *Options) {
		options.MaxSize = maxSize
	}
}

func WithMaxAge(maxAge int) Option {
	return func(options *Options) {
		options.MaxSize = maxAge
	}
}

func WithStacktrace(stackTrace string) Option {
	return func(options *Options) {
		options.Stacktrace = stackTrace
	}
}

func WithIsStdout(isStdOut string) Option {
	return func(options *Options) {
		options.IsStdOut = isStdOut
	}
}

func WithProjectName(projectName string) Option {
	return func(options *Options) {
		options.ProjectName = projectName
	}
}
