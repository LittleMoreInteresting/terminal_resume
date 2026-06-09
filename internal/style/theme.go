package style

import (
	"github.com/charmbracelet/lipgloss"
)

// 复古终端配色方案 - 经典绿黑主题
var (
	// 基础颜色
	Black   = lipgloss.Color("#0a0a0a")
	Green   = lipgloss.Color("#00ff41")
	DimGreen = lipgloss.Color("#008f11")
	DarkGreen = lipgloss.Color("#003b00")
	White   = lipgloss.Color("#e0e0e0")
	Gray    = lipgloss.Color("#808080")
	Amber   = lipgloss.Color("#ffb000")
	Red     = lipgloss.Color("#ff3333")

	// 常用样式
	TitleStyle = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true).
		PaddingTop(1).
		PaddingBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
		Foreground(DimGreen).
		Italic(true)

	NormalText = lipgloss.NewStyle().
		Foreground(White)

	DimText = lipgloss.NewStyle().
		Foreground(Gray)

	HighlightText = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	BoxStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(DimGreen).
		Padding(1, 2)

	SelectedItem = lipgloss.NewStyle().
		Foreground(Black).
		Background(Green).
		Bold(true).
		Padding(0, 1)

	UnselectedItem = lipgloss.NewStyle().
		Foreground(Green).
		Padding(0, 1)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(Amber).
		Bold(true)

	KeyStyle = lipgloss.NewStyle().
		Foreground(Green).
		Bold(true)

	ValueStyle = lipgloss.NewStyle().
		Foreground(White)

	SeparatorStyle = lipgloss.NewStyle().
		Foreground(DimGreen)

	FooterStyle = lipgloss.NewStyle().
		Foreground(Gray).
		Italic(true)

	CursorStyle = lipgloss.NewStyle().
		Foreground(Green).
		Blink(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(Red).
		Bold(true)
)

// PageType 定义页面类型
type PageType int

const (
	DashboardPage PageType = iota
	AboutPage
	ExperiencePage
	SkillsPage
	ProjectsPage
	ContactPage
)

func (p PageType) String() string {
	switch p {
	case DashboardPage:
		return "Dashboard"
	case AboutPage:
		return "About"
	case ExperiencePage:
		return "Experience"
	case SkillsPage:
		return "Skills"
	case ProjectsPage:
		return "Projects"
	case ContactPage:
		return "Contact"
	default:
		return "Unknown"
	}
}

// AllPages 返回所有页面类型
func AllPages() []PageType {
	return []PageType{
		AboutPage,
		ExperiencePage,
		SkillsPage,
		ProjectsPage,
		ContactPage,
	}
}

// PageTitles 页面标题映射
var PageTitles = map[PageType]string{
	AboutPage:      "About Me",
	ExperiencePage: "Experience",
	SkillsPage:     "Skills",
	ProjectsPage:   "Projects",
	ContactPage:    "Contact",
}
