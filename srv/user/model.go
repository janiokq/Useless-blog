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
	Id        int64  `json:"id" db:"id"`
	CreateAt  string `json:"create_at" db:"create_at"`
	UpdateAt  string `json:"update_at" db:"update_at"`
	UserName  string `json:"user_name" db:"user_name"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
	Token     string `json:"token" db:"token" valid:"required~token必须存在"`
	Phone     string `json:"phone" db:"phone" valid:"required~phone必须存在"`
	Password  string `json:"password" db:"password"`
}

// 验证手机号格式
func (u *User) validatePhone() error {
	result, _ := regexp.MatchString(utils.PhoneChinaRegular, u.Phone)
	if !result {
		return errors.New(PhoneNumberFormat)
	}
	return nil
}

func (u *User) validatePassword() error {
	if len(u.Password) <= 5 {
		return errors.New(MinimumPasswordLength)
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
	//设置UserEntity 缓存到redis
	cache.CacheSetHM(ctx, cinit.UserEntityRedisPrefix, u.Id, utils.Struct2Map(*u), cinit.UserEntityExpirationtime*3600)
	go msgNotify(ctx, "修改用户:"+strconv.FormatInt(u.Id, 10))
	return nil
}
func (u *User) FormMap(data map[string]string) {
	//u.Id = strconv.FormatInt()
}

func (u *User) LoadDataForId(ctx context.Context) error {
	err := u.validateId()
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	_u, err := cache.CacheGetHM(ctx, cinit.UserEntityRedisPrefix, u.Id)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	utils.Map2Struct(_u, u)

	return nil
}

func (u *User) Register(ctx context.Context) error {
	err := utils.V(u.validatePhone, u.validatePassword)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	count := 0
	cinit.Mysql.QueryRow("SELECT count(1) FROM user WHERE phone=?", u.Phone).Scan(&count)
	if count > 0 {
		logx.Info(PhoneNumberExist, ctx)
		return errors.New(PhoneNumberExist)
	}
	insertTime := time.Now().Format(cinit.TimeFormatting)
	u.CreateAt = insertTime
	u.UpdateAt = insertTime

	insertResult, err := cinit.Mysql.Exec("INSERT INTO user (phone,password,create_at,update_at) values(?,?,?,?) ", u.Phone, utils.MD5(u.Password), insertTime, insertTime)
	if err != nil {
		logx.Error(err.Error(), ctx)
		return err
	}
	id, err := insertResult.LastInsertId()
	u.Id = id
	u.UserName = ""
	j := new(jwt.Msg)
	j.UserName = u.UserName
	j.Id = u.Id
	token, err := jwt.Encode(*j)
	if err != nil {
		logx.Error(err.Error(), ctx)
	} else {
		u.Token = token
		//设置token 缓存得到redis
		cache.CacheSet(ctx, cinit.TokenRedisCachePrefix, u.Id, token, cinit.TokenExpirationtime*3600)
		//设置UserEntity 缓存到redis

		cache.CacheSetHM(ctx, cinit.UserEntityRedisPrefix, u.Id, utils.Struct2Map(*u), cinit.UserEntityExpirationtime*3600)
	}
	return nil
}

func (u *User) LoginFormPhoneAndPassword(ctx context.Context, phone string, password string) error {
	err := utils.V(u.validatePhone, u.validatePassword)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}

	row := cinit.Mysql.QueryRow(`SELECT id,UNIX_TIMESTAMP(create_at),UNIX_TIMESTAMP(update_at),user_name,avatar_url,phone,password FROM user WHERE phone=?`, phone)
	err = row.Scan(&u.Id, &u.CreateAt, &u.UpdateAt, &u.UserName, &u.AvatarUrl, &u.Phone, &u.Password)
	switch {
	case err == sql.ErrNoRows:
		logx.Info(err.Error(), ctx)
		return errors.New(PhoneDoesNotExist)
	case err != nil:
		logx.Error(err.Error(), ctx)
		return err
	default:
	}

	t, err := strconv.ParseInt(u.CreateAt, 10, 64)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	cs := time.Unix(t, 0).Format(cinit.TimeFormatting)
	u.CreateAt = cs

	ut, err := strconv.ParseInt(u.UpdateAt, 10, 64)
	if err != nil {
		logx.Info(err.Error(), ctx)
		return err
	}
	us := time.Unix(ut, 0).Format(cinit.TimeFormatting)
	u.UpdateAt = us

	if !strings.EqualFold(utils.MD5(password), u.Password) {
		return errors.New(PasswordMistake)
	}
	j := new(jwt.Msg)
	j.UserName = u.UserName
	j.Id = u.Id
	token, err := jwt.Encode(*j)
	if err != nil {
		logx.Error(err.Error(), ctx)
	} else {
		u.Token = token
		//设置token 缓存得到redis
		cache.CacheSet(ctx, cinit.TokenRedisCachePrefix, u.Id, token, cinit.TokenExpirationtime*3600)
		//设置UserEntity 缓存到redis
		cache.CacheSetHM(ctx, cinit.UserEntityRedisPrefix, u.Id, utils.Struct2Map(*u), cinit.UserEntityExpirationtime*3600)
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
