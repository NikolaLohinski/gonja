package loaders

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// filesystemLoader represents a local filesystem loader with basic
// BaseDirectory capabilities. The access to the local filesystem is unrestricted.
type filesystemLoader struct {
	root string
}

// MustNewFileSystemLoader creates a new FilesystemLoader instance
// and panics if there's any error during instantiation. The parameters
// are the same like NewFileSystemLoader.
func MustNewFileSystemLoader(root string) Loader {
	fs, err := NewFileSystemLoader(root)
	if err != nil {
		log.Panic(err)
	}
	return fs
}

// NewFileSystemLoader creates a new FilesystemLoader and allows
// templates to be loaded from disk (unrestricted). If any base directory
// is given (or being set using SetBaseDir), this base directory is being used
// for path calculation in template inclusions/imports. Otherwise the path
// is calculated relatively to the current working directory.
func NewFileSystemLoader(root string) (Loader, error) {
	fs := &filesystemLoader{}
	if root != "" {
		// Make the path absolute
		if !filepath.IsAbs(root) {
			abs, err := filepath.Abs(root)
			if err != nil {
				return nil, err
			}
			root = abs
		}

		// Check for existence
		fi, err := os.Stat(root)
		if err != nil {
			return nil, err
		}
		if !fi.IsDir() {
			return nil, errors.Errorf("The given root '%s' is not a directory.", root)
		}

		fs.root = root
	}
	return fs, nil
}

func (fs *filesystemLoader) Inherit(root string) (Loader, error) {
	if root == "" {
		root = fs.root
	}
	return NewFileSystemLoader(root)
}

// Get reads the path's content from your local filesystem.
func (fs *filesystemLoader) Read(path string) (io.Reader, error) {
	realPath, err := fs.Resolve(path)
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(realPath)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}

// Path resolves a filename relative to the base directory. Absolute paths are allowed.
// When there's no base dir set, the absolute path to the filename
// will be calculated based on either the provided base directory (which
// might be a path of a template which includes another template) or
// the current working directory.
func (fs *filesystemLoader) Resolve(name string) (string, error) {
	if filepath.IsAbs(name) {
		return name, nil
	}

	// root := fs.root
	if fs.root == "" {
		root, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return filepath.Join(root, name), nil
	} else {
		return filepath.Join(fs.root, name), nil
	}
}
