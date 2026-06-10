# 09 - 本地测试与调试

## 为什么需要本地入口？

每次测试都通过 SSH 连接效率太低。本地入口直接运行 Bubble Tea，无需 SSH 服务器：

```
┌──────────────┐              ┌──────────────┐
│  本地调试     │   对比       │  SSH 模式    │
├──────────────┤              ├──────────────┤
│ go run       │              │ 生成 Host Key │
│ ./cmd/local  │  更简单      │ 启动服务器   │
│              │  ───────►    │ ssh 连接     │
│ 直接看到界面  │              │              │
└──────────────┘              └──────────────┘
```

## 实现本地入口

创建 `cmd/local/main.go`：

```go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "terminal_resume/internal/app"
)

func main() {
    m := app.NewModel()
    p := tea.NewProgram(m, tea.WithAltScreen())

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

**与 SSH 入口的区别**：

| 方面 | 本地入口 | SSH 入口 |
|------|----------|----------|
| 代码 | `tea.NewProgram(m, tea.WithAltScreen())` | `bubbletea.Middleware(teaHandler)` |
| Host Key | 不需要 | 必须生成 |
| 网络 | 无 | 监听 TCP 端口 |
| 适用场景 | 开发调试 | 生产部署 |

## 运行本地版本

```bash
go run ./cmd/local
```

直接看到 TUI 界面，按 `q` 退出。

## 单元测试

### 测试 YAML 加载

已创建的 `internal/data/loader_test.go`：

```go
package data

import "testing"

func TestLoadResume(t *testing.T) {
    resume, err := LoadResume()
    if err != nil {
        t.Fatalf("LoadResume failed: %v", err)
    }

    if resume.Name == "" {
        t.Error("expected Name to be non-empty")
    }
    if resume.Title == "" {
        t.Error("expected Title to be non-empty")
    }
    if len(resume.Experience) == 0 {
        t.Error("expected Experience to be non-empty")
    }
    if len(resume.Projects) == 0 {
        t.Error("expected Projects to be non-empty")
    }

    t.Logf("Loaded resume for: %s (%s)", resume.Name, resume.Title)
    t.Logf("Experience count: %d", len(resume.Experience))
    t.Logf("Projects count: %d", len(resume.Projects))
}
```

运行：
```bash
go test ./internal/data/...
```

### 扩展测试思路

```go
// 测试默认数据完整性
func TestDefaultResume(t *testing.T) {
    r := DefaultResume()
    if r.Name == "" {
        t.Error("DefaultResume Name is empty")
    }
    // 更多断言...
}

// 测试页面切换逻辑（需要导出或重构）
func TestPageNavigation(t *testing.T) {
    m := NewModel()
    // 模拟 nextPage
    m.nextPage()
    if m.currentPage != style.AboutPage {
        t.Errorf("expected AboutPage, got %v", m.currentPage)
    }
}
```

## 调试技巧

### 技巧 1：日志输出

Bubble Tea 占用标准输出，日志需要重定向到文件：

```go
// main.go 中添加
f, _ := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
log.SetOutput(f)
log.Println("Debug message")
```

### 技巧 2：模型状态打印

```go
// 在 Update 中临时添加调试输出
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    log.Printf("Update: msg=%T, page=%v, cursor=%d", msg, m.currentPage, m.cursor)
    // ...
}
```

### 技巧 3：使用 `tea.WithOutput()`

将 Bubble Tea 输出到 bytes.Buffer 进行测试：

```go
var buf bytes.Buffer
p := tea.NewProgram(m,
    tea.WithOutput(&buf),
    tea.WithoutRenderer(),  // 不渲染到真实终端
)
```

### 技巧 4：最小复现

遇到 Bug 时，创建最小可复现程序：

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
)

type testModel struct{}

func (m testModel) Init() tea.Cmd { return nil }
func (m testModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if key, ok := msg.(tea.KeyMsg); ok {
        println("Key:", key.String())  // 直接打印到 stderr
    }
    return m, nil
}
func (m testModel) View() string { return "Press keys, Ctrl+C to quit\n" }

func main() {
    tea.NewProgram(testModel{}, tea.WithAltScreen()).Run()
}
```

## 常见问题排查

| 现象 | 可能原因 | 解决 |
|------|----------|------|
| 界面闪烁 | 终端不支持双缓冲 | 换用现代终端 |
| 颜色异常 | 终端色深不足 | 使用 256 色或真彩色终端 |
| 中文乱码 | 编码不是 UTF-8 | 设置 `LANG=en_US.UTF-8` |
| 鼠标无效 | 未启用 MouseCellMotion | 添加 `tea.WithMouseCellMotion()` |
| 窗口尺寸为 0 | 未处理 WindowSizeMsg | 在 Update 中处理 |

## 本章小结

- `cmd/local/main.go` 提供无需 SSH 的快速调试入口
- `go test ./...` 运行所有单元测试
- 日志重定向到文件避免与 TUI 输出冲突
- 最小复现程序是调试的利器

## 下一步

[10 - 部署与运维](10-deployment.md)
