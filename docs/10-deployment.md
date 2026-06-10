# 10 - 部署与运维

## 部署方案对比

| 方案 | 难度 | 成本 | 适合 |
|------|------|------|------|
| 云服务器 (VPS) | 中 | 低 | 长期运行，完全控制 |
| Fly.io | 低 | 免费额度 | 快速部署，自动 HTTPS |
| Docker | 中 | 取决于托管 | 标准化部署，易于迁移 |
| GitHub Actions | 低 | 免费 | CI/CD 自动化 |

## 方案 1：Docker 部署

### Dockerfile

```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o terminal-resume .

# 运行阶段
FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app

# 复制二进制文件和 SSH Host Key
COPY --from=builder /app/terminal-resume .
COPY --from=builder /app/.ssh ./.ssh

EXPOSE 23234
CMD ["./terminal-resume"]
```

**多阶段构建优势**：
- 构建阶段包含 Go 工具链（~1GB）
- 运行阶段只有 Alpine + 二进制（~20MB）
- 大幅减少攻击面和镜像体积

### 构建与运行

```bash
# 构建镜像
docker build -t terminal-resume .

# 运行容器
docker run -d \
  --name resume \
  -p 23234:23234 \
  -e HOST=0.0.0.0 \
  terminal-resume

# 测试连接
ssh localhost -p 23234
```

## 方案 2：Fly.io 部署（推荐）

Fly.io 提供免费的边缘部署，支持自定义域名。

### 安装 flyctl

```bash
# macOS
brew install flyctl

# Linux
curl -L https://fly.io/install.sh | sh

# Windows
pwsh -Command "iwr https://fly.io/install.ps1 -useb | iex"
```

### 部署步骤

```bash
# 登录
fly auth login

# 初始化应用（首次）
fly launch
# 按提示选择区域、应用名

# 部署
fly deploy

# 查看状态
fly status

# 查看日志
fly logs
```

### fly.toml 配置

```toml
app = 'your-app-name'
primary_region = 'sin'

[build]

[env]
  HOST = '0.0.0.0'
  PORT = '23234'

[[services]]
  internal_port = 23234
  protocol = 'tcp'

  [[services.ports]]
    port = 23234
    handlers = []
```

**注意**：Fly.io 的免费额度包括 3 个共享 CPU-1x 256MB VM，足够运行本项目。

## 方案 3：云服务器 (VPS)

### 服务器准备

```bash
# 购买 VPS 后（如阿里云、腾讯云、DigitalOcean）
# 通过 SSH 登录服务器
ssh root@your-server-ip

# 安装 Go（如果未安装）
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### 部署脚本

```bash
#!/bin/bash
# deploy.sh

APP_DIR="/opt/terminal-resume"
SERVICE_NAME="terminal-resume"

cd $APP_DIR

# 拉取最新代码
git pull origin main

# 构建
go build -o terminal-resume .

# 生成 Host Key（如果不存在）
if [ ! -f .ssh/term_info_ed25519 ]; then
    mkdir -p .ssh
    ssh-keygen -t ed25519 -f .ssh/term_info_ed25519 -N ""
fi

# 重启服务
sudo systemctl restart $SERVICE_NAME

echo "Deployed successfully!"
```

### Systemd 服务配置

创建 `/etc/systemd/system/terminal-resume.service`：

```ini
[Unit]
Description=Terminal Resume SSH Server
After=network.target

[Service]
Type=simple
User=resume
WorkingDirectory=/opt/terminal-resume
ExecStart=/opt/terminal-resume/terminal-resume
Restart=always
RestartSec=5
Environment="HOST=0.0.0.0"
Environment="PORT=23234"

[Install]
WantedBy=multi-user.target
```

启用服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable terminal-resume
sudo systemctl start terminal-resume
sudo systemctl status terminal-resume
```

## 域名配置

### SSH 默认端口（22）

如果希望用户使用 `ssh resume.yourdomain.com`（不加 `-p`），需要将服务绑定到 22 端口：

```bash
# 需要 root 权限，因为 22 是特权端口
sudo setcap cap_net_bind_service=+ep ./terminal-resume
# 或
sudo ./terminal-resume  # 不推荐
```

**更安全的做法**：使用反向代理或端口转发。

### Cloudflare + 域名

1. 购买域名（如 Namecheap、Cloudflare Registrar）
2. 添加 A 记录指向服务器 IP
3. 等待 DNS 传播（通常几分钟到几小时）

```
Type: A
Name: resume
Content: your-server-ip
TTL: Auto
```

连接测试：
```bash
ssh resume.yourdomain.com -p 23234
```

## 监控与日志

### 查看运行状态

```bash
# Systemd
sudo journalctl -u terminal-resume -f

# Docker
docker logs -f resume

# Fly.io
fly logs
```

### 健康检查

```bash
# 检查端口监听
ss -tlnp | grep 23234

# 测试连接
ssh -o ConnectTimeout=5 your-domain.com -p 23234
```

## 安全建议

1. **防火墙**：仅开放 SSH 端口，关闭其他端口
   ```bash
   sudo ufw allow 23234/tcp
   sudo ufw enable
   ```

2. **Rate Limiting**：防止暴力连接
   ```bash
   # 使用 fail2ban 或其他工具
   ```

3. **非 Root 运行**：Systemd 配置中的 `User=resume`

4. **Host Key 备份**：`.ssh/term_info_ed25519` 文件备份，更换服务器时保持一致（避免用户看到 host key 变更警告）

## 本章小结

- Docker 多阶段构建实现最小镜像
- Fly.io 适合快速免费部署
- VPS 提供完全控制和最低长期成本
- Systemd 确保服务崩溃后自动重启
- 域名 + 标准 22 端口提升用户体验

## 课程结束 🎉

恭喜完成本课程！你现在拥有一个可以通过 SSH 访问的交互式简历应用。

回顾所学：
1. ✅ Wish 构建 SSH 服务器
2. ✅ Bubble Tea 实现 TUI 交互
3. ✅ Lipgloss 设计复古终端样式
4. ✅ go:embed 嵌入配置文件
5. ✅ 双入口设计（SSH + 本地）
6. ✅ Docker + 云部署

## 扩展方向

- 添加动画效果（打字机、闪烁光标）
- 支持自定义主题（配置文件切换配色）
- 添加更多页面（Education、Certifications）
- 集成 analytics（记录访问日志）
- 支持多语言简历
