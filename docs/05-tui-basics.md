# 05 - TUI 基础：Bubble Tea 框架入门

## 什么是 TUI？

TUI（Terminal User Interface）是在终端中渲染的图形界面，区别于：

- **CLI**：命令行输入输出，无交互界面
- **GUI**：图形窗口界面，需要桌面环境
- **TUI**：终端内的交互界面，键盘驱动

## Elm 架构

Bubble Tea 采用 Elm 架构，将程序分为三个纯函数：

```
┌─────────────┐     Msg      ┌─────────────┐
│   Init()    │ ───────────► │   Update()  │
│  (初始状态)  │              │ (处理消息)   │
└─────────────┘              └──────┬──────┘
                                    │
                                    │ New Model
                                    ▼
┌─────────────┐     Cmd      ┌─────────────┐
│   View()    │ ◄─────────── │   Model     │
│  (渲染界面)  │              │  (当前状态)  │
└─────────────┘              └─────────────┘
        │
        │ 用户看到界面，按下按键
        ▼
     New Msg (循环)
```

### 核心接口

```go
// Model 接口：任何 struct 只要实现这三个方法就是 Bubble Tea 模型
type Model interface {
    Init() Cmd           // 初始化命令（如启动定时器）
    Update(Msg) (Model, Cmd)  // 处理消息，返回新状态和后续命令
    View() string        // 渲染当前状态为字符串
}
```

## 最简单的 Bubble Tea 程序

```go
package main

import (
    "fmt"
    "os"
    tea "github.com/charmbracelet/bubbletea"
)

// 定义模型
type model struct {
    count int
}

func (m model) Init() tea.Cmd {
    return nil  // 无初始命令
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit  // 退出程序
        case "up":
            m.count++
        case "down":
            m.count--
        }
    }
    return m, nil
}

func (m model) View() string {
    return fmt.Sprintf("Count: %d\n\nPress ↑/↓ to change, q to quit", m.count)
}

func main() {
    p := tea.NewProgram(model{})
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## 消息类型（Msg）

Bubble Tea 中一切都是消息：

### 内置消息

| 消息类型 | 触发时机 | 用途 |
|----------|----------|------|
| `tea.KeyMsg` | 用户按键 | 处理键盘输入 |
| `tea.WindowSizeMsg` | 终端尺寸变化 | 获取 width/height |
| `tea.MouseMsg` | 鼠标事件 | 处理鼠标点击/滚动 |
| `tea.QuitMsg` | 退出信号 | 清理资源 |

### 按键消息详解

```go
case tea.KeyMsg:
    switch msg.String() {
    case "q", "Q":           // q 或 Q
    case "ctrl+c":           // Ctrl+C
    case "up", "k":          // ↑ 或 k（Vim 风格）
    case "down", "j":        // ↓ 或 j
    case "left", "h":        // ← 或 h
    case "right", "l":       // → 或 l
    case "1", "2", "3":      // 数字键
    case "enter":            // 回车
    case "esc":              // ESC
    case "tab":              // Tab
    case " ":                // 空格
    }
```

### 自定义消息

```go
// 定义自定义消息
type tickMsg time.Time

// 发送消息的命令
func tickCmd() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

// 在 Update 中处理
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tickMsg:
        m.count++
        return m, tickCmd()  // 继续定时器
    }
    // ...
}
```

## 程序选项

```go
p := tea.NewProgram(model{},
    tea.WithAltScreen(),        // 使用备用屏幕（全屏，退出后恢复）
    tea.WithMouseCellMotion(),  // 启用鼠标支持
    tea.WithInput(sshInput),    // 自定义输入（SSH 场景）
)
```

| 选项 | 作用 |
|------|------|
| `WithAltScreen()` | 进入全屏模式，退出后终端内容不被破坏 |
| `WithMouseCellMotion()` | 启用鼠标事件，支持点击和滚动 |

## 本项目的模型设计

```go
type Model struct {
    width       int         // 终端宽度（用于自适应布局）
    height      int         // 终端高度
    resume      *data.Resume // 简历数据
    currentPage style.PageType // 当前页面
    cursor      int         // 滚动位置 / 选中项
    ready       bool        // 是否已完成初始化（获取到窗口尺寸）
}
```

### 为什么需要 `ready`？

终端尺寸是异步获取的。程序启动时 `width`/`height` 为 0，收到 `tea.WindowSizeMsg` 后才变为真实值。`ready` 标记避免在尺寸未知时渲染：

```go
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }
    // ... 正常渲染
}
```

## 本章小结

- Bubble Tea 使用 Elm 架构：Model/Update/View/Cmd
- 所有交互通过消息（Msg）驱动
- `tea.KeyMsg` 处理键盘，`tea.WindowSizeMsg` 处理尺寸变化
- `WithAltScreen()` 提供全屏 TUI 体验
- `ready` 标记处理异步初始化

## 下一步

[06 - 应用模型：页面与导航](06-app-model.md)
