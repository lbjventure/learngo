package genModel

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

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
	UserAuthsModel interface {
		Insert(data *UserAuths) (sql.Result, error)
		FindOne(id int64) (*UserAuths, error)
		FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuths, error)
		FindOneByUserIdAuthType(userId int64, authType string) (*UserAuths, error)
		Update(data *UserAuths) error
		Delete(id int64) error
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

func NewUserAuthsModel(conn sqlx.SqlConn, c cache.CacheConf) UserAuthsModel {
	return &defaultUserAuthsModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`user_auths`",
	}
}

func (m *defaultUserAuthsModel) Insert(data *UserAuths) (sql.Result, error) {
	looklookTestUserAuthsIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	looklookTestUserAuthsAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookTestUserAuthsUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?)", m.table, userAuthsRowsExpectAutoSet)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType)
	}, looklookTestUserAuthsIdKey, looklookTestUserAuthsAuthTypeAuthKeyKey, looklookTestUserAuthsUserIdAuthTypeKey)
	return ret, err
}

func (m *defaultUserAuthsModel) FindOne(id int64) (*UserAuths, error) {
	looklookTestUserAuthsIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, id)
	var resp UserAuths
	err := m.QueryRow(&resp, looklookTestUserAuthsIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthsRows, m.table)
		return conn.QueryRow(v, query, id)
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

func (m *defaultUserAuthsModel) FindOneByAuthTypeAuthKey(authType string, authKey string) (*UserAuths, error) {
	looklookTestUserAuthsAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, authType, authKey)
	var resp UserAuths
	err := m.QueryRowIndex(&resp, looklookTestUserAuthsAuthTypeAuthKeyKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `auth_type` = ? and `auth_key` = ? limit 1", userAuthsRows, m.table)
		if err := conn.QueryRow(&resp, query, authType, authKey); err != nil {
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

func (m *defaultUserAuthsModel) FindOneByUserIdAuthType(userId int64, authType string) (*UserAuths, error) {
	looklookTestUserAuthsUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, userId, authType)
	var resp UserAuths
	err := m.QueryRowIndex(&resp, looklookTestUserAuthsUserIdAuthTypeKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `user_id` = ? and `auth_type` = ? limit 1", userAuthsRows, m.table)
		if err := conn.QueryRow(&resp, query, userId, authType); err != nil {
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

func (m *defaultUserAuthsModel) Update(data *UserAuths) error {
	looklookTestUserAuthsIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	looklookTestUserAuthsAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookTestUserAuthsUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, userAuthsRowsWithPlaceHolder)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.UserId, data.AuthKey, data.AuthType, data.Id)
	}, looklookTestUserAuthsIdKey, looklookTestUserAuthsAuthTypeAuthKeyKey, looklookTestUserAuthsUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthsModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	looklookTestUserAuthsIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, id)
	looklookTestUserAuthsAuthTypeAuthKeyKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsAuthTypeAuthKeyPrefix, data.AuthType, data.AuthKey)
	looklookTestUserAuthsUserIdAuthTypeKey := fmt.Sprintf("%s%v:%v", cacheLooklookTestUserAuthsUserIdAuthTypePrefix, data.UserId, data.AuthType)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, looklookTestUserAuthsIdKey, looklookTestUserAuthsAuthTypeAuthKeyKey, looklookTestUserAuthsUserIdAuthTypeKey)
	return err
}

func (m *defaultUserAuthsModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, primary)
}

func (m *defaultUserAuthsModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", userAuthsRows, m.table)
	return conn.QueryRow(v, query, primary)
}
