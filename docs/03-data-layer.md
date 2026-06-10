# 03 - 数据层：结构定义与 YAML 配置

## 设计目标

简历内容需要：
1. **可配置** — 用户只需修改 YAML 文件即可更新简历
2. **类型安全** — Go struct 提供编译时检查
3. **自包含** — 编译后二进制独立运行，不依赖外部文件

## 步骤 1：定义数据结构

创建 `internal/data/resume.go`：

```go
package data

// Resume 简历顶层结构
type Resume struct {
    Name       string       `yaml:"name"`
    Title      string       `yaml:"title"`
    Location   string       `yaml:"location"`
    Email      string       `yaml:"email"`
    GitHub     string       `yaml:"github"`
    Website    string       `yaml:"website"`
    About      About        `yaml:"about"`
    Experience []Experience `yaml:"experience"`
    Skills     Skills       `yaml:"skills"`
    Projects   []Project    `yaml:"projects"`
    Contact    Contact      `yaml:"contact"`
}

type About struct {
    Summary string   `yaml:"summary"`
    Details []string `yaml:"details"`
}

type Experience struct {
    Company    string   `yaml:"company"`
    Role       string   `yaml:"role"`
    Period     string   `yaml:"period"`
    Location   string   `yaml:"location"`
    Highlights []string `yaml:"highlights"`
}

type Skills struct {
    Languages  []string `yaml:"languages"`
    Frameworks []string `yaml:"frameworks"`
    Tools      []string `yaml:"tools"`
    Others     []string `yaml:"others"`
}

type Project struct {
    Name        string   `yaml:"name"`
    Description string   `yaml:"description"`
    TechStack   []string `yaml:"techStack"`
    URL         string   `yaml:"url"`
    Highlights  []string `yaml:"highlights"`
}

type Contact struct {
    Email    string `yaml:"email"`
    GitHub   string `yaml:"github"`
    Website  string `yaml:"website"`
    LinkedIn string `yaml:"linkedIn"`
    Twitter  string `yaml:"twitter"`
    Message  string `yaml:"message"`
}
```

### Struct Tag 详解

`` `yaml:"name"` `` 告诉 `yaml.v3` 在解析 YAML 时，将键 `name` 的值映射到此字段。

**驼峰命名问题**：
- YAML 键 `techStack` → Struct tag `yaml:"techStack"`
- YAML 键 `linkedIn` → Struct tag `yaml:"linkedIn"`
- 若 tag 写错，字段将保持零值，不会报错！这是常见 Bug 来源。

## 步骤 2：创建 YAML 配置文件

创建 `internal/data/resume.yaml`：

```yaml
name: "张三"
title: "全栈工程师"
location: "中国 · 北京"
email: "zhangsan@example.com"
github: "github.com/zhangsan"
website: "zhangsan.dev"

about:
  summary: "热爱技术的全栈开发者..."
  details:
    - "擅长 Go、Rust、TypeScript"
    - "熟悉 Kubernetes、Docker"

experience:
  - company: "字节跳动"
    role: "高级后端工程师"
    period: "2022.03 - 至今"
    location: "北京"
    highlights:
      - "负责核心业务系统架构设计"

skills:
  languages: ["Go", "Rust", "TypeScript"]
  frameworks: ["Gin", "React"]
  tools: ["Docker", "K8s"]
  others: ["Microservices"]

projects:
  - name: "Go-Cache"
    description: "高性能分布式缓存"
    techStack: ["Go", "Raft", "gRPC"]
    url: "github.com/zhangsan/go-cache"
    highlights:
      - "单机 QPS 达到 15 万"

contact:
  email: "zhangsan@example.com"
  github: "github.com/zhangsan"
  message: "欢迎联系我！"
```

## 步骤 3：实现 YAML 加载器

创建 `internal/data/loader.go`：

```go
package data

import (
    _ "embed"
    "fmt"
    "gopkg.in/yaml.v3"
)

//go:embed resume.yaml
var resumeYAML []byte

// LoadResume 从嵌入的 YAML 加载简历
func LoadResume() (*Resume, error) {
    var resume Resume
    if err := yaml.Unmarshal(resumeYAML, &resume); err != nil {
        return nil, fmt.Errorf("failed to parse resume YAML: %w", err)
    }
    return &resume, nil
}

// LoadResumeOrDefault 加载失败时返回默认数据
func LoadResumeOrDefault() *Resume {
    resume, err := LoadResume()
    if err != nil {
        return DefaultResume()
    }
    return resume
}
```

### Fallback 模式

`LoadResumeOrDefault()` 实现**优雅降级**：
- YAML 正常 → 使用 YAML 数据
- YAML 损坏/不存在 → 使用代码中的默认数据
- 保证应用始终能启动，不会崩溃

## 步骤 4：提供默认数据

在 `internal/data/resume.go` 中添加：

```go
// DefaultResume 返回默认简历数据
func DefaultResume() *Resume {
    return &Resume{
        Name:  "张三",
        Title: "全栈工程师",
        // ... 完整默认数据
    }
}
```

默认数据的作用：
1. YAML 加载失败时的后备
2. 展示数据结构的预期格式
3. 新用户直接编译即可看到效果

## 步骤 5：编写测试

创建 `internal/data/loader_test.go`：

```go
package data

import "testing"

func TestLoadResume(t *testing.T) {
    resume, err := LoadResume()
    if err != nil {
        t.Fatalf("LoadResume failed: %v", err)
    }

    if resume.Name == "" {
        t.Error("expected Name to be non-empty")
    }
    if len(resume.Experience) == 0 {
        t.Error("expected Experience to be non-empty")
    }

    t.Logf("Loaded resume for: %s", resume.Name)
}
```

运行测试：
```bash
go test ./internal/data/...
```

## 常见问题

### Q1: `expected ';', found embed`

原因：`//go:embed` 指令与 `import` 之间存在空行或其他内容。

解决：确保 `//go:embed` 紧贴在变量声明上方，无空行。

### Q2: `"embed" imported and not used`

原因：使用了 `import "embed"` 而非空白导入。

解决：改为 `import _ "embed"`，因为只需要注册 embed 功能，不直接使用包内符号。

### Q3: YAML 字段解析为空

原因：struct tag 与 YAML 键名不匹配（大小写、拼写）。

解决：仔细检查 tag 是否与 YAML 键完全一致，尤其是驼峰命名如 `techStack`。

## 本章小结

- 使用 struct + yaml tag 定义类型安全的简历数据结构
- `go:embed` 将 YAML 编译进二进制
- Fallback 模式确保应用鲁棒性
- 测试验证 YAML 解析正确性

## 下一步

[04 - 样式系统：复古终端主题](04-style-system.md)
