package model

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
 "context"
	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"

	"looklook/common/globalkey"
	"looklook/common/xerr"
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
	usersModel interface {
		Insert(ctx context.Context, session sqlx.Session, data *Users) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*Users, error)
		FindOneByMobile(ctx context.Context, mobile string) (*Users, error)
		Update(ctx context.Context, session sqlx.Session, data *Users) (sql.Result, error)
		UpdateWithVersion(ctx context.Context, session sqlx.Session, data *Users) error
		Delete(ctx context.Context, session sqlx.Session, id int64) error
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

func newUsersModel(conn sqlx.SqlConn, c cache.CacheConf) *defaultUsersModel {
	return &defaultUsersModel{
		CachedConn: sqlc.NewConn(conn, c),
		table:      "`users`",
	}
}


func (m *defaultUsersModel) Insert(ctx context.Context, session sqlx.Session, data *Users) (sql.Result, error) {
	data.DeleteTime = time.Unix(0, 0)
	fmt.Println(1,data.DeleteTime)
	fmt.Printf("2 %+v",m.table)
	looklookUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, data.Id)
	looklookUsercenterUserMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table, usersRowsExpectAutoSet)
		fmt.Println(3,query)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info)
	}, looklookUsercenterUserIdKey, looklookUsercenterUserMobileKey)
}

func (m *defaultUsersModel) FindOne(ctx context.Context, id int64) (*Users, error) {
	looklookUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, id)
	var resp Users
	err := m.QueryRowCtx(ctx, &resp, looklookUsercenterUserIdKey, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", usersRows, m.table)
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

func (m *defaultUsersModel) FindOneByMobile(ctx context.Context, mobile string) (*Users, error) {
	looklookUsercenterUserMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, mobile)
	var resp Users
	err := m.QueryRowIndexCtx(ctx, &resp, looklookUsercenterUserMobileKey, m.formatPrimary, func(ctx context.Context, conn sqlx.SqlConn, v interface{}) (i interface{}, e error) {
		query := fmt.Sprintf("select %s from %s where `mobile` = ? and del_state = ? limit 1", usersRows, m.table)
		if err := conn.QueryRowCtx(ctx, &resp, query, mobile, globalkey.DelStateNo); err != nil {
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

func (m *defaultUsersModel) Update(ctx context.Context, session sqlx.Session, data *Users) (sql.Result, error) {
	looklookUsercenterUserMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	looklookUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUserAuthsIdPrefix, data.Id)
	return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, usersRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id)
	}, looklookUsercenterUserIdKey, looklookUsercenterUserMobileKey)
}

func (m *defaultUsersModel) UpdateWithVersion(ctx context.Context, session sqlx.Session, data *Users) error {

	oldVersion := data.Version
	data.Version += 1

	var sqlResult sql.Result
	var err error

	looklookUsercenterUserMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	looklookUsercenterUserIdKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, data.Id)
	sqlResult, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where `id` = ? and version = ? ", m.table, usersRowsWithPlaceHolder)
		if session != nil {
			return session.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id, oldVersion)
		}
		return conn.ExecCtx(ctx, query, data.DeleteTime, data.DelState, data.Version, data.Mobile, data.Password, data.Nickname, data.Sex, data.Avatar, data.Info, data.Id, oldVersion)
	}, looklookUsercenterUserIdKey, looklookUsercenterUserMobileKey)
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

func (m *defaultUsersModel) Delete(ctx context.Context, session sqlx.Session, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil {
		return err
	}

	looklookUsercenterUserIdKey := fmt.Sprintf("%s%v",  cacheLooklookTestUsersIdPrefix,data.Id)
	looklookUsercenterUserMobileKey := fmt.Sprintf("%s%v", cacheLooklookTestUsersMobilePrefix, data.Mobile)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		if session != nil {
			return session.ExecCtx(ctx, query, id)
		}
		return conn.ExecCtx(ctx, query, id)
	}, looklookUsercenterUserIdKey, looklookUsercenterUserMobileKey)
	return err
}

func (m *defaultUsersModel) formatPrimary(primary interface{}) string {
	return fmt.Sprintf("%s%v", cacheLooklookTestUsersIdPrefix, primary)
}
func (m *defaultUsersModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary interface{}) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? and del_state = ? limit 1", usersRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary, globalkey.DelStateNo)
}

func (m *defaultUsersModel) tableName() string {
	return m.table
}
