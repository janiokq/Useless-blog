#define variables
GOPROXY=https://goproxy.io
GO111MODULE=on

.PHONY : run
run :
	@echo "测试22"
.PHONY : proto
proto :
	@echo "生成proto开始"
	@chmod +x ./scripts/proto.sh && ./scripts/proto.sh
	@echo "生成proto结束"

