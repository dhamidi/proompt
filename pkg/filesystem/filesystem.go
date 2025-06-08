package filesystem

import (
	"io/fs"
	"os"
	"testing/fstest"
)

// ReadFS combines standard fs interfaces for reading
type ReadFS interface {
	fs.ReadFileFS
	fs.StatFS
	fs.ReadDirFS
}

// WriteFS provides write operations
type WriteFS interface {
	WriteFile(path string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
	Remove(path string) error
	TempFile(dir, pattern string) (*os.File, error)
}

// Filesystem combines read and write operations
type Filesystem interface {
	ReadFS
	WriteFS
	Getwd() (string, error)
	UserConfigDir() (string, error)
}

// RealFilesystem wraps os.DirFS for reads + standard library for writes
type RealFilesystem struct {
	readFS fs.FS
}

// NewRealFilesystem creates a new RealFilesystem
func NewRealFilesystem(rootDir string) *RealFilesystem {
	return &RealFilesystem{
		readFS: os.DirFS(rootDir),
	}
}

// FakeFilesystem uses testing/fstest.MapFS + in-memory writes
type FakeFilesystem struct {
	fstest.MapFS
	cwd string
}

// NewFakeFilesystem creates a new FakeFilesystem
func NewFakeFilesystem() *FakeFilesystem {
	return &FakeFilesystem{
		MapFS: make(fstest.MapFS),
		cwd:   "/",
	}
}
