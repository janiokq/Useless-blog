package user

import (
	"context"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/cache"
	"github.com/janiokq/Useless-blog/internal/utils"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"time"
)

type User struct {
	Id        int64     `json:"id" db:"id"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
	UpdateAt  time.Time `json:"update_at" db:"update_at"`
	Name      int64     `json:"name" db:"name"`
	AvatarUrl string    `json:"avatar_url" db:"avatar_url"`
	Token     string    `json:"token" db:"token" valid:"required~token必须存在"`
	Phone     string    `json:"phone" db:"phone" valid:"required~phone必须存在"`
}

const cachePrefix = "User_Prefix"

func (u *User) validateId() error {
	if u.Id <= 0 {
		return errors.New("id必须大于0")
	}
	return nil
}
func (u *User) validate() error {
	_, err := govalidator.ValidateStruct(u)
	return err
}

func (u *User) beforUpdate(ctx context.Context) error {
	err := utils.V(u.validate, u.validateId)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) afterUpdate(ctx context.Context) error {
	go cache.CacheDel(ctx, cachePrefix, u.Id)
	//go msgNotify(ctx, "修改用户:"+m.UserName)

	return nil
}

func (u *User) update(ctx context.Context) error {
	err := u.beforUpdate(ctx)
	if err != nil {
		//log.Info(err.Error(), ctx)
		//TODO 日志记录
		return err
	}
	r, err := cinit.Mysql.Exec("UPDATE user set create_at=?,update_at=?,name=?,avatar_url=?,phone=?", u.CreateAt, u.UpdateAt, u.Name, u.AvatarUrl, u.Phone)
	err = utils.R(r, err)
	if err != nil {
		logx.Error(err.Error(), ctx)
		return err
	}

	return u.afterUpdate(ctx)
}
