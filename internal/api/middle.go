package api

import (
	"github.com/gin-gonic/gin"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/cache"
	"github.com/janiokq/Useless-blog/internal/jwt"
	metricsInternal "github.com/janiokq/Useless-blog/internal/metrics"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/prometheus/common/log"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()
		// 从请求头获取token信息
		jwtString := c.Request.Header.Get(cinit.JWTName)
		if jwtString == " " {
			HandleError(c, ReqTokenError, "token验证失败")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		jwtmsg, err := jwt.Decode(strings.Trim(jwtString, " "))
		if err != nil {
			log.Info(err.Error(), ctx)
			HandleError(c, ReqTokenError, "token验证失败")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		key, err := cache.CacheGetBuyKey(ctx, cache.GetIdKey(cinit.TokenRedisCachePrefix, jwtmsg.Id))
		if err != nil || key == "" {
			log.Info(err.Error(), ctx)
			HandleError(c, ReqTokenError, "token验证失败")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(cinit.JWTMsg, jwtmsg)
		logx.Infof("header:%+v", c.Request.Header)
		c.Next()

	}
}

func TraceHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		logx.Infof("header:%+v", c.Request.Header)
		c.Next()
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func MetricsFunc(m *metricsInternal.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		stop := time.Now()
		latency := stop.Sub(start)
		status := c.Writer.Status()
		// Total request count.
		meter := metrics.GetOrRegisterMeter("status."+strconv.Itoa(status), m.GetRegistry())
		meter.Mark(c.Request.ContentLength)

		// Request size in bytes.
		// meter = metrics.GetOrRegisterMeter(m.WithPrefix("status."+strconv.Itoa(status)), m.GetRegistry())
		// meter.Mark(req.ContentLength)

		// Request duration in nanoseconds.
		h := metrics.GetOrRegisterHistogram("h_status."+strconv.Itoa(status), m.GetRegistry(),
			metrics.NewExpDecaySample(1028, 0.015))
		h.Update(latency.Nanoseconds())

	}
}
