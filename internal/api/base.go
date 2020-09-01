package api

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
)

func Validate(ctx context.Context, req interface{}) error {

	formatresult, err := govalidator.ValidateStruct(req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), ctx)
		return err
	}
	if !formatresult {
		logx.Error(ReturnMsg[BusParamError], ctx)
		return errors.New(ReturnMsg[BusParamError])
	}
	return nil
}
