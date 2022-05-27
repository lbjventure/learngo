package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
	"looklook/common/globalkey"
	"looklook/common/xerr"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userAuthsFieldNames          = builder.RawFieldNames(&UserAuths{})
	userAuthsRows                = strings.Join(userAuthsFieldNames, ",")
	userAuthsRowsExpectAutoSet   = strings.Join(stringx.Remove(userAuthsFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	userAuthsRowsWithPlaceHolder = strings.Join(stringx.Remove(userAuthsFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheLooklookTestUserAuthsIdPrefix              = "cache:looklookTest:userAuths:id:"
	cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix = "cache:looklookTest:userAuths:authType:authKey:"
	cacheLooklookTestUserAuthsUserIdAuthTypePrefix  = "cache:looklookTest:userAuths:userId:authType:"
)

type (
	userAuthsModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *UserAuths) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*UserAuths, error)
		FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuths, error)
		FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuths, error)
		Update(ctx context.Context, session sqlx.Session, data *UserAuths) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *UserAuths) error
		Delete(ctx context.Context, session sqlx.Session, id int64) error
	}

	defaultUserAuthsModel struct {
		sqlc.CachedConn
		table string
	}

	UserAuths struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		Version    int64     `db:"version"` // 版本号
		UserId     int64     `db:"user_id"`
		AuthKey    string    `db:"auth_key"`  // 平台唯一id
		AuthType   string    `db:"auth_type"` // 平台类型
	}
)

func newUserAuthsModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUserAuthsModel {
	return &defaultUserAuthsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auths`",
	}
}

func (m *defaultUserAuthsModel) Insert(ctx context.Context, session sqlx.Session, data *UserAuths) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, userAuthsRowsExpectAutoSet)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
	}, looklookUsercenterUserAuthUserIdAuthTypeKey, looklookUsercenterUserAuthIdKey, looklookUsercenterUserAuthAuthTypeAuthKeyKey)
}

func (m *defaultUserAuthsModel) FindOne(ctx context.Context, id int64) (*UserAuths, error) {
	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, id)
	var resp UserAuths
	err := m.QueryRowCtx(ctx, &resp, looklookUsercenterUserAuthIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", userAuthsRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id, globalkey.DelStateNo)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthsModel) FindOneByAuthTypeAuthKey(ctx context.Context, authType string, authKey string) (*UserAuths, error) {
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, authType, authKey)
	var resp UserAuths
	err := m.QueryRowIndexCtx(ctx, &resp, looklookUsercenterUserAuthAuthTypeAuthKeyKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_type` = ? and `auth_key` = ? and del_state = ? limit 1", userAuthsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, authType, authKey, globalkey.DelStateNo); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthsModel) FindOneByUserIdAuthType(ctx context.Context, userId int64, authType string) (*UserAuths, error) {
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, userId, authType)
	var resp UserAuths
	err := m.QueryRowIndexCtx(ctx, &resp, looklookUsercenterUserAuthUserIdAuthTypeKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `auth_type` = ? and del_state = ? limit 1", userAuthsRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, userId, authType, globalkey.DelStateNo); err != nil {
			return nil, err
		}
		return resp.Id, nil
	}, m.queryPrimary)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserAuthsModel) Update(ctx context.Context, session sqlx.Session, data *UserAuths) (sql.Result, error) {
	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthsRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id)
	}, looklookUsercenterUserAuthUserIdAuthTypeKey, looklookUsercenterUserAuthIdKey, looklookUsercenterUserAuthAuthTypeAuthKeyKey)
}

func (m *defaultUserAuthsModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *UserAuths) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, userAuthsRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id, oldVersion)
	}, looklookUsercenterUserAuthUserIdAuthTypeKey, looklookUsercenterUserAuthIdKey, looklookUsercenterUserAuthAuthTypeAuthKeyKey)
	if err != nil {
		return err
	}
	updateCount, err := sqlResult.RowsAffected()
	if err != nil {
		return err
	}
	if updateCount == 0 {
		return xerr.NewErrCode(xerr.DB_UPDATE_AFFECTED_ZERO_ERROR)
	}

	return nil
}

func (m *defaultUserAuthsModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	looklookUsercenterUserAuthIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, id)
	looklookUsercenterUserAuthAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookUsercenterUserAuthUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, looklookUsercenterUserAuthIdKey, looklookUsercenterUserAuthAuthTypeAuthKeyKey, looklookUsercenterUserAuthUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, primary)
}
func (m *defaultUserAuthsModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", userAuthsRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary, globalkey.DelStateNo)
}

func (m *defaultUserAuthsModel) tableName() string {
	return m.table
}
