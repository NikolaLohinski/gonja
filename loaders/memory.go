package loaders

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// MemoryLoader represents a basic in-memory loader
type memoryLoader struct {
	root    string
	content map[string]string
}

func MustNewMemoryLoader(content map[string]string) Loader {
	loader, err := NewMemoryLoader(content)
	if err != nil {
		log.Panic(err)
	}
	return loader
}

// NewMemoryLoader creates a new MemoryLoader and allows
// templates to be loaded from memory.
func NewMemoryLoader(content map[string]string) (Loader, error) {
	root := ""
	for key := range content {
		if !strings.HasPrefix(key, string(os.PathSeparator)) {
			return nil, fmt.Errorf("all keys must start with '%s' but the following does not: '%s'", string(os.PathSeparator), key)
		}
		if root == "" {
			root = key
			continue
		}
		for !strings.HasPrefix(key, root) && root != "" {
			root = root[:len(root)-1]
		}
	}
	return &memoryLoader{
		root:    root,
		content: content,
	}, nil
}

func (m *memoryLoader) Inherit(root string) (Loader, error) {
	if root == "" {
		root = m.root
	}
	return &memoryLoader{
		content: m.content,
		root:    root,
	}, nil
}

func (m *memoryLoader) Read(name string) (io.Reader, error) {
	resolved, err := m.Resolve(name)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve name '%s': %s", name, err)
	}

	data, ok := m.content[resolved]
	if !ok {
		return nil, fmt.Errorf("unknown path: '%s'", resolved)
	}
	return strings.NewReader(data), nil
}

func (m *memoryLoader) Resolve(name string) (string, error) {
	if strings.HasPrefix(name, string(os.PathSeparator)) {
		return name, nil
	}
	resolved := filepath.Join(m.root, name)
	if _, ok := m.content[resolved]; !ok {
		return "", fmt.Errorf("unknown resolved path: '%s'", resolved)
	}

	return resolved, nil
}
