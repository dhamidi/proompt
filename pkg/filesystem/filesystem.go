package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
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

// Open implements fs.FS
func (rfs *RealFilesystem) Open(name string) (fs.File, error) {
	return rfs.readFS.Open(name)
}

// ReadFile implements fs.ReadFileFS
func (rfs *RealFilesystem) ReadFile(name string) ([]byte, error) {
	if filepath.IsAbs(name) {
		return os.ReadFile(name)
	}
	return fs.ReadFile(rfs.readFS, name)
}

// Stat implements fs.StatFS
func (rfs *RealFilesystem) Stat(name string) (fs.FileInfo, error) {
	if filepath.IsAbs(name) {
		return os.Stat(name)
	}
	return fs.Stat(rfs.readFS, name)
}

// ReadDir implements fs.ReadDirFS
func (rfs *RealFilesystem) ReadDir(name string) ([]fs.DirEntry, error) {
	if filepath.IsAbs(name) {
		return os.ReadDir(name)
	}
	return fs.ReadDir(rfs.readFS, name)
}

// WriteFile implements WriteFS
func (rfs *RealFilesystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// MkdirAll implements WriteFS
func (rfs *RealFilesystem) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Remove implements WriteFS
func (rfs *RealFilesystem) Remove(path string) error {
	return os.Remove(path)
}

// TempFile implements WriteFS
func (rfs *RealFilesystem) TempFile(dir, pattern string) (*os.File, error) {
	return os.CreateTemp(dir, pattern)
}

// Getwd implements Filesystem
func (rfs *RealFilesystem) Getwd() (string, error) {
	return os.Getwd()
}

// UserConfigDir implements Filesystem
func (rfs *RealFilesystem) UserConfigDir() (string, error) {
	return os.UserConfigDir()
}

// FakeFilesystem uses testing/fstest.MapFS + in-memory writes
type FakeFilesystem struct {
	fstest.MapFS
	cwd           string
	userConfigDir string
}

// NewFakeFilesystem creates a new FakeFilesystem
func NewFakeFilesystem() *FakeFilesystem {
	return &FakeFilesystem{
		MapFS:         make(fstest.MapFS),
		cwd:           "/",
		userConfigDir: "/home/user/.config",
	}
}

// WriteFile implements WriteFS for FakeFilesystem
func (ffs *FakeFilesystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	ffs.MapFS[path] = &fstest.MapFile{
		Data: data,
		Mode: perm,
	}
	return nil
}

// MkdirAll implements WriteFS for FakeFilesystem
func (ffs *FakeFilesystem) MkdirAll(path string, perm os.FileMode) error {
	// For FakeFilesystem, directories are implicitly created
	return nil
}

// Remove implements WriteFS for FakeFilesystem
func (ffs *FakeFilesystem) Remove(path string) error {
	if _, exists := ffs.MapFS[path]; !exists {
		return os.ErrNotExist
	}
	delete(ffs.MapFS, path)
	return nil
}

// TempFile implements WriteFS for FakeFilesystem
func (ffs *FakeFilesystem) TempFile(dir, pattern string) (*os.File, error) {
	// For testing, we'll simulate a temp file by creating a fake file
	tempPath := dir + "/" + pattern + "123456"
	ffs.MapFS[tempPath] = &fstest.MapFile{
		Data: []byte{},
		Mode: 0600,
	}
	// Return a nil file pointer since we can't create a real *os.File in memory
	// This is acceptable for testing purposes where the file path is more important
	return nil, nil
}

// Getwd implements Filesystem for FakeFilesystem
func (ffs *FakeFilesystem) Getwd() (string, error) {
	return ffs.cwd, nil
}

// UserConfigDir implements Filesystem for FakeFilesystem
func (ffs *FakeFilesystem) UserConfigDir() (string, error) {
	return ffs.userConfigDir, nil
}

// SetCwd sets the current working directory for FakeFilesystem
func (ffs *FakeFilesystem) SetCwd(cwd string) {
	ffs.cwd = cwd
}

// SetUserConfigDir sets the user config directory for FakeFilesystem
func (ffs *FakeFilesystem) SetUserConfigDir(dir string) {
	ffs.userConfigDir = dir
}
