package types

import "runtime/debug"

type Service interface {
	Run() error
}

type ServiceManager struct {
	ctx      Context
	services []Service
}

type sync func() error

func NewServiceManager(ctx Context, services ...Service) *ServiceManager {
	return &ServiceManager{
		ctx:      ctx,
		services: services,
	}
}

func (m *ServiceManager) Run() {
	for _, service := range m.services {
		go m.run(service.Run)
	}
}

func (m *ServiceManager) run(s sync) {
	defer func() {
		if v := recover(); v != nil {
			m.ctx.Logger().Error("detected panic recover: ", v)
			debug.PrintStack()
		}
	}()

	if err := s(); err != nil {
		m.ctx.Logger().Error(err)
	}
}
