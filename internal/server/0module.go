package server

import (
	"github.com/daronenko/backend-template/internal/server/httpserver"
	"github.com/daronenko/backend-template/internal/server/svr"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("server",
		fx.Provide(httpserver.New),
		fx.Provide(svr.CreateEndpointGroups))
}
