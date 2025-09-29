package scheduler

import (
	"github.com/robfig/cron/v3"
	"github.com/tinh-tinh/tinhtinh/v2/core"
)

const REGISTRY core.Provide = "SCHEDULER_REGISTRY"

func ForRoot(opts ...cron.Option) core.Modules {
	cron := cron.New(opts...)
	return func(module core.Module) core.Module {
		schedulerModule := module.New(core.NewModuleOptions{})

		schedulerModule.NewProvider(core.ProviderOptions{
			Name:  REGISTRY,
			Value: cron,
		})
		schedulerModule.Export(REGISTRY)

		return schedulerModule
	}
}

func Inject(ref core.RefProvider) *cron.Cron {
	cron, ok := ref.Ref(REGISTRY).(*cron.Cron)
	if !ok {
		return nil
	}

	return cron
}
