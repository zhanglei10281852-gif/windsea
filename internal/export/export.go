package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/zhanglei10281852-gif/windsea/internal/domain"
	"io"
	"path/filepath"
	"time"
)

type Artifact struct {
	Name    string
	Content []byte
}

func BuildBundle(artifacts []Artifact) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	writer := zip.NewWriter(buffer)
	seen := map[string]bool{}
	for _, artifact := range artifacts {
		if artifact.Name == "" || seen[artifact.Name] {
			return nil, domain.ErrValidation
		}
		seen[artifact.Name] = true
		header := &zip.FileHeader{Name: filepath.Clean(artifact.Name), Method: zip.Deflate, Modified: time.Now()}
		entry, err := writer.CreateHeader(header)
		if err != nil {
			return nil, err
		}
		if _, err := entry.Write(artifact.Content); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func ReadBundle(raw []byte) (map[string][]byte, error) {
	reader, err := zip.NewReader(bytes.NewReader(raw), int64(len(raw)))
	if err != nil {
		return nil, err
	}
	result := map[string][]byte{}
	for _, file := range reader.File {
		entry, err := file.Open()
		if err != nil {
			return nil, err
		}
		content, err := io.ReadAll(entry)
		_ = entry.Close()
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", file.Name, err)
		}
		result[file.Name] = content
	}
	return result, nil
}
