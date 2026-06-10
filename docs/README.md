# 项目实战课程：构建 SSH TUI 个人简历应用

## 课程概览

本课程将带你从零开始，使用 Go 语言和 Charm 生态库（Wish / Bubble Tea / Lipgloss）构建一个复古终端风格的 SSH TUI 简历应用。最终效果：访问者通过 `ssh your-domain.com` 即可在终端中浏览你的交互式简历。

## 目标受众

- 已掌握 Go 基础语法（struct、interface、method、goroutine）
- 想学习 TUI（Terminal User Interface）开发
- 对 SSH 应用、复古终端美学感兴趣

## 最终效果

```
ssh resume.example.com
```

一个复古绿黑配色的终端界面，包含：
- 🏠 首页 Dashboard：ASCII 艺术 + 名片信息 + 快速导航
- 📋 About：个人简介与亮点
- 💼 Experience：工作经历（支持光标高亮）
- 🛠 Skills：技能矩阵
- 📁 Projects：项目展示
- 📬 Contact：联系方式

## 课程大纲

| 章节 | 主题 | 核心知识点 |
|------|------|-----------|
| [01](01-introduction.md) | 课程介绍与技术选型 | Wish, Bubble Tea, Lipgloss 定位与关系 |
| [02](02-setup.md) | 环境准备与项目初始化 | go mod, 目录结构, go:embed |
| [03](03-data-layer.md) | 数据层：结构定义与 YAML 配置 | struct tags, yaml.v3, embed, fallback 模式 |
| [04](04-style-system.md) | 样式系统：复古终端主题 | Lipgloss 基础, 配色方案, 可复用样式 |
| [05](05-tui-basics.md) | TUI 基础：Bubble Tea 框架入门 | Elm 架构, Model/Update/View/Cmd |
| [06](06-app-model.md) | 应用模型：页面与导航 | 状态机, 输入处理, 页面切换逻辑 |
| [07](07-pages.md) | 页面渲染：各板块实现 | 字符串构建, 布局, 滚动高亮 |
| [08](08-ssh-server.md) | SSH 服务器：Wish 框架集成 | SSH 协议基础, Wish Middleware, teaHandler |
| [09](09-local-testing.md) | 本地测试与调试 | 本地入口, 单元测试, 常见问题 |
| [10](10-deployment.md) | 部署与运维 | Docker, 云服务器, Fly.io, 域名配置 |

## 练习题

- [练习 01](exercises/exercise-01.md)：修改配色主题
- [练习 02](exercises/exercise-02.md)：添加新页面（Education）
- [练习 03](exercises/exercise-03.md)：实现打字机动画效果

## 参考资源

- [Bubble Tea 官方文档](https://github.com/charmbracelet/bubbletea)
- [Lipgloss 官方文档](https://github.com/charmbracelet/lipgloss)
- [Wish 官方文档](https://github.com/charmbracelet/wish)
- [terminal.shop](https://terminal.shop) — 灵感来源
