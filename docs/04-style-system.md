# 04 - 样式系统：复古终端主题

## 设计目标

创建一套**可复用**的复古终端样式，统一整个应用的视觉风格。灵感来源：老式 CRT 显示器的磷光绿 + 黑色背景。

## 配色方案

```
┌─────────────────────────────────────────────┐
│          复古终端配色表                       │
├─────────────────────────────────────────────┤
│  #0a0a0a  Black     │  背景色               │
│  #00ff41  Green      │  主色调（高亮文字）    │
│  #008f11  DimGreen   │  次要绿（边框、分隔线） │
│  #003b00  DarkGreen  │  暗绿（装饰）          │
│  #e0e0e0  White      │  正文文字             │
│  #808080  Gray       │  淡化文字             │
│  #ffb000  Amber      │  强调色（标题）        │
│  #ff3333  Red        │  错误提示             │
└─────────────────────────────────────────────┘
```

## Lipgloss 基础

Lipgloss 是一个声明式样式库，核心概念：

```go
// 创建样式
style := lipgloss.NewStyle().
    Foreground(lipgloss.Color("#00ff41")).
    Bold(true).
    Padding(1, 2).           // 上下 1，左右 2
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("#008f11"))

// 应用样式
styledText := style.Render("Hello, Terminal!")
```

### 常用链式方法

| 方法 | 作用 | 示例值 |
|------|------|--------|
| `Foreground` | 文字颜色 | `lipgloss.Color("#00ff41")` |
| `Background` | 背景颜色 | `lipgloss.Color("#003b00")` |
| `Bold` | 加粗 | `true` / `false` |
| `Italic` | 斜体 | `true` / `false` |
| `Padding` | 内边距 | `Padding(1, 2)` |
| `Margin` | 外边距 | `Margin(1, 2)` |
| `BorderStyle` | 边框样式 | `RoundedBorder()`, `NormalBorder()` |
| `BorderForeground` | 边框颜色 | `lipgloss.Color(...)` |
| `Width` / `Height` | 固定宽高 | `Width(40)` |

### 布局方法

```go
// 水平拼接
line := lipgloss.JoinHorizontal(lipgloss.Top, box1, box2, box3)

// 垂直拼接
column := lipgloss.JoinVertical(lipgloss.Left, header, body, footer)

// 居中对齐
centered := lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
```

## 实现主题文件

创建 `internal/style/theme.go`：

### 1. 颜色常量

```go
var (
    Black     = lipgloss.Color("#0a0a0a")
    Green     = lipgloss.Color("#00ff41")
    DimGreen  = lipgloss.Color("#008f11")
    DarkGreen = lipgloss.Color("#003b00")
    White     = lipgloss.Color("#e0e0e0")
    Gray      = lipgloss.Color("#808080")
    Amber     = lipgloss.Color("#ffb000")
    Red       = lipgloss.Color("#ff3333")
)
```

### 2. 预定义样式

```go
var (
    // 标题样式：大号绿色加粗
    TitleStyle = lipgloss.NewStyle().
        Foreground(Green).
        Bold(true).
        PaddingTop(1).
        PaddingBottom(1)

    // 副标题：暗绿色斜体
    SubtitleStyle = lipgloss.NewStyle().
        Foreground(DimGreen).
        Italic(true)

    // 正文：白色
    NormalText = lipgloss.NewStyle().
        Foreground(White)

    // 淡化文字：灰色
    DimText = lipgloss.NewStyle().
        Foreground(Gray)

    // 高亮文字：绿色加粗
    HighlightText = lipgloss.NewStyle().
        Foreground(Green).
        Bold(true)

    // 盒子样式：圆角边框
    BoxStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.RoundedBorder()).
        BorderForeground(DimGreen).
        Padding(1, 2)

    // 选中项：绿底黑字（反色）
    SelectedItem = lipgloss.NewStyle().
        Foreground(Black).
        Background(Green).
        Bold(true).
        Padding(0, 1)

    // 未选中项：绿色文字
    UnselectedItem = lipgloss.NewStyle().
        Foreground(Green).
        Padding(0, 1)

    // 页面标题：琥珀色
    HeaderStyle = lipgloss.NewStyle().
        Foreground(Amber).
        Bold(true)

    // 键名：绿色加粗
    KeyStyle = lipgloss.NewStyle().
        Foreground(Green).
        Bold(true)

    // 值：白色
    ValueStyle = lipgloss.NewStyle().
        Foreground(White)

    // 分隔线：暗绿色
    SeparatorStyle = lipgloss.NewStyle().
        Foreground(DimGreen)

    // 底部帮助：灰色斜体
    FooterStyle = lipgloss.NewStyle().
        Foreground(Gray).
        Italic(true)

    // 错误：红色加粗
    ErrorStyle = lipgloss.NewStyle().
        Foreground(Red).
        Bold(true)
)
```

### 3. 页面类型枚举

```go
// PageType 定义页面类型
type PageType int

const (
    DashboardPage PageType = iota  // 首页
    AboutPage                       // 关于我
    ExperiencePage                  // 工作经历
    SkillsPage                      // 技能
    ProjectsPage                    // 项目
    ContactPage                     // 联系方式
)

func (p PageType) String() string {
    switch p {
    case DashboardPage: return "Dashboard"
    case AboutPage:     return "About"
    // ... 其他 case
    default:            return "Unknown"
    }
}

// AllPages 返回所有内容页面（不含 Dashboard）
func AllPages() []PageType {
    return []PageType{AboutPage, ExperiencePage, SkillsPage, ProjectsPage, ContactPage}
}
```

## 样式设计原则

1. **语义化命名**：`TitleStyle` 而不是 `GreenBoldStyle`
2. **集中管理**：所有样式放在 `style` 包，避免分散
3. **复用优先**：定义通用样式，而非每个页面单独写
4. **约束宽度**：使用 `.Width()` 确保内容不溢出终端

## 实战技巧

### 技巧 1：动态宽度

```go
// 根据终端宽度自适应
contentWidth := m.width - 4  // 留边距
rendered := style.BoxStyle.Width(contentWidth).Render(content)
```

### 技巧 2：字符串重复做分隔线

```go
separator := style.SeparatorStyle.Render(strings.Repeat("─", width))
```

### 技巧 3：条件样式

```go
var name string
if isSelected {
    name = style.SelectedItem.Render(item)
} else {
    name = style.UnselectedItem.Render(item)
}
```

## 本章小结

- Lipgloss 提供声明式终端样式 API
- 设计复古绿黑配色方案
- 预定义语义化样式，集中管理
- 页面类型使用 Go 枚举（iota）

## 下一步

[05 - TUI 基础：Bubble Tea 框架入门](05-tui-basics.md)
