# Watchs - 文件变更监控工具

一个基于DDD（领域驱动设计）架构的文件变更监控工具，可以监控指定目录下的文件变化，并在文件变化时执行指定命令。

## 功能特点

* 监控指定目录下的文件变化（递归）
* 可以按文件类型过滤（支持多类型）
* 可以排除特定目录或文件
* 文件变化时执行指定命令
* 支持配置文件和命令行参数
* 支持通过命令行生成配置文件
* 支持交互式配置向导
* 基于DDD架构，代码结构清晰，易于维护和扩展
* 使用命令模式实现可扩展的命令行界面
* 集成GitHub Actions自动化构建和发布

## 最新版本

当前最新版本为 [v1.0.0](https://github.com/fly32101/watchs/releases/tag/v1.0.0)，已通过GitHub Actions自动构建并发布。

## 项目架构

项目采用DDD（领域驱动设计）架构，分为以下几层：

* **领域层（Domain）**：包含核心业务逻辑和实体  
   * `entity`：领域实体，如配置和文件事件  
   * `service`：领域服务接口  
   * `repository`：仓储接口
* **应用层（Application）**：协调领域对象完成用户任务  
   * 应用服务，如文件监控服务
* **基础设施层（Infrastructure）**：提供技术实现  
   * `persistence`：配置持久化实现  
   * `watcher`：文件监控和命令执行实现
* **表示层（Presentation）**：处理用户交互  
   * `cli`：命令行界面，使用命令模式实现

### 设计模式

项目使用了以下设计模式：

* **命令模式（Command Pattern）**：将命令行操作封装为对象，实现命令的可扩展性和可组合性
* **依赖注入（Dependency Injection）**：通过构造函数注入依赖，降低组件间耦合
* **仓储模式（Repository Pattern）**：抽象数据访问逻辑，实现持久化与领域逻辑的分离
* **工厂方法（Factory Method）**：创建复杂对象，封装对象创建逻辑

## 自动化构建与发布

项目使用GitHub Actions进行自动化构建和发布：

* **持续集成（CI）**：在每次代码推送和PR时，自动运行测试和构建检查
* **自动发布**：仅在创建新的标签（如v1.0.0）时，自动构建多平台二进制文件并创建GitHub Release
* **手动构建**：可通过GitHub Actions手动触发构建，用于开发测试

### 发布新版本

要发布新版本，只需创建并推送一个新的标签：

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

GitHub Actions将自动使用GoReleaser构建多平台二进制文件并创建Release。

### 手动构建

如需手动构建开发版本，可以：
1. 访问GitHub仓库的Actions页面
2. 选择"Build"工作流
3. 点击"Run workflow"手动触发构建

## 安装

### 从GitHub Releases安装

访问[GitHub Releases页面](https://github.com/fly32101/watchs/releases)，下载适合你系统的二进制文件：
- Linux: `watchs_Linux_x86_64.tar.gz` (Intel/AMD) 或 `watchs_Linux_arm64.tar.gz` (ARM)
- Windows: `watchs_Windows_x86_64.zip`
- macOS: `watchs_Darwin_x86_64.tar.gz` (Intel) 或 `watchs_Darwin_arm64.tar.gz` (Apple Silicon)

### 从源码安装

```bash
go install github.com/fly32101/watchs/cmd/watchs@latest
```

或者从源码编译：

```bash
git clone https://github.com/fly32101/watchs.git
cd watchs
go build -o watchs ./cmd/watchs
```

## 使用方法

### 查看帮助信息

```bash
watchs help
```

或查看特定命令的帮助信息：

```bash
watchs help <命令名称>
watchs <命令名称> --help
```

### 查看版本信息

```bash
watchs version
```

### 交互式配置

使用交互式向导创建配置文件（推荐新用户使用）：

```bash
watchs interactive
```

向导将引导你完成所有配置选项，并可以选择立即启动监控。

### 命令行生成配置文件

使用 `init` 命令生成配置文件：

```bash
watchs init -config watchs.json -dir ./ -types .go,.js -exclude vendor,node_modules -cmd "go run main.go"
```

参数说明：

* `-config`: 配置文件路径（默认为 `watchs.json`）
* `-dir`: 要监控的目录（默认为 `./`）
* `-types`: 要监控的文件类型，以逗号分隔
* `-exclude`: 要排除的路径，以逗号分隔
* `-cmd`: 文件变化时执行的命令（默认为 `echo 文件已更新`）
* `-force`: 是否强制覆盖已存在的配置文件

### 使用配置文件

创建 `watchs.json` 配置文件后，直接运行：

```bash
watchs
```

或者指定配置文件路径：

```bash
watchs -config custom-watchs.json
```

也可以使用 watch 命令（与直接运行相同）：

```bash
watchs watch -config watchs.json
```

### 使用命令行参数

也可以直接通过命令行参数运行，无需配置文件：

```bash
watchs -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

或者使用 watch 命令：

```bash
watchs watch -dir ./ -types .go,.json -exclude vendor,node_modules,.git -cmd "go run main.go"
```

## 命令行参数

### 监控命令参数 (watch)

* `-config`: 配置文件路径（默认为 `watchs.json`）
* `-dir`: 要监控的目录（覆盖配置文件）
* `-types`: 要监控的文件类型，以逗号分隔（覆盖配置文件）
* `-exclude`: 要排除的路径，以逗号分隔（覆盖配置文件）
* `-cmd`: 文件变化时执行的命令（覆盖配置文件）
* `-debounce`: 防抖时间，单位毫秒（默认为500）

### 初始化命令参数 (init)

* `-config`: 配置文件路径（默认为 `watchs.json`）
* `-dir`: 要监控的目录（默认为 `./`）
* `-types`: 要监控的文件类型，以逗号分隔
* `-exclude`: 要排除的路径，以逗号分隔
* `-cmd`: 文件变化时执行的命令（默认为 `echo 文件已更新`）
* `-force`: 是否强制覆盖已存在的配置文件

## 示例

### 查看帮助

```bash
# 显示所有可用命令
watchs help

# 显示特定命令的帮助信息
watchs help interactive
watchs init --help
```

### 交互式配置

```bash
# 启动交互式配置向导
watchs interactive
```

### 生成配置文件

```bash
# 生成默认配置文件
watchs init

# 生成自定义配置文件
watchs init -config frontend.json -dir ./frontend -types .js,.jsx,.ts,.tsx,.css -exclude node_modules -cmd "npm run build"
```

### 监控文件变化

```bash
# 使用配置文件监控
watchs

# 监控当前目录下的所有 .go 文件，排除 vendor 目录，当文件变化时运行测试
watchs -dir ./ -types .go -exclude vendor -cmd "go test ./..."

# 监控前端项目并自动重新构建
watchs -dir ./frontend -types .js,.jsx,.ts,.tsx,.css -exclude node_modules -cmd "npm run build"
```

## 扩展命令

如果你想添加新的命令，只需实现 `Command` 接口并在 `CLI` 初始化时注册即可：

```go
// 实现命令接口
type MyCommand struct {
    // 依赖项
}

func (c *MyCommand) Name() string {
    return "mycommand"
}

func (c *MyCommand) Description() string {
    return "我的自定义命令"
}

func (c *MyCommand) Execute(args []string) error {
    // 命令实现
    return nil
}

// 在 CLI 初始化时注册
registry.Register(NewMyCommand(...))
```

## 注意事项

* 命令会在监控目录下执行
* 如果命令是长时间运行的进程，当文件再次变化时，之前的进程会被终止并重新启动
* 使用防抖机制避免频繁触发命令执行

## 开源协议

本项目采用 [MIT 许可证](LICENSE)。

## 贡献

欢迎提交 Issue 和 Pull Request。

## 相关链接

- [GitHub 仓库](https://github.com/fly32101/watchs)
- [发布页面](https://github.com/fly32101/watchs/releases) 