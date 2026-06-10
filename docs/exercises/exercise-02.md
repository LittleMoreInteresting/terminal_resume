# 练习 02：添加 Education 页面

## 目标

在现有 5 个内容页基础上，添加第 6 个页面 **Education（教育经历）**。

## 任务

### 1. 更新数据结构

修改 `internal/data/resume.go`：

```go
// 添加 Education 结构
type Education struct {
    School     string `yaml:"school"`
    Degree     string `yaml:"degree"`
    Major      string `yaml:"major"`
    Period     string `yaml:"period"`
    Location   string `yaml:"location"`
    Highlights []string `yaml:"highlights"`
}

// Resume 结构添加字段
type Resume struct {
    // ... 现有字段
    Education []Education `yaml:"education"`
    // ... 现有字段
}
```

### 2. 更新 YAML

在 `internal/data/resume.yaml` 中添加：

```yaml
education:
  - school: "清华大学"
    degree: "硕士"
    major: "计算机科学与技术"
    period: "2016.09 - 2019.06"
    location: "北京"
    highlights:
      - "GPA: 3.8/4.0"
      - "发表 SCI 论文 2 篇"
  - school: "北京大学"
    degree: "本科"
    major: "软件工程"
    period: "2012.09 - 2016.06"
    location: "北京"
    highlights:
      - "优秀毕业生"
```

### 3. 更新样式包

修改 `internal/style/theme.go`：
- 在 `PageType` 枚举中添加 `EducationPage`
- 在 `AllPages()` 中插入新页面（建议放在 About 之后）
- 在 `PageTitles` 中添加映射

### 4. 更新模型

修改 `internal/app/model.go`：
- `Update()` 中添加数字键 `6` 的跳转
- `View()` 的 switch 中添加 `EducationPage` 分支
- 实现 `renderEducation()` 方法

### 5. 更新底部导航

修改 `renderFooter()` 中的页码显示（Education 可能是 `[2]` 或 `[6]`，取决于你的排序选择）。

## 提示

- `AllPages()` 的顺序决定了 ←/→ 导航的顺序
- 可以参考 `renderExperience()` 的实现方式
- 别忘了更新 Dashboard 的快速导航提示

## 验证

```bash
go run ./cmd/local
```

按 `6` 或导航到 Education 页面，确认显示正确。
