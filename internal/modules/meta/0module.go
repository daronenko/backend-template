package meta

import (
	httpdelivery "github.com/daronenko/backend-template/internal/modules/meta/delivery/http/v1"
	"github.com/daronenko/backend-template/internal/modules/meta/usecase/v1"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("meta",
		httpdelivery.Module(),
		usecase.Module(),
	)
}
