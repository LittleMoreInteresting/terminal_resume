package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"terminal_resume/internal/data"
	"terminal_resume/internal/style"
)

// Model 是 Bubble Tea 的主模型
type Model struct {
	width       int
	height      int
	resume      *data.Resume
	currentPage style.PageType
	cursor      int      // 用于详情页内的滚动位置
	ready       bool
}

// NewModel 创建新的 TUI 模型
func NewModel() Model {
	return Model{
		resume:      data.LoadResumeOrDefault(),
		currentPage: style.DashboardPage,
		cursor:      0,
	}
}

// Init 初始化
func (m Model) Init() tea.Cmd {
	return nil
}

// Update 处理消息和输入
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q":
			if m.currentPage == style.DashboardPage {
				return m, tea.Quit
			}
			m.currentPage = style.DashboardPage
			m.cursor = 0
			return m, nil

		case "ctrl+c":
			return m, tea.Quit

		// 数字键快速跳转
		case "1":
			m.currentPage = style.AboutPage
			m.cursor = 0
			return m, nil
		case "2":
			m.currentPage = style.ExperiencePage
			m.cursor = 0
			return m, nil
		case "3":
			m.currentPage = style.SkillsPage
			m.cursor = 0
			return m, nil
		case "4":
			m.currentPage = style.ProjectsPage
			m.cursor = 0
			return m, nil
		case "5":
			m.currentPage = style.ContactPage
			m.cursor = 0
			return m, nil

		// 左右切换板块
		case "left", "h":
			m.prevPage()
			m.cursor = 0
			return m, nil
		case "right", "l":
			m.nextPage()
			m.cursor = 0
			return m, nil

		// 上下滚动
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

// View 渲染界面
func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	// 确保最小尺寸
	if m.width < 60 || m.height < 20 {
		return style.ErrorStyle.Render("Terminal too small. Please resize to at least 60x20.")
	}

	var b strings.Builder

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

// prevPage 切换到上一页
func (m *Model) prevPage() {
	pages := style.AllPages()
	for i := len(pages) - 1; i >= 0; i-- {
		if pages[i] < m.currentPage {
			m.currentPage = pages[i]
			return
		}
	}
	// 如果已经在最前面，回到 Dashboard
	m.currentPage = style.DashboardPage
}

// nextPage 切换到下一页
func (m *Model) nextPage() {
	pages := style.AllPages()
	for i := 0; i < len(pages); i++ {
		if pages[i] > m.currentPage {
			m.currentPage = pages[i]
			return
		}
	}
	// 如果已经在最后，循环到 Dashboard
	m.currentPage = style.DashboardPage
}

// renderFooter 渲染底部导航栏
func (m Model) renderFooter() string {
	var items []string

	// 页面指示器
	if m.currentPage == style.DashboardPage {
		items = append(items, style.SelectedItem.Render("[Home]"))
	} else {
		items = append(items, style.UnselectedItem.Render("[Home]"))
	}

	for _, page := range style.AllPages() {
		label := fmt.Sprintf("[%d] %s", int(page), page.String())
		if m.currentPage == page {
			items = append(items, style.SelectedItem.Render(label))
		} else {
			items = append(items, style.UnselectedItem.Render(label))
		}
	}

	nav := strings.Join(items, "  ")

	// 帮助提示
	help := style.FooterStyle.Render("←/→ or h/l: switch  |  ↑/↓ or j/k: scroll  |  q: back/quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		style.SeparatorStyle.Render(strings.Repeat("─", m.width-2)),
		nav,
		help,
	)
}

// renderDashboard 渲染首页
func (m Model) renderDashboard() string {
	r := m.resume

	// ASCII 艺术字标题
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

	// 快速导航提示
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

	// 水平居中
	return lipgloss.Place(m.width, m.height-3,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// renderAbout 渲染关于我页面
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

// renderExperience 渲染工作经历页面
func (m Model) renderExperience() string {
	var b strings.Builder
	b.WriteString(style.HeaderStyle.Render("💼 Experience") + "\n\n")

	for i, exp := range m.resume.Experience {
		// 高亮当前选中的工作经历（基于 cursor）
		header := style.KeyStyle.Render(fmt.Sprintf("▸ %s", exp.Company)) + "  " +
			style.HighlightText.Render(exp.Role) + "  " +
			style.DimText.Render(exp.Period)

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

// renderSkills 渲染技能页面
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

// renderProjects 渲染项目页面
func (m Model) renderProjects() string {
	var b strings.Builder
	b.WriteString(style.HeaderStyle.Render("📁 Projects") + "\n\n")

	for i, proj := range m.resume.Projects {
		// 高亮当前选中的项目
		name := style.KeyStyle.Render(fmt.Sprintf("▸ %s", proj.Name))
		if i == m.cursor%len(m.resume.Projects) {
			name = style.SelectedItem.Render(fmt.Sprintf(" %s ", proj.Name))
		}

		b.WriteString(name + "\n")
		b.WriteString(style.NormalText.Render(fmt.Sprintf("  %s", proj.Description)) + "\n")
		b.WriteString(style.DimText.Render(fmt.Sprintf("  🔗 %s", proj.URL)) + "\n\n")

		// Tech Stack
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

// renderContact 渲染联系方式页面
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
