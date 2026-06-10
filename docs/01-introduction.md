# 01 - 课程介绍与技术选型

## 项目背景

传统简历通常以 PDF 或网页形式呈现。本项目探索一种更具极客风格的展示方式：通过 SSH 协议直接在终端中展示交互式简历。这种形式有以下优势：

1. **差异化** — 在 HR 和面试官面前留下深刻印象
2. **技术展示** — 本身就是技术能力的体现
3. **轻量访问** — 任何有终端的设备都能访问，无需浏览器
4. **复古美学** — 终端风格自带程序员文化认同感

灵感来源：[terminal.shop](https://terminal.shop) — 一个完全通过 SSH 访问的咖啡订购商店。

## 技术选型

### 核心库：Charm 生态

本项目基于 [Charm](https://charm.sh) 公司开源的三个核心库构建：

```
┌─────────────────────────────────────────────────────────────┐
│                      应用架构分层                             │
├─────────────────────────────────────────────────────────────┤
│  表现层 (View)          │  Lipgloss — 声明式终端样式          │
│  ─────────────────────  │  颜色、边框、布局、文字样式           │
│  交互层 (Update/Model)  │  Bubble Tea — TUI 框架              │
│  ─────────────────────  │  Elm 架构：Model/Update/View/Cmd    │
│  传输层 (SSH Server)    │  Wish — SSH 服务器框架              │
│  ─────────────────────  │  将 SSH 会话包装为 Bubble Tea 程序  │
└─────────────────────────────────────────────────────────────┘
```

### 各库定位详解

#### 1. Lipgloss — 终端 CSS

```go
// 声明式样式定义，类似 CSS
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#00ff41")).  // 文字颜色
    Bold(true).                              // 加粗
    Padding(1, 2).                          // 内边距
    BorderStyle(liposs.RoundedBorder())      // 圆角边框
```

**解决的问题**：ANSI 转义码难以手工编写和维护。

#### 2. Bubble Tea — TUI 框架

采用 [Elm 架构](https://guide.elm-lang.org/architecture/)，核心三个组件：

- **Model**：应用状态（当前页面、光标位置、窗口尺寸）
- **Update(msg)**：接收消息（按键、窗口变化），返回新 Model
- **View()**：根据 Model 渲染界面

```go
type Model struct {
    currentPage PageType
    cursor      int
    width       int
    height      int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 处理输入，更新状态
}

func (m Model) View() string {
    // 根据状态渲染字符串
}
```

**解决的问题**：终端输入处理和屏幕刷新复杂，手工实现容易出错。

#### 3. Wish — SSH 应用服务器

```go
server, _ := wish.NewServer(
    wish.WithAddress("localhost:23234"),
    wish.WithHostKeyPath(".ssh/host_key"),
    wish.WithMiddleware(
        bubbletea.Middleware(teaHandler),  // 核心：将 SSH 会话转为 TUI
    ),
)
```

**解决的问题**：传统 SSH 服务器需要处理认证、会话管理、PTY 分配等复杂逻辑。

### 辅助库

| 库 | 用途 |
|---|---|
| `yaml.v3` | 解析简历配置文件 |
| `go:embed` | 将 YAML 嵌入编译后的二进制文件 |

## 为什么不选其他方案？

| 方案 | 缺点 | 本方案优势 |
|------|------|-----------|
| ncurses (C) | API 老旧，Go 绑定不成熟 | 纯 Go，现代 API |
| tview (Go) | 功能强大但学习曲线陡 | Bubble Tea 更简单直观 |
| 纯 Web | 需要浏览器，不够极客 | SSH 直接访问，终端原生 |
| 纯文本 + cat | 无交互，无法滚动 | 完全交互式体验 |

## 本章小结

- 项目目标是构建 SSH 访问的交互式简历 TUI
- 核心技术栈：Wish (SSH) + Bubble Tea (TUI) + Lipgloss (样式)
- 三个库分工明确：传输层 → 交互层 → 表现层

## 下一步

[02 - 环境准备与项目初始化](02-setup.md)
