package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/api"
	pb "github.com/janiokq/Useless-blog/internal/proto/v1/user"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/jinzhu/copier"
)

func UserInfo(c *gin.Context) {
	_req := &pb.UserToken{}
	val, _ := c.Get(cinit.JWTMsg)
	err := copier.Copy(_req, val)
	_req.Token = c.Request.Header.Get(cinit.JWTName)

	if err != nil {
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())
		return
	}
	err = api.Validate(c.Request.Context(), _req)
	if err != nil {
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())
		return
	}
	_rsp, err := UserClient.UserInfo(c.Request.Context(), _req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.RPCErr(c, err)
		return
	}
	api.HandleSuccess(c, _rsp)
}

func Register(c *gin.Context) {
	var _req = new(pb.LoginRequset)
	err := c.Bind(&_req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	err = api.Validate(c.Request.Context(), _req)
	if err != nil {
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	_rsp, err := UserClient.UserRegister(c.Request.Context(), _req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.RPCErr(c, err)

		return
	}
	api.HandleSuccess(c, _rsp)
}

func Login(c *gin.Context) {
	_req := new(pb.LoginRequset)
	err := c.Bind(&_req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	err = api.Validate(c.Request.Context(), _req)
	if err != nil {
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	_rsp, err := UserClient.UserLogin(c.Request.Context(), _req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.RPCErr(c, err)

		return
	}
	api.HandleSuccess(c, _rsp)

}

func UserUpdate(c *gin.Context) {
	_req := new(pb.UserEntity)
	err := c.Bind(&_req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	err = api.Validate(c.Request.Context(), _req)
	if err != nil {
		logx.Error(err.Error(), c.Request.Context())
		api.HandleError(c, api.BusParamConvertError, err.Error())

		return
	}
	_rsp, err := UserClient.UserUpdateInfo(c.Request.Context(), _req)
	if err != nil {
		//  解析返回的错误信息
		logx.Error(err.Error(), c.Request.Context())
		api.RPCErr(c, err)

		return
	}
	api.HandleSuccess(c, _rsp)
}
