package parser

import (
	"strings"
	"testing"
)

func BenchmarkParseStandardFormat(b *testing.B) {
	parser := NewLogParser()
	input := "2024-01-01 12:00:00 [INFO] 这是一条标准格式的测试日志消息"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseISOFormat(b *testing.B) {
	parser := NewLogParser()
	input := "2024-01-01T12:00:00Z [WARN] 这是一条ISO格式的测试日志消息"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseSimpleFormat(b *testing.B) {
	parser := NewLogParser()
	input := "[DEBUG] 这是一条简单格式的测试日志消息"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseUnknownFormat(b *testing.B) {
	parser := NewLogParser()
	input := "这是一条未知格式的测试日志消息，应该被解析为UNKNOWN级别"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMixedFormats(b *testing.B) {
	parser := NewLogParser()
	inputs := []string{
		"2024-01-01 12:00:00 [INFO] 标准格式日志",
		"2024-01-01T12:00:00Z [ERROR] ISO格式日志",
		"[DEBUG] 简单格式日志",
		"未知格式日志",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_, err := parser.Parse(input)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkParseLongMessage(b *testing.B) {
	parser := NewLogParser()
	message := strings.Repeat("这是一个很长的测试日志消息，包含很多字符来测试性能 ", 100)
	input := "2024-01-01 12:00:00 [INFO] " + message

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parser.Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}
