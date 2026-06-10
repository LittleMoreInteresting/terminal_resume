# 练习 03：实现打字机动画效果

## 目标

在 Dashboard 页面添加打字机效果：标题文字逐字显示，营造复古终端氛围。

## 效果示意

```
终端显示过程：
T          (0ms)
Te         (50ms)
Ter        (100ms)
Term       (150ms)
...
Terminal Resume  (完成)
```

## 任务

### 1. 添加打字机消息

在 `internal/app/model.go` 中添加：

```go
// 自定义消息
type typewriterMsg struct {
    text  string
    index int
}
```

### 2. 添加打字机命令

```go
func typewriterCmd(text string, index int) tea.Cmd {
    return tea.Tick(time.Millisecond*50, func(t time.Time) tea.Msg {
        return typewriterMsg{text: text, index: index}
    })
}
```

### 3. 修改 Model

```go
type Model struct {
    // ... 现有字段
    typewriterText  string  // 当前显示的打字机文字
    typewriterDone  bool    // 是否完成
}
```

### 4. 修改 Init

```go
func (m Model) Init() tea.Cmd {
    return typewriterCmd("Terminal Resume", 0)
}
```

### 5. 修改 Update

```go
case typewriterMsg:
    if msg.index < len(msg.text) {
        m.typewriterText = msg.text[:msg.index+1]
        return m, typewriterCmd(msg.text, msg.index+1)
    }
    m.typewriterDone = true
    return m, nil
```

### 6. 修改 Dashboard 渲染

在 `renderDashboard()` 中，用 `m.typewriterText` 替代静态标题的一部分。

## 进阶挑战

### 挑战 1：光标闪烁

在打字机文字末尾添加闪烁光标：

```go
// 使用 Lipgloss 的 Blink 样式
CursorStyle = lipgloss.NewStyle().
    Foreground(Green).
    Blink(true)

// 渲染时
if !m.typewriterDone {
    text += CursorStyle.Render("█")
}
```

### 挑战 2：多行打字机

让名片信息也逐行显示：
1. 先显示名字
2. 再显示职位
3. 再显示联系方式

### 挑战 3：跳过动画

按任意键跳过打字机动画，立即显示完整内容。

```go
case tea.KeyMsg:
    if !m.typewriterDone {
        m.typewriterDone = true
        m.typewriterText = "Terminal Resume"
        return m, nil
    }
    // ... 原有按键处理
```

## 提示

- `tea.Tick` 是 Bubble Tea 的定时器，返回 `Cmd`
- 每次 `Tick` 触发后发送消息，Update 处理后再发送下一个 `Tick`，形成链条
- 注意：打字机动画期间，其他按键应仍能响应（或按挑战 3 处理）

## 参考实现思路

```
Init() ──► typewriterCmd("T", 0)
              │
              ▼
Update(typewriterMsg{"T", 0}) ──► 显示 "T"
              │
              ▼
typewriterCmd("Te", 1) ──► Update(...) ──► 显示 "Te"
              │
              ▼
            ...
              │
              ▼
Update(typewriterMsg{"Terminal Resume", 14}) ──► 完成
```
