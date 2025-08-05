package di

import (
	"github.com/watchs/application/interfaces"
	"github.com/watchs/application/services"
	"github.com/watchs/domain/repository"
	"github.com/watchs/infrastructure/persistence"
)

// Container 依赖注入容器
type Container struct {
	configRepo              repository.ConfigRepository
	configApplicationService interfaces.ConfigApplicationService
	watchApplicationService  interfaces.WatchApplicationService
}

// NewContainer 创建新的依赖注入容器
func NewContainer() *Container {
	container := &Container{}
	container.initializeDependencies()
	return container
}

// initializeDependencies 初始化所有依赖关系
func (c *Container) initializeDependencies() {
	// 基础设施层
	c.configRepo = persistence.NewJsonConfigRepository()

	// 应用服务层
	c.configApplicationService = services.NewConfigApplicationService(c.configRepo)
	c.watchApplicationService = services.NewWatchApplicationService(c.configApplicationService)
}

// GetConfigRepository 获取配置仓储
func (c *Container) GetConfigRepository() repository.ConfigRepository {
	return c.configRepo
}

// GetConfigApplicationService 获取配置应用服务
func (c *Container) GetConfigApplicationService() interfaces.ConfigApplicationService {
	return c.configApplicationService
}

// GetWatchApplicationService 获取监控应用服务
func (c *Container) GetWatchApplicationService() interfaces.WatchApplicationService {
	return c.watchApplicationService
}
