package cinit

const (
	//message queue
	TopicUserSrvChange = "user_msg_change"

	//redis
	TokenRedisCachePrefix = "user_token"
	TokenExpirationtime   = 72

	UserEntityRedisPrefix    = "user_entity"
	UserEntityExpirationtime = 72

	//other
	TimeFormatting = "2006-01-02 15:04:05"

	ReqParam        = "req_param"     //  请求参数绑定
	JWTName         = "Authorization" //  JWT请求头名称
	JWTMsg          = "JWT-MSG"       //  JWT自定义的消息
	FloatComputeBit = 2               //  浮点计算位数

)
