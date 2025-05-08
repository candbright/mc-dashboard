package utils

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	// 增强正则表达式以匹配整个键值对
	kvRegex = regexp.MustCompile(`^\s*([\w\-_\s]+)\s*=\s*(.*?)\s*$`)

	// 错误类型
	ErrInvalidLineFormat = fmt.Errorf("invalid line format")
	ErrInvalidKey        = fmt.Errorf("invalid key format")
	ErrInvalidValue      = fmt.Errorf("invalid value format")
)

// TemplateConfig 模板配置
type TemplateConfig struct {
	// 是否保留注释
	KeepComments bool
	// 是否保留空行
	KeepEmptyLines bool
	// 是否保留原始缩进
	KeepIndentation bool
}

// ConvertToTemplate 将配置文件转换为模板
func ConvertToTemplate(inputContent string, config ...TemplateConfig) (string, error) {
	var cfg TemplateConfig
	if len(config) > 0 {
		cfg = config[0]
	}

	var outputBuilder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(inputContent))

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		processed, err := processLine(line, cfg)
		if err != nil {
			return "", fmt.Errorf("line %d: %w", lineNum, err)
		}

		if processed != "" {
			outputBuilder.WriteString(processed + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan error: %w", err)
	}

	return strings.TrimSuffix(outputBuilder.String(), "\n"), nil
}

// processLine 处理单行内容
func processLine(line string, config TemplateConfig) (string, error) {
	// 保留原始缩进
	leadingSpace := ""
	if config.KeepIndentation {
		leadingSpace = getLeadingSpace(line)
	}

	trimmedLine := strings.TrimLeftFunc(line, unicode.IsSpace)

	// 处理注释和空行
	if isCommentOrEmpty(trimmedLine) {
		if (config.KeepComments && strings.HasPrefix(trimmedLine, "#")) ||
			(config.KeepEmptyLines && trimmedLine == "") {
			return line, nil
		}
		return "", nil
	}

	// 匹配键值对
	matches := kvRegex.FindStringSubmatch(trimmedLine)
	if len(matches) < 3 {
		return "", ErrInvalidLineFormat
	}

	originalKey := matches[1]
	templateVar := toCamelCase(originalKey)

	// 验证键和值
	if err := validateKey(originalKey); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidKey, err)
	}

	if err := validateValue(matches[2]); err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidValue, err)
	}

	// 构建模板行
	return fmt.Sprintf("%s%s= {{.%s}}", leadingSpace, originalKey, templateVar), nil
}

// getLeadingSpace 获取行首的空格
func getLeadingSpace(line string) string {
	var space strings.Builder
	for _, r := range line {
		if unicode.IsSpace(r) {
			space.WriteRune(r)
		} else {
			break
		}
	}
	return space.String()
}

// isCommentOrEmpty 判断是否为注释或空行
func isCommentOrEmpty(line string) bool {
	trimmed := strings.TrimSpace(line)
	return len(trimmed) == 0 || strings.HasPrefix(trimmed, "#")
}

// toCamelCase 转换为小驼峰命名
func toCamelCase(s string) string {
	// 统一处理多种分隔符
	re := regexp.MustCompile(`[\s\-_]+`)
	words := re.Split(s, -1)

	// 处理数字开头的情况
	if len(words) > 0 && len(words[0]) > 0 {
		if unicode.IsDigit(rune(words[0][0])) {
			words[0] = "_" + words[0]
		}
	}

	var builder strings.Builder
	for i, word := range words {
		if len(word) == 0 {
			continue
		}

		// 统一转为小写处理
		lowerWord := strings.ToLower(word)
		if i > 0 {
			builder.WriteString(strings.Title(lowerWord))
		} else {
			// 首字母保持小写
			if len(lowerWord) > 0 {
				builder.WriteByte(lowerWord[0])
				if len(lowerWord) > 1 {
					builder.WriteString(lowerWord[1:])
				}
			}
		}
	}

	return sanitizeVarName(builder.String())
}

// validateKey 验证键的格式
func validateKey(key string) error {
	if len(key) == 0 {
		return fmt.Errorf("key cannot be empty")
	}

	// 检查是否包含非法字符
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return fmt.Errorf("key contains invalid character: %c", r)
		}
	}

	return nil
}

// validateValue 验证值的格式
func validateValue(value string) error {
	// 允许空值
	if len(value) == 0 {
		return nil
	}

	// 检查是否包含非法字符
	for _, r := range value {
		if !unicode.IsPrint(r) {
			return fmt.Errorf("value contains non-printable character")
		}
	}

	return nil
}

// sanitizeVarName 清理变量名
func sanitizeVarName(name string) string {
	var validChars []rune
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			validChars = append(validChars, r)
		}
	}
	return string(validChars)
}
