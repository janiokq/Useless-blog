package user

import (
	"context"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/cache"
	pb "github.com/janiokq/Useless-blog/internal/proto/v1/user"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct{}

func (s *Server) UserRegister(ctx context.Context, in *pb.LoginRequset) (out *pb.UserEntity, outerr error) {
	m := new(User)
	out = new(pb.UserEntity)
	err := copier.Copy(m, in)
	if err != nil {
		outerr = status.Error(codes.InvalidArgument, err.Error())
		return
	}

	err = m.Register(ctx)
	if err != nil {
		outerr = status.Error(codes.PermissionDenied, err.Error())
		return
	}
	err = copier.Copy(out, m)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	return
}

func (s *Server) UserInfo(ctx context.Context, in *pb.UserToken) (out *pb.UserEntity, outerr error) {
	m := new(User)
	out = new(pb.UserEntity)
	err := copier.Copy(m, in)
	if err != nil {
		outerr = status.Error(codes.InvalidArgument, err.Error())
		return
	}
	err = m.LoadDataForId(ctx)
	if err != nil {
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	err = copier.Copy(out, m)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	return
}

func (s *Server) UserLogin(ctx context.Context, in *pb.LoginRequset) (out *pb.UserEntity, outerr error) {
	m := new(User)
	out = new(pb.UserEntity)
	err := copier.Copy(m, in)
	if err != nil {
		outerr = status.Error(codes.PermissionDenied, err.Error())
		return
	}
	err = m.LoginFormPhoneAndPassword(ctx, in.Phone, in.Password)
	if err != nil {
		outerr = status.Error(codes.PermissionDenied, err.Error())
		return
	}
	err = copier.Copy(out, m)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}

	return
}

func (s *Server) UserLogout(ctx context.Context, in *pb.UserToken) (out *pb.UserToken, outerr error) {
	out = new(pb.UserToken)
	err := copier.Copy(out, in)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	key := cache.GetIdKey(cinit.TokenRedisCachePrefix, in.Id)
	_, err = cache.CacheGetBuyKey(ctx, key)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.OutOfRange, err.Error())
		return
	}
	cache.CacheDel(ctx, cinit.TokenRedisCachePrefix, in.Id)
	return
}

func (s *Server) UserUpdateInfo(ctx context.Context, in *pb.UserEntity) (out *pb.UserEntity, outerr error) {
	out = new(pb.UserEntity)
	u := new(User)
	err := copier.Copy(u, in)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	err = u.update(ctx)
	if err != nil {
		outerr = status.Error(codes.InvalidArgument, err.Error())
		return
	}
	err = copier.Copy(out, u)
	if err != nil {
		logx.Error(err.Error(), ctx)
		outerr = status.Error(codes.Internal, err.Error())
		return
	}
	return
}
