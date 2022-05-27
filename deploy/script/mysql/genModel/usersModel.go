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
	usersFieldNames          = builder.RawFieldNames(&Users{})
	usersRows                = strings.Join(usersFieldNames, ",")
	usersRowsExpectAutoSet   = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_time`", "`update_time`"), ",")
	usersRowsWithPlaceHolder = strings.Join(stringx.Remove(usersFieldNames, "`id`", "`create_time`", "`update_time`"), "=?,") + "=?"

	cacheLooklookTestUsersIdPrefix     = "cache:looklookTest:users:id:"
	cacheLooklookTestUsersMobilePrefix = "cache:looklookTest:users:mobile:"
)

type (
	UsersModel interface {
		Insert(data *Users) (sql.Result, error)
		FindOne(id int64) (*Users, error)
		FindOneByMobile(mobile string) (*Users, error)
		Update(data *Users) error
		Delete(id int64) error
	}

	defaultUsersModel struct {
		sqlc.CachedConn
		table string
	}

	Users struct {
		Id         int64     `db:"id"`
		CreateTime time.Time `db:"create_time"`
		UpdateTime time.Time `db:"update_time"`
		DeleteTime time.Time `db:"delete_time"`
		DelState   int64     `db:"del_state"`
		Version    int64     `db:"version"` // 版本号
		Mobile     string    `db:"mobile"`
		Password   string    `db:"password"`
		Nickname   string    `db:"nickname"`
		Sex        int64     `db:"sex"` // 性别 0:男 1:女
		Avatar     string    `db:"avatar"`
		Info       string    `db:"info"`
	}
)

func NewUsersModel(conn sqlx.SqlConn, c cache.CacheConf) UsersModel {
	return &defaultUsersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`users`",
	}
}

func (m *defaultUsersModel) Insert(data *Users) (sql.Result, error) {
	looklookTestUsersIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, data.Id)
	looklookTestUsersMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	ret, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info)
	}, looklookTestUsersIdKey, looklookTestUsersMobileKey)
	return ret, err
}

func (m *defaultUsersModel) FindOne(id int64) (*Users, error) {
	looklookTestUsersIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, id)
	var resp Users
	err := m.QueryRow(&resp, looklookTestUsersIdKey, func(conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
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

func (m *defaultUsersModel) FindOneByMobile(mobile string) (*Users, error) {
	looklookTestUsersMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, mobile)
	var resp Users
	err := m.QueryRowIndex(&resp, looklookTestUsersMobileKey, m.formatPrimary, func(conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `mobile` = ? limit 1", usersRows, m.table)
		if err := conn.QueryRow(&resp, query, mobile); err != nil {
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

func (m *defaultUsersModel) Update(data *Users) error {
	looklookTestUsersIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, data.Id)
	looklookTestUsersMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	_, err := m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
		return conn.Exec(query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id)
	}, looklookTestUsersIdKey, looklookTestUsersMobileKey)
	return err
}

func (m *defaultUsersModel) Delete(id int64) error {
	data, err := m.FindOne(id)
	if err != nil {
		return err
	}

	looklookTestUsersIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, id)
	looklookTestUsersMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	_, err = m.Exec(func(conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.Exec(query, id)
	}, looklookTestUsersMobileKey, looklookTestUsersIdKey)
	return err
}

func (m *defaultUsersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, primary)
}

func (m *defaultUsersModel) queryPrimary(conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", usersRows, m.table)
	return conn.QueryRow(v, query, primary)
}
