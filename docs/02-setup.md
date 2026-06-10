# 02 - 环境准备与项目初始化

## 前置要求

- Go 1.21+ 已安装
- 基本的 Go 语法知识（struct、method、interface）
- 熟悉 Git 基本操作
- 一个终端（Windows Terminal、iTerm2、GNOME Terminal 等）

## 步骤 1：创建项目目录

```bash
mkdir terminal_resume
cd terminal_resume
go mod init terminal_resume
```

模块名 `terminal_resume` 将作为后续所有内部包的导入路径前缀。

## 步骤 2：添加依赖

```bash
go get github.com/charmbracelet/wish@latest
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/lipgloss@latest
go get gopkg.in/yaml.v3@latest
```

执行后 `go.mod` 将自动记录依赖版本。

## 步骤 3：规划目录结构

```
terminal_resume/
├── main.go                      # SSH 服务器入口
├── cmd/
│   └── local/
│       └── main.go              # 本地测试入口（无需 SSH）
├── internal/
│   ├── app/
│   │   └── model.go             # Bubble Tea 模型（核心）
│   ├── style/
│   │   └── theme.go             # 样式定义
│   └── data/
│       ├── resume.go            # 数据结构
│       ├── resume.yaml          # 简历内容（可配置）
│       ├── loader.go            # YAML 加载器
│       └── loader_test.go       # 加载测试
├── docs/                        # 本课程文档
├── go.mod
├── go.sum
└── README.md
```

### 目录设计原则

- **`internal/`**：Go 标准约定，表示内部包，不允许外部项目导入
- **`cmd/local/`**：本地测试入口，生产环境使用根目录 `main.go`
- **双入口设计**：
  - 根目录 `main.go` — Wish SSH 服务器（生产）
  - `cmd/local/main.go` — 直接运行 Bubble Tea（开发调试）

## 步骤 4：创建目录

```bash
mkdir -p cmd/local
mkdir -p internal/app
mkdir -p internal/style
mkdir -p internal/data
```

## 步骤 5：验证环境

创建临时文件验证编译环境：

```go
// main.go
package main

import "fmt"

func main() {
    fmt.Println("Environment OK!")
}
```

```bash
go run main.go
# 输出：Environment OK!
```

## 关键技术点：`go:embed`

本项目使用 `go:embed` 将 YAML 配置文件嵌入二进制：

```go
import _ "embed"

//go:embed resume.yaml
var resumeYAML []byte
```

**注意**：
- `//go:embed` 是编译器指令，必须紧跟在变量声明上方
- 文件路径相对于当前 `.go` 文件的位置
- 需要 `_ "embed"` 空白导入以注册 embed 包
- 嵌入后二进制文件独立运行，无需外部 YAML 文件

## 本章小结

- 初始化 Go 模块并安装 Charm 生态依赖
- 采用 `internal/` 标准目录结构
- 双入口设计：SSH 服务器 + 本地调试
- `go:embed` 实现配置内嵌

## 下一步

[03 - 数据层：结构定义与 YAML 配置](03-data-layer.md)
