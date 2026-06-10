# 08 - SSH 服务器：Wish 框架集成

## SSH 协议基础

SSH（Secure Shell）通常用于远程登录服务器执行命令。但 SSH 协议支持 **PTY（伪终端）分配**，允许服务器向客户端发送交互式界面。

传统 SSH 服务器流程：
```
用户 ──ssh──► SSHD ──exec──► bash/shell
```

Wish 的 SSH 服务器流程：
```
用户 ──ssh──► Wish ──Bubble Tea──► TUI 界面
```

## Wish 核心概念

### Middleware 模式

Wish 使用中间件（Middleware）处理 SSH 会话。每个中间件可以对会话进行加工：

```go
wish.WithMiddleware(
    logging.Middleware(),      // 记录日志
    activeterm.Middleware(),   // 强制要求交互式终端
    bubbletea.Middleware(teaHandler),  // 将会话转为 TUI
)
```

中间件执行顺序：**从外到内**，后添加的中间件先处理原始会话。

### teaHandler 签名

```go
// bubbletea.Middleware 需要这个签名的函数
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption)
```

- 参数 `s ssh.Session`：当前 SSH 会话对象
- 返回值 `tea.Model`：Bubble Tea 模型
- 返回值 `[]tea.ProgramOption`：程序选项（如 AltScreen）

## 实现 SSH 服务器

创建根目录 `main.go`：

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/activeterm"
    "github.com/charmbracelet/wish/bubbletea"
    "github.com/charmbracelet/wish/logging"
    tea "github.com/charmbracelet/bubbletea"
    "terminal_resume/internal/app"
)

func main() {
    // 从环境变量读取配置，提供默认值
    port := os.Getenv("PORT")
    if port == "" {
        port = "23234"
    }
    host := os.Getenv("HOST")
    if host == "" {
        host = "localhost"
    }

    // 创建 Wish 服务器
    server, err := wish.NewServer(
        wish.WithAddress(fmt.Sprintf("%s:%s", host, port)),
        wish.WithHostKeyPath(".ssh/term_info_ed25519"),
        wish.WithMiddleware(
            logging.Middleware(),
            activeterm.Middleware(),
            bubbletea.Middleware(teaHandler),
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    // 优雅关闭：监听系统信号
    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    log.Printf("Starting SSH server on %s:%s", host, port)

    go func() {
        if err := server.ListenAndServe(); err != nil {
            log.Fatalf("Server error: %v", err)
        }
    }()

    <-done
    log.Println("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Shutdown error: %v", err)
    }

    log.Println("Server stopped")
}

// teaHandler：每个 SSH 会话创建独立的 Bubble Tea 程序
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    m := app.NewModel()

    return m, []tea.ProgramOption{
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    }
}
```

## 关键配置详解

### Host Key

SSH 服务器需要 Host Key 用于身份验证（不是用户认证，是服务器身份验证）：

```bash
# 生成 ED25519 密钥（推荐）
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N ""

# 生成两个文件：
# .ssh/term_info_ed25519      (私钥，保密)
# .ssh/term_info_ed25519.pub  (公钥，可公开)
```

**为什么用 ED25519？**
- 比 RSA 更安全
- 密钥更短（256 bit 等效 RSA 3000+ bit）
- 生成和验证更快

### activeterm.Middleware()

```go
wish.WithMiddleware(
    activeterm.Middleware(),  // 拒绝非交互式会话
)
```

这个中间件会拒绝没有请求 PTY 的 SSH 连接（如 `ssh host command` 的脚本调用），确保用户确实在交互式终端中。

### logging.Middleware()

自动记录每个 SSH 会话的 connect/disconnect 事件，便于监控和调试。

## 运行与测试

### 1. 生成 Host Key

```bash
mkdir -p .ssh
ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N ""
```

### 2. 启动服务器

```bash
go run .
```

输出：
```
2024/xx/xx xx:xx:xx Starting SSH server on localhost:23234
2024/xx/xx xx:xx:xx Connect with: ssh localhost -p 23234
```

### 3. 客户端连接

```bash
# 本地测试
ssh localhost -p 23234

# 首次连接会提示确认 host key fingerprint，输入 yes
```

### 4. 退出

在 TUI 中按 `q` 或 `Ctrl+C` 退出，SSH 会话自动关闭。

## 常见问题

### Q1: `Failed to create server: open .ssh/term_info_ed25519: no such file`

原因：Host Key 文件不存在。

解决：运行 `ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N ""` 生成。

### Q2: 连接后显示空白或乱码

原因：客户端终端不支持某些 ANSI 转义码，或字体缺少某些字符。

解决：
- 使用现代终端（Windows Terminal、iTerm2、GNOME Terminal）
- 确保使用支持 Unicode 的字体（如 JetBrains Mono、Fira Code）
- 检查 `LANG` 环境变量是否包含 UTF-8

### Q3: 如何配置公钥认证？

Wish 默认接受任何连接（无密码/密钥认证）。如需添加认证：

```go
import "github.com/charmbracelet/wish/git"

wish.WithMiddleware(
    git.AccessMiddleware(git.AccessParams{
        PublicKeyPath: ".ssh/authorized_keys",
    }),
    // ... 其他中间件
)
```

但本项目的简历应用通常不需要认证，保持开放访问即可。

## 本章小结

- Wish 将 SSH 会话包装为 Bubble Tea 程序
- `teaHandler` 为每个会话创建独立的 Model
- Host Key 用于服务器身份验证，必须预先生成
- 优雅关闭通过信号监听 + `context.WithTimeout` 实现
- `activeterm.Middleware()` 确保交互式终端

## 下一步

[09 - 本地测试与调试](09-local-testing.md)
