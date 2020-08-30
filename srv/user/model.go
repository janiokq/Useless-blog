package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/cache"
	"github.com/janiokq/Useless-blog/internal/jwt"
	"github.com/janiokq/Useless-blog/internal/utils"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id        int64     `json:"id" db:"id"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
	UpdateAt  time.Time `json:"update_at" db:"update_at"`
	UserName  string    `json:"user_name" db:"name"`
	AvatarUrl string    `json:"avatar_url" db:"avatar_url"`
	Token     string    `json:"token" db:"token" valid:"required~token必须存在"`
	Phone     string    `json:"phone" db:"phone" valid:"required~phone必须存在"`
	Password  string    `json:"password" db:"password"`
}

const cachePrefix = "User_Prefix"

// 验证手机号格式
func (u *User) validatePhone() error {
	//m.Iphone
	result, _ := regexp.MatchString(utils.PhoneChinaRegular, u.Phone)
	if !result {
		return errors.New(PhoneNumberFormat)
	}
	return nil
}
func (u *User) validateId() error {
	if u.Id <= 0 {
		return errors.New(IdMinimum)
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
	go msgNotify(ctx, "修改用户:"+strconv.FormatInt(u.Id, 10))
	return nil
}

func (u *User) LoginFormPhoneAndPassword(ctx context.Context, phone string, password string) error {
	err := u.validatePhone()
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	err = cinit.Mysql.Get(u, `SELECT id,create_at,update_at,user_name,avatar_url,phone,password FROM user WHERE phone=?`, phone)
	switch {
	case err == sql.ErrNoRows:
		logx.Info(err.Error(), ctx)
		return errors.New(PhoneDoesNotExist)
	case err != nil:
		logx.Error(err.Error(), ctx)
		return err
	default:
	}
	if !strings.EqualFold(utils.MD5(password), u.Password) {
		return errors.New(PasswordMistake)
	}
	j := new(jwt.Msg)
	j.UserName = u.UserName
	j.UserID = u.Id
	token, err := jwt.Encode(*j)
	if err != nil {
		logx.Error(err.Error(), ctx)
	} else {
		u.Token = token
		cache.CacheSet(ctx, cinit.TokenRedisCachePrefix, u.Id, token, cinit.TokenExpirationtime*3600)
	}
	return nil
}

func (u *User) update(ctx context.Context) error {
	err := u.beforUpdate(ctx)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	r, err := cinit.Mysql.Exec("UPDATE user set create_at=?,update_at=?,user_name=?,avatar_url=?,phone=?", u.CreateAt, u.UpdateAt, u.UserName, u.AvatarUrl, u.Phone)
	err = utils.R(r, err)
	if err != nil {
		logx.Error(err.Error(), ctx)
		return err
	}
	return u.afterUpdate(ctx)
}
