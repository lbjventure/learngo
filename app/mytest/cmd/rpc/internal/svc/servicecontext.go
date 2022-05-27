package svc

import (
	"looklook/app/mytest/cmd/rpc/internal/config"
	"looklook/app/mytest/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	)

type ServiceContext struct {
	Config config.Config
	RedisClient *redis.Redis

	UserModel     model.UsersModel
	UserAuthModel model.UserAuthsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),

		UserAuthModel: model.NewUserAuthsModel(sqlConn, c.Cache),
		UserModel:     model.NewUsersModel(sqlConn, c.Cache),
	}
}
