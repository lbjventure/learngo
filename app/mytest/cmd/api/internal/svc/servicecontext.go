package svc


import (
	"looklook/app/mytest/cmd/api/internal/config"
	"looklook/app/mytest/cmd/rpc/mytest"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config
	MytestRpc mytest.Mytest

	SetUidToCtxMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		MytestRpc: mytest.NewMytest(zrpc.MustNewClient(c.MytestRpcConf)),
	}
}
