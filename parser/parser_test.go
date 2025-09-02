package parser

import (
	"testing"
	"time"
)

func TestLogParser(t *testing.T) {
	parser := NewLogParser()

	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*LogEntry) bool
	}{
		{
			name:    "标准格式",
			input:   "2024-01-01 12:00:00 [INFO] 这是一条测试日志",
			wantErr: false,
			check: func(entry *LogEntry) bool {
				return entry.Level == "INFO" && entry.Message == "这是一条测试日志"
			},
		},
		{
			name:    "带毫秒的标准格式",
			input:   "2024-01-01 12:00:00.123 [ERROR] 错误信息",
			wantErr: false,
			check: func(entry *LogEntry) bool {
				return entry.Level == "ERROR" && entry.Message == "错误信息"
			},
		},
		{
			name:    "ISO 8601 格式",
			input:   "2024-01-01T12:00:00Z [WARN] 警告信息",
			wantErr: false,
			check: func(entry *LogEntry) bool {
				return entry.Level == "WARN" && entry.Message == "警告信息"
			},
		},
		{
			name:    "简单格式",
			input:   "[DEBUG] 调试信息",
			wantErr: false,
			check: func(entry *LogEntry) bool {
				return entry.Level == "DEBUG" && entry.Message == "调试信息"
			},
		},
		{
			name:    "未知格式",
			input:   "这是一条未知格式的日志",
			wantErr: false,
			check: func(entry *LogEntry) bool {
				return entry.Level == "UNKNOWN" && entry.Message == "这是一条未知格式的日志"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := parser.Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.check(entry) {
				t.Errorf("Parse() result validation failed for input: %s", tt.input)
			}
		})
	}
}

func TestLogEntryString(t *testing.T) {
	entry := &LogEntry{
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Level:     "INFO",
		Message:   "测试消息",
		Fields:    make(map[string]string),
		Raw:       "原始日志",
	}

	result := entry.String()
	expected := "[2024-01-01 12:00:00.000] INFO: 测试消息"
	if result != expected {
		t.Errorf("String() = %v, want %v", result, expected)
	}
}
