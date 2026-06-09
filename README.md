# 🖥️ SSH Resume

一个基于 **SSH** 的交互式个人简历终端应用，灵感来自 [terminal.shop](https://terminal.shop) 的复古终端美学。

无需浏览器，直接通过 `ssh` 命令即可访问你的个人简历！

```bash
ssh resume.yourdomain.com
```

![技术栈](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Wish](https://img.shields.io/badge/Wish-SSH%20Server-ff69b4)
![BubbleTea](https://img.shields.io/badge/Bubble%20Tea-TUI-7B68EE)

---

## ✨ 特性

- 🚀 **SSH 直接访问** - 无需浏览器，终端即简历
- 🎨 **复古终端风格** - 经典绿黑配色 + CRT 扫描线氛围
- ⌨️ **键盘驱动导航** - Vim 风格快捷键，极客体验
- 📄 **YAML 配置** - 修改 `resume.yaml` 即可更新简历内容
- 🖥️ **自适应终端** - 自动适配各种终端尺寸

---

## 🚀 快速开始

### 本地运行（无需 SSH）

```bash
go run ./cmd/local
```

### 启动 SSH 服务器

```bash
# 生成 SSH Host Key（首次运行）
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N ""

# 启动服务器
go run .

# 另开终端连接
ssh localhost -p 23234
```

---

## 📁 项目结构

```
ssh-resume/
├── main.go                    # SSH 服务器入口 (Wish)
├── cmd/local/main.go          # 本地测试入口（无需 SSH）
├── internal/
│   ├── app/
│   │   └── model.go           # Bubble Tea TUI 模型
│   ├── style/
│   │   └── theme.go           # 复古终端配色 & 样式
│   └── data/
│       ├── resume.go          # 数据结构与默认值
│       ├── resume.yaml        # 📄 简历配置文件（修改这里！）
│       ├── loader.go          # YAML 加载器（go:embed）
│       └── loader_test.go     # 测试
├── go.mod
└── README.md
```

---

## 📝 自定义简历

编辑 `internal/data/resume.yaml`：

```yaml
name: "你的名字"
title: "你的职位"
location: "你的城市"
email: "your@email.com"
github: "github.com/yourname"
website: "yourdomain.com"

about:
  summary: "一句话介绍自己"
  details:
    - "亮点 1"
    - "亮点 2"

experience:
  - company: "公司名"
    role: "职位"
    period: "2020 - 至今"
    location: "城市"
    highlights:
      - "工作成就 1"
      - "工作成就 2"

skills:
  languages: ["Go", "Rust", "Python"]
  frameworks: ["Gin", "React"]
  tools: ["Docker", "K8s"]
  others: ["Microservices"]

projects:
  - name: "项目名称"
    description: "项目简介"
    techStack: ["Go", "Redis"]
    url: "github.com/you/project"
    highlights:
      - "项目亮点"

contact:
  email: "your@email.com"
  github: "github.com/you"
  message: "欢迎联系我！"
```

修改后重新编译即可：

```bash
go build -o ssh-resume .
```

---

## 🎮 操作指南

| 按键 | 功能 |
|------|------|
| `←` / `h` | 切换到上一板块 |
| `→` / `l` | 切换到下一板块 |
| `↑` / `k` | 向上滚动 |
| `↓` / `j` | 向下滚动 |
| `1` ~ `5` | 快速跳转到对应板块 |
| `q` | 返回首页 / 退出 |
| `Ctrl+C` | 强制退出 |

---

## 🐳 Docker 部署

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o ssh-resume .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/ssh-resume .
COPY --from=builder /app/.ssh ./.ssh
EXPOSE 23234
CMD ["./ssh-resume"]
```

```bash
docker build -t ssh-resume .
docker run -p 23234:23234 ssh-resume
```

---

## ☁️ Fly.io 部署（推荐）

```bash
# 安装 flyctl
# https://fly.io/docs/hands-on/install-flyctl/

fly launch
fly deploy
```

---

## 🛠️ 技术栈

| 库 | 用途 |
|---|---|
| [wish](https://github.com/charmbracelet/wish) | SSH 服务器框架 |
| [bubbletea](https://github.com/charmbracelet/bubbletea) | TUI 应用框架 |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | 终端样式美化 |
| [yaml.v3](https://gopkg.in/yaml.v3) | 简历配置解析 |

---

## 📄 License

MIT License © 2024
