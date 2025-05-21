package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(src, dst string) error {
	// 创建目标文件
	zipFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	// 创建 zip writer
	writer := zip.NewWriter(zipFile)
	defer writer.Close()

	// 获取源文件/目录信息
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to get source info: %w", err)
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(src)
	}

	// 遍历源目录
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk path %s: %w", path, err)
		}

		// 创建 zip header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header for %s: %w", path, err)
		}

		// 计算文件在 zip 中的路径
		if baseDir != "" {
			relPath, err := filepath.Rel(src, path)
			if err != nil {
				return fmt.Errorf("failed to get relative path: %w", err)
			}
			if relPath == "." {
				return nil
			}
			header.Name = filepath.Join(baseDir, relPath)
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		// 创建文件条目
		writer, err := writer.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		// 如果是目录，跳过写入内容
		if info.IsDir() {
			return nil
		}

		// 打开源文件
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", path, err)
		}
		defer file.Close()

		// 复制文件内容
		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to write file content: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create zip archive: %w", err)
	}

	return nil
}

func Unzip(src, dst string) error {
	// 打开 zip 文件
	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// 确保目标目录存在
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 遍历 zip 文件中的所有文件
	for _, file := range reader.File {
		// 构建目标文件路径
		filePath := filepath.Join(dst, file.Name)

		// 检查文件路径是否在目标目录内
		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", filePath)
		}

		// 如果是目录，创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// 确保父目录存在
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		// 打开源文件
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}

		// 创建目标文件
		f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create destination file: %w", err)
		}

		// 复制文件内容
		_, err = io.Copy(f, rc)
		f.Close()
		rc.Close()
		if err != nil {
			return fmt.Errorf("failed to copy file content: %w", err)
		}
	}

	return nil
}
