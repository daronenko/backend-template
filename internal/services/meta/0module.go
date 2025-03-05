package meta

import (
	httpdelivery "github.com/daronenko/backend-template/internal/services/meta/delivery/http/v1"
	"github.com/daronenko/backend-template/internal/services/meta/usecase/v1"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("meta",
		httpdelivery.Module(),
		usecase.Module(),
	)
}
