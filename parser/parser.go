package parser

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// LogEntry 表示解析后的日志条目
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
	Fields    map[string]string
	Raw       string
}

// LogParser 日志解析器
type LogParser struct {
	patterns []logPattern
}

// logPattern 日志模式匹配器
type logPattern struct {
	name    string
	regex   *regexp.Regexp
	handler func([]string) (*LogEntry, error)
}

// NewLogParser 创建新的日志解析器
func NewLogParser() *LogParser {
	parser := &LogParser{}
	parser.initPatterns()
	return parser
}

// initPatterns 初始化支持的日志格式模式
func (p *LogParser) initPatterns() {
	// 标准日志格式: 2024-01-01 12:00:00 [INFO] 消息内容
	standardPattern := logPattern{
		name:  "standard",
		regex: regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}(?:\.\d+)?)\s+\[(\w+)\]\s+(.+)$`),
		handler: func(matches []string) (*LogEntry, error) {
			timestamp, err := time.Parse("2006-01-02 15:04:05", strings.TrimSpace(matches[1]))
			if err != nil {
				// 尝试带毫秒的格式
				timestamp, err = time.Parse("2006-01-02 15:04:05.000", strings.TrimSpace(matches[1]))
				if err != nil {
					return nil, fmt.Errorf("时间解析失败: %v", err)
				}
			}

			return &LogEntry{
				Timestamp: timestamp,
				Level:     matches[2],
				Message:   matches[3],
				Fields:    make(map[string]string),
				Raw:       matches[0],
			}, nil
		},
	}

	// ISO 8601 格式: 2024-01-01T12:00:00.000Z [INFO] 消息内容
	isoPattern := logPattern{
		name:  "iso8601",
		regex: regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?Z?)\s+\[(\w+)\]\s+(.+)$`),
		handler: func(matches []string) (*LogEntry, error) {
			timestamp, err := time.Parse(time.RFC3339, matches[1])
			if err != nil {
				// 尝试不带 Z 的格式
				timestamp, err = time.Parse("2006-01-02T15:04:05.000", matches[1])
				if err != nil {
					return nil, fmt.Errorf("时间解析失败: %v", err)
				}
			}

			return &LogEntry{
				Timestamp: timestamp,
				Level:     matches[2],
				Message:   matches[3],
				Fields:    make(map[string]string),
				Raw:       matches[0],
			}, nil
		},
	}

	// JSON 风格格式: {"timestamp": "2024-01-01T12:00:00Z", "level": "INFO", "message": "消息内容"}
	jsonPattern := logPattern{
		name:  "json",
		regex: regexp.MustCompile(`^\s*\{\s*"timestamp"\s*:\s*"([^"]+)"\s*,\s*"level"\s*:\s*"([^"]+)"\s*,\s*"message"\s*:\s*"([^"]+)"\s*\}\s*$`),
		handler: func(matches []string) (*LogEntry, error) {
			timestamp, err := time.Parse(time.RFC3339, matches[1])
			if err != nil {
				return nil, fmt.Errorf("时间解析失败: %v", err)
			}

			return &LogEntry{
				Timestamp: timestamp,
				Level:     matches[2],
				Message:   matches[3],
				Fields:    make(map[string]string),
				Raw:       matches[0],
			}, nil
		},
	}

	// 简单格式: [INFO] 消息内容
	simplePattern := logPattern{
		name:  "simple",
		regex: regexp.MustCompile(`^\[(\w+)\]\s+(.+)$`),
		handler: func(matches []string) (*LogEntry, error) {
			return &LogEntry{
				Timestamp: time.Now(),
				Level:     matches[1],
				Message:   matches[2],
				Fields:    make(map[string]string),
				Raw:       matches[0],
			}, nil
		},
	}

	p.patterns = []logPattern{standardPattern, isoPattern, jsonPattern, simplePattern}
}

// Parse 解析单行日志
func (p *LogParser) Parse(line string) (*LogEntry, error) {
	for _, pattern := range p.patterns {
		matches := pattern.regex.FindStringSubmatch(line)
		if matches != nil {
			return pattern.handler(matches)
		}
	}

	// 如果没有匹配的模式，返回原始行作为消息
	return &LogEntry{
		Timestamp: time.Now(),
		Level:     "UNKNOWN",
		Message:   line,
		Fields:    make(map[string]string),
		Raw:       line,
	}, nil
}

// String 返回格式化的日志条目字符串
func (e *LogEntry) String() string {
	timestamp := e.Timestamp.Format("2006-01-02 15:04:05.000")
	
	if len(e.Fields) > 0 {
		var fieldStrs []string
		for k, v := range e.Fields {
			fieldStrs = append(fieldStrs, fmt.Sprintf("%s=%s", k, v))
		}
		return fmt.Sprintf("[%s] %s: %s {%s}", timestamp, e.Level, e.Message, strings.Join(fieldStrs, ", "))
	}
	
	return fmt.Sprintf("[%s] %s: %s", timestamp, e.Level, e.Message)
}
