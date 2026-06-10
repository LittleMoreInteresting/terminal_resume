# 练习 01：修改配色主题

## 目标

将应用的复古绿黑主题改为**琥珀色主题**（Amber on Black），类似老式 IBM 单色显示器。

## 参考配色

```
背景: #0a0a0a (Black)
主色: #ffb000 (Amber)
次色: #cc8800 (DimAmber)
暗色: #553300 (DarkAmber)
强调: #00ff41 (Green, 保留用于选中状态)
错误: #ff4444 (Red)
```

## 任务

1. 修改 `internal/style/theme.go` 中的颜色常量
2. 确保所有页面正确显示新配色
3. 运行本地版本验证效果

## 提示

- 只需要修改 `theme.go` 文件顶部的颜色定义
- 样式对象（如 `TitleStyle`、`BoxStyle`）使用颜色常量，会自动生效
- 建议保留 `SelectedItem` 的反色设计（深色背景 + 亮色文字）

## 进阶

尝试实现**蓝色主题**（Phosphor Blue）：
```
主色: #00bfff (DeepSkyBlue)
次色: #0066aa (DimBlue)
暗色: #002244 (DarkBlue)
```
