# 07 - 页面渲染：各板块实现

本章实现 6 个页面的渲染逻辑。所有渲染方法返回 `string`，由 Lipgloss 样式包装。

## Dashboard 首页

Dashboard 是应用的"门面"，包含：ASCII 艺术标题 + 名片信息 + 快速导航提示。

```go
func (m Model) renderDashboard() string {
    r := m.resume

    // ASCII 艺术字（可用 figlet 生成）
    title := `
   ____  _           _         __  __
  / ___|(_)_ __ ___ | | ___   |  \/  | ___ _ __   __ _  ___ _   _ ___
  \___ \| | '_ ' _ \| |/ _ \  | |\/| |/ _ \ '_ \ / _' |/ _ \ | | / __|
   ___) | | | | | | | |  __/  | |  | |  __/ | | | (_| |  __/ |_| \__ \
  |____/|_|_| |_| |_|_|\___|  |_|  |_|\___|_| |_|\__, |\___|\__,_|___/
                                                 |___/
`
    titleStyled := style.TitleStyle.Render(title)

    // 名片信息
    info := lipgloss.JoinVertical(lipgloss.Left,
        style.HeaderStyle.Render(fmt.Sprintf("👤  %s", r.Name)),
        style.SubtitleStyle.Render(fmt.Sprintf("    %s", r.Title)),
        "",
        style.KeyStyle.Render("📍  ")+style.ValueStyle.Render(r.Location),
        style.KeyStyle.Render("📧  ")+style.ValueStyle.Render(r.Email),
        style.KeyStyle.Render("🐙  ")+style.ValueStyle.Render(r.GitHub),
        style.KeyStyle.Render("🌐  ")+style.ValueStyle.Render(r.Website),
    )

    // 快速导航
    quickNav := lipgloss.JoinVertical(lipgloss.Left,
        "",
        style.SeparatorStyle.Render(strings.Repeat("─", 50)),
        style.HighlightText.Render("🚀 Quick Navigation"),
        "",
        style.KeyStyle.Render("[1]")+" About Me    "+style.KeyStyle.Render("[2]")+" Experience",
        style.KeyStyle.Render("[3]")+" Skills      "+style.KeyStyle.Render("[4]")+" Projects",
        style.KeyStyle.Render("[5]")+" Contact",
        "",
        style.DimText.Render("Use arrow keys or h/j/k/l to navigate"),
    )

    content := lipgloss.JoinVertical(lipgloss.Center,
        titleStyled,
        style.BoxStyle.Render(info),
        quickNav,
    )

    // 整体居中
    return lipgloss.Place(m.width, m.height-3,
        lipgloss.Center, lipgloss.Center, content)
}
```

**要点**：
- `lipgloss.Place()` 将内容在指定区域内居中
- 高度减 3 是为底部导航栏留空间
- ASCII 艺术可用 [figlet](http://www.figlet.org/) 生成

## About 页面

```go
func (m Model) renderAbout() string {
    r := m.resume.About

    var b strings.Builder
    b.WriteString(style.HeaderStyle.Render("📋 About Me") + "\n\n")
    b.WriteString(style.NormalText.Render(r.Summary) + "\n\n")
    b.WriteString(style.SeparatorStyle.Render(strings.Repeat("─", m.width-4)) + "\n\n")
    b.WriteString(style.HighlightText.Render("✨ Highlights") + "\n\n")

    for _, detail := range r.Details {
        b.WriteString(style.KeyStyle.Render("  ▸ ") + style.ValueStyle.Render(detail) + "\n")
    }

    return style.BoxStyle.Width(m.width - 4).Render(b.String())
}
```

## Experience 页面（带光标高亮）

```go
func (m Model) renderExperience() string {
    var b strings.Builder
    b.WriteString(style.HeaderStyle.Render("💼 Experience") + "\n\n")

    for i, exp := range m.resume.Experience {
        // 构建表头
        header := style.KeyStyle.Render(fmt.Sprintf("▸ %s", exp.Company)) + "  " +
            style.HighlightText.Render(exp.Role) + "  " +
            style.DimText.Render(exp.Period)

        // 光标高亮当前项
        if i == m.cursor%len(m.resume.Experience) {
            header = style.SelectedItem.Render(fmt.Sprintf(" %s ", exp.Company)) + "  " +
                style.HighlightText.Render(exp.Role) + "  " +
                style.DimText.Render(exp.Period)
        }

        b.WriteString(header + "\n")
        b.WriteString(style.DimText.Render(fmt.Sprintf("  📍 %s", exp.Location)) + "\n\n")

        for _, highlight := range exp.Highlights {
            b.WriteString(style.NormalText.Render(fmt.Sprintf("    • %s", highlight)) + "\n")
        }

        b.WriteString("\n" + style.SeparatorStyle.Render(strings.Repeat("─", m.width-8)) + "\n\n")
    }

    return style.BoxStyle.Width(m.width - 4).Render(b.String())
}
```

**光标高亮原理**：
- `m.cursor` 记录当前选中索引
- `i == m.cursor % len(list)` 确保光标在列表范围内循环
- 选中项使用 `SelectedItem` 样式（绿底黑字反色）

## Skills 页面

```go
func (m Model) renderSkills() string {
    s := m.resume.Skills
    var b strings.Builder
    b.WriteString(style.HeaderStyle.Render("🛠  Skills") + "\n\n")

    sections := []struct {
        Title string
        Items []string
    }{
        {"Programming Languages", s.Languages},
        {"Frameworks & Libraries", s.Frameworks},
        {"Tools & Platforms", s.Tools},
        {"Other Expertise", s.Others},
    }

    for _, section := range sections {
        b.WriteString(style.HighlightText.Render(section.Title) + "\n")
        b.WriteString(style.DimText.Render(strings.Repeat("─", len(section.Title))) + "\n")

        for _, item := range section.Items {
            b.WriteString(style.KeyStyle.Render("  ◆ ") + style.ValueStyle.Render(item) + "\n")
        }
        b.WriteString("\n")
    }

    return style.BoxStyle.Width(m.width - 4).Render(b.String())
}
```

## Projects 页面

```go
func (m Model) renderProjects() string {
    var b strings.Builder
    b.WriteString(style.HeaderStyle.Render("📁 Projects") + "\n\n")

    for i, proj := range m.resume.Projects {
        name := style.KeyStyle.Render(fmt.Sprintf("▸ %s", proj.Name))
        if i == m.cursor%len(m.resume.Projects) {
            name = style.SelectedItem.Render(fmt.Sprintf(" %s ", proj.Name))
        }

        b.WriteString(name + "\n")
        b.WriteString(style.NormalText.Render(fmt.Sprintf("  %s", proj.Description)) + "\n")
        b.WriteString(style.DimText.Render(fmt.Sprintf("  🔗 %s", proj.URL)) + "\n\n")

        // Tech Stack 标签
        var techs []string
        for _, tech := range proj.TechStack {
            techs = append(techs, style.UnselectedItem.Render(tech))
        }
        b.WriteString("  " + strings.Join(techs, " ") + "\n\n")

        for _, highlight := range proj.Highlights {
            b.WriteString(style.NormalText.Render(fmt.Sprintf("    • %s", highlight)) + "\n")
        }

        b.WriteString("\n" + style.SeparatorStyle.Render(strings.Repeat("─", m.width-8)) + "\n\n")
    }

    return style.BoxStyle.Width(m.width - 4).Render(b.String())
}
```

## Contact 页面

```go
func (m Model) renderContact() string {
    c := m.resume.Contact
    var b strings.Builder
    b.WriteString(style.HeaderStyle.Render("📬 Contact") + "\n\n")

    contacts := []struct {
        Icon string
        Key  string
        Val  string
    }{
        {"📧", "Email", c.Email},
        {"🐙", "GitHub", c.GitHub},
        {"🌐", "Website", c.Website},
        {"💼", "LinkedIn", c.LinkedIn},
        {"🐦", "Twitter", c.Twitter},
    }

    for _, ct := range contacts {
        b.WriteString(fmt.Sprintf("%s  ", ct.Icon))
        b.WriteString(style.KeyStyle.Render(fmt.Sprintf("%-10s", ct.Key)))
        b.WriteString(style.ValueStyle.Render(ct.Val) + "\n\n")
    }

    b.WriteString("\n" + style.SeparatorStyle.Render(strings.Repeat("─", m.width-8)) + "\n\n")
    b.WriteString(style.HighlightText.Render(c.Message) + "\n")

    return style.BoxStyle.Width(m.width - 4).Render(b.String())
}
```

## 渲染技巧总结

### 1. strings.Builder vs 直接拼接

```go
// ✅ 推荐：使用 strings.Builder，高效且清晰
var b strings.Builder
b.WriteString("...")
b.WriteString("...")
result := b.String()

// ❌ 避免：反复创建字符串，效率低
result := ""
result += "..."
result += "..."
```

### 2. 宽度约束

```go
// 统一留 4 字符边距（左右各 2）
style.BoxStyle.Width(m.width - 4).Render(content)
```

### 3. 分隔线动态长度

```go
// 根据可用宽度动态调整
sep := strings.Repeat("─", m.width-8)  // 再减 4 给内边距
```

### 4. 条件样式模式

```go
var style lipgloss.Style
if isSelected {
    style = style.SelectedItem
} else {
    style = style.UnselectedItem
}
rendered := style.Render(text)
```

## 本章小结

- 每个页面是独立的渲染方法，返回样式化字符串
- 使用 `strings.Builder` 高效构建多行内容
- `cursor % len(list)` 实现循环高亮
- `BoxStyle.Width()` 约束内容宽度防止溢出

## 下一步

[08 - SSH 服务器：Wish 框架集成](08-ssh-server.md)
