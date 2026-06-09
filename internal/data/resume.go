package data

// Resume 简历数据结构
type Resume struct {
	Name        string       `yaml:"name"`
	Title       string       `yaml:"title"`
	Location    string       `yaml:"location"`
	Email       string       `yaml:"email"`
	GitHub      string       `yaml:"github"`
	Website     string       `yaml:"website"`
	About       About        `yaml:"about"`
	Experience  []Experience `yaml:"experience"`
	Skills      Skills       `yaml:"skills"`
	Projects    []Project    `yaml:"projects"`
	Contact     Contact      `yaml:"contact"`
}

// About 关于我
type About struct {
	Summary string   `yaml:"summary"`
	Details []string `yaml:"details"`
}

// Experience 工作经历
type Experience struct {
	Company    string   `yaml:"company"`
	Role       string   `yaml:"role"`
	Period     string   `yaml:"period"`
	Location   string   `yaml:"location"`
	Highlights []string `yaml:"highlights"`
}

// Skills 技能
type Skills struct {
	Languages  []string `yaml:"languages"`
	Frameworks []string `yaml:"frameworks"`
	Tools      []string `yaml:"tools"`
	Others     []string `yaml:"others"`
}

// Project 项目
type Project struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	TechStack   []string `yaml:"techStack"`
	URL         string   `yaml:"url"`
	Highlights  []string `yaml:"highlights"`
}

// Contact 联系方式
type Contact struct {
	Email    string `yaml:"email"`
	GitHub   string `yaml:"github"`
	Website  string `yaml:"website"`
	LinkedIn string `yaml:"linkedIn"`
	Twitter  string `yaml:"twitter"`
	Message  string `yaml:"message"`
}

// DefaultResume 返回默认简历数据（作为 YAML 加载失败的后备）
func DefaultResume() *Resume {
	return &Resume{
		Name:     "张三",
		Title:    "全栈工程师",
		Location: "中国 · 北京",
		Email:    "zhangsan@example.com",
		GitHub:   "github.com/zhangsan",
		Website:  "zhangsan.dev",
		About: About{
			Summary: "热爱技术的全栈开发者，专注于构建高性能、可扩展的分布式系统。拥有5年后端开发经验，对云原生技术和开源社区充满热情。",
			Details: []string{
				"擅长 Go、Rust、TypeScript 等现代编程语言",
				"熟悉 Kubernetes、Docker 等云原生技术栈",
				"积极参与开源社区，多个项目的核心贡献者",
				"喜欢探索新技术，热衷于技术分享和写作",
			},
		},
		Experience: []Experience{
			{
				Company:  "字节跳动",
				Role:     "高级后端工程师",
				Period:   "2022.03 - 至今",
				Location: "北京",
				Highlights: []string{
					"负责核心业务系统的架构设计与开发，日均 QPS 超过 100 万",
					"主导微服务架构升级，将单体应用拆分为 50+ 微服务",
					"设计并实现高可用缓存系统，P99 延迟降低 60%",
					"培养团队 5 名初级工程师，建立代码审查和技术分享机制",
				},
			},
			{
				Company:  "美团",
				Role:     "后端开发工程师",
				Period:   "2019.07 - 2022.02",
				Location: "北京",
				Highlights: []string{
					"参与外卖订单系统核心模块开发，服务千万级日活用户",
					"优化数据库查询性能，核心接口响应时间减少 45%",
					"开发内部 RPC 框架，被 3 个团队采纳使用",
					"获得 2021 年度优秀员工称号",
				},
			},
			{
				Company:  "阿里巴巴",
				Role:     "Java 开发工程师（实习）",
				Period:   "2018.06 - 2018.09",
				Location: "杭州",
				Highlights: []string{
					"参与淘宝商品详情页性能优化项目",
					"学习大型互联网公司的工程实践和开发流程",
				},
			},
		},
		Skills: Skills{
			Languages: []string{
				"Go (精通)", "Rust (熟练)", "TypeScript (熟练)", "Python (熟练)", "Java (熟悉)", "C/C++ (了解)",
			},
			Frameworks: []string{
				"Gin", "Echo", "React", "Vue.js", "gRPC", "Docker", "Kubernetes",
			},
			Tools: []string{
				"Linux", "Git", "Vim", "VS Code", "DataGrip", "Postman", "K8s",
			},
			Others: []string{
				"Microservices", "Distributed Systems", "CI/CD", "TDD", "Agile",
			},
		},
		Projects: []Project{
			{
				Name:        "Go-Cache",
				Description: "高性能分布式缓存系统，兼容 Redis 协议",
				TechStack:   []string{"Go", "Raft", "gRPC", "Docker"},
				URL:         "github.com/zhangsan/go-cache",
				Highlights: []string{
					"支持集群模式和主从复制，数据一致性采用 Raft 算法",
					"单机 QPS 达到 15 万，内存占用比 Redis 降低 20%",
					"GitHub Stars: 2.3k，被多家公司用于生产环境",
				},
			},
			{
				Name:        "Terminal Resume",
				Description: "基于 SSH 的交互式简历终端应用",
				TechStack:   []string{"Go", "Wish", "Bubble Tea", "Lipgloss"},
				URL:         "github.com/zhangsan/ssh-resume",
				Highlights: []string{
					"支持通过 SSH 直接访问，无需浏览器",
					"采用 TUI 设计，提供沉浸式的终端体验",
					"开源项目，获得 500+ Stars",
				},
			},
			{
				Name:        "LogInsight",
				Description: "轻量级日志收集与分析平台",
				TechStack:   []string{"Go", "ClickHouse", "Vue.js", "WebSocket"},
				URL:         "github.com/zhangsan/loginsight",
				Highlights: []string{
					"支持实时日志收集和关键字搜索，查询延迟 < 100ms",
					"内置可视化仪表盘，支持自定义告警规则",
					"日处理日志量超过 10TB",
				},
			},
		},
		Contact: Contact{
			Email:    "zhangsan@example.com",
			GitHub:   "github.com/zhangsan",
			Website:  "zhangsan.dev",
			LinkedIn: "linkedin.com/in/zhangsan",
			Twitter:  "@zhangsan_dev",
			Message:  "欢迎通过邮件或 GitHub 与我联系！",
		},
	}
}
