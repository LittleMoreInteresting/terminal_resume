# 06 - 应用模型：页面与导航

## 模型结构

创建 `internal/app/model.go`，实现 Bubble Tea 的三个核心方法。

### 模型定义

```go
package app

import (
    "fmt"
    "strings"
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "terminal_resume/internal/data"
    "terminal_resume/internal/style"
)

type Model struct {
    width       int
    height      int
    resume      *data.Resume
    currentPage style.PageType
    cursor      int      // 详情页内滚动位置
    ready       bool
}

func NewModel() Model {
    return Model{
        resume:      data.LoadResumeOrDefault(),
        currentPage: style.DashboardPage,
        cursor:      0,
    }
}
```

### Init 方法

```go
func (m Model) Init() tea.Cmd {
    return nil  // 本应用无需初始命令
}
```

## Update 方法：消息处理核心

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // ========== 窗口尺寸变化 ==========
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.ready = true
        return m, nil

    // ========== 键盘输入 ==========
    case tea.KeyMsg:
        switch msg.String() {

        // --- 退出逻辑 ---
        case "q", "Q":
            if m.currentPage == style.DashboardPage {
                return m, tea.Quit  // 在首页按 q 退出
            }
            m.currentPage = style.DashboardPage  // 其他页按 q 回首页
            m.cursor = 0
            return m, nil

        case "ctrl+c":
            return m, tea.Quit  // 强制退出

        // --- 数字键快速跳转 ---
        case "1":
            m.currentPage = style.AboutPage
            m.cursor = 0
            return m, nil
        case "2":
            m.currentPage = style.ExperiencePage
            m.cursor = 0
            return m, nil
        // ... 3, 4, 5 同理

        // --- 左右切换页面 ---
        case "left", "h":
            m.prevPage()
            m.cursor = 0
            return m, nil
        case "right", "l":
            m.nextPage()
            m.cursor = 0
            return m, nil

        // --- 上下滚动 ---
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }
            return m, nil
        case "down", "j":
            m.cursor++
            return m, nil
        }
    }

    return m, nil
}
```

### 导航逻辑设计

```
┌──────────────────────────────────────────────┐
│                  导航模式                       │
├──────────────────────────────────────────────┤
│                                              │
│   Dashboard ◄────────────────────────────►   │
│      │                                       │
│      │ 按 →/l 或数字键                        │
│      ▼                                       │
│   About ◄──► Experience ◄──► Skills ◄──► ... │
│                                              │
│   按 q 返回 Dashboard                         │
│                                              │
└──────────────────────────────────────────────┘
```

### prevPage / nextPage 实现

```go
func (m *Model) prevPage() {
    pages := style.AllPages()  // [About, Experience, Skills, Projects, Contact]
    for i := len(pages) - 1; i >= 0; i-- {
        if pages[i] < m.currentPage {
            m.currentPage = pages[i]
            return
        }
    }
    // 已经在最前面，回 Dashboard
    m.currentPage = style.DashboardPage
}

func (m *Model) nextPage() {
    pages := style.AllPages()
    for i := 0; i < len(pages); i++ {
        if pages[i] > m.currentPage {
            m.currentPage = pages[i]
            return
        }
    }
    // 已经在最后，回 Dashboard
    m.currentPage = style.DashboardPage
}
```

**设计要点**：
- 使用 `PageType` 的整数值进行大小比较
- `AllPages()` 返回有序的页面列表
- Dashboard 是特殊的"家"页面，不属于内容页序列

## View 方法：页面分发

```go
func (m Model) View() string {
    if !m.ready {
        return "Loading..."
    }

    // 最小尺寸检查
    if m.width < 60 || m.height < 20 {
        return style.ErrorStyle.Render("Terminal too small. Please resize to at least 60x20.")
    }

    var b strings.Builder

    // 根据当前页面分发渲染
    switch m.currentPage {
    case style.DashboardPage:
        b.WriteString(m.renderDashboard())
    case style.AboutPage:
        b.WriteString(m.renderAbout())
    case style.ExperiencePage:
        b.WriteString(m.renderExperience())
    case style.SkillsPage:
        b.WriteString(m.renderSkills())
    case style.ProjectsPage:
        b.WriteString(m.renderProjects())
    case style.ContactPage:
        b.WriteString(m.renderContact())
    }

    // 添加底部导航栏
    b.WriteString("\n")
    b.WriteString(m.renderFooter())

    return b.String()
}
```

### 底部导航栏

```go
func (m Model) renderFooter() string {
    var items []string

    // Home 按钮
    if m.currentPage == style.DashboardPage {
        items = append(items, style.SelectedItem.Render("[Home]"))
    } else {
        items = append(items, style.UnselectedItem.Render("[Home]"))
    }

    // 页面列表 [1] About [2] Experience ...
    for _, page := range style.AllPages() {
        label := fmt.Sprintf("[%d] %s", int(page), page.String())
        if m.currentPage == page {
            items = append(items, style.SelectedItem.Render(label))
        } else {
            items = append(items, style.UnselectedItem.Render(label))
        }
    }

    nav := strings.Join(items, "  ")
    help := style.FooterStyle.Render("←/→ or h/l: switch  |  ↑/↓ or j/k: scroll  |  q: back/quit")

    return lipgloss.JoinVertical(lipgloss.Left,
        style.SeparatorStyle.Render(strings.Repeat("─", m.width-2)),
        nav,
        help,
    )
}
```

## 导航设计决策

### 为什么用混合导航？

| 导航模式 | 优点 | 缺点 | 适用场景 |
|----------|------|------|----------|
| 纯命令行（`:about`） | Vim 风格，极客感强 | 对非技术用户不友好 | 个人工具 |
| 纯菜单（方向键选择） | 直观 | 长列表操作繁琐 | 简单应用 |
| **混合（本方案）** | 兼顾直观与效率 | 实现稍复杂 | **简历展示** |

本方案的混合导航：
- **Dashboard** 展示总览和快捷提示
- **←/→ 或 h/l** 在页面间快速切换
- **数字键 1-5** 直接跳转
- **↑/↓ 或 j/k** 在页面内滚动

## 本章小结

- `Update()` 是消息处理中心，分发到不同按键处理逻辑
- 页面导航使用 `prevPage()` / `nextPage()` 在有序页面间切换
- `View()` 使用 switch 分发到不同渲染方法
- 底部导航栏提供视觉反馈和操作提示

## 下一步

[07 - 页面渲染：各板块实现](07-pages.md)
