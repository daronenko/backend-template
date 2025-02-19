package ctx

const (
	EnvServer Env = iota
	EnvWorker
)

type Env int

type Ctx struct {
	Env Env
}

func Declare(env Env) Ctx {
	return Ctx{env}
}
