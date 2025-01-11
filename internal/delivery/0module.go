package delivery

import (
	"go.uber.org/fx"

	metaDelivery "github.com/daronenko/backend-template/internal/delivery/meta"
)

type opt int

const (
	OptIncludeSwagger opt = iota
)

func Module(o ...opt) fx.Option {
	opts := []fx.Option{
		// Delivery (meta)
		metaDelivery.Module(),
	}

	for _, opt := range o {
		switch opt {
		case OptIncludeSwagger:
			// opts = append(opts, fx.Invoke(metaDelivery.RegisterSwagger))
		}
	}

	return fx.Module("delivery",
		// options
		opts...,
	)
}
