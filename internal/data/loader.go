package data

import (
    _ "embed"
    "fmt"

    "gopkg.in/yaml.v3"
)

//go:embed resume.yaml
var resumeYAML []byte

// LoadResume 从嵌入的 YAML 文件加载简历数据
func LoadResume() (*Resume, error) {
    var resume Resume
    if err := yaml.Unmarshal(resumeYAML, &resume); err != nil {
        return nil, fmt.Errorf("failed to parse resume YAML: %w", err)
    }
    return &resume, nil
}

// LoadResumeOrDefault 尝试加载 YAML，失败则返回默认数据
func LoadResumeOrDefault() *Resume {
    resume, err := LoadResume()
    if err != nil {
        return DefaultResume()
    }
    return resume
}
