package fs

import (
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hairyhenderson/go-fsimpl"
	"github.com/hairyhenderson/go-fsimpl/blobfs"
	"github.com/hairyhenderson/go-fsimpl/filefs"
	"github.com/hairyhenderson/go-fsimpl/gitfs"
	"github.com/hairyhenderson/go-fsimpl/httpfs"
)

var mux = sync.OnceValue(func() fsimpl.FSMux {
	mux := fsimpl.NewMux()
	mux.Add(filefs.FS)
	mux.Add(httpfs.FS)
	mux.Add(blobfs.FS)
	mux.Add(gitfs.FS)
	return mux
})

type Getter interface {
	Get(src string) (string, error)
	GetFile(src string) (string, error)
}

type gogetter struct {
	tmp string
}

func NewGoGetter(tmp string) gogetter {
	return gogetter{
		tmp: tmp,
	}
}

func (g gogetter) Get(src string) (string, error) {
	base, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	if base.Scheme == "" {
		return src, nil
	}
	fsys, err := mux().New(base)
	if err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(g.tmp, "*")
	if err != nil {
		return "", err
	}
	if err := os.CopyFS(tmp, fsys); err != nil {
		return "", err
	}
	return tmp, nil
}

func (g gogetter) GetFile(src string) (string, error) {
	base, err := url.Parse(src)
	if err != nil {
		return "", err
	}
	if base.Scheme == "" {
		return src, nil
	}
	repopath, subpath, _ := strings.Cut(base.Path, "//")
	base.Path = repopath
	fsys, err := mux().New(base)
	if err != nil {
		return "", err
	}
	tmp, err := os.MkdirTemp(g.tmp, "*")
	if err != nil {
		return "", err
	}
	data, err := fs.ReadFile(fsys, subpath)
	if err != nil {
		return "", err
	}
	file := filepath.Join(tmp, filepath.Base(subpath))
	if err := os.WriteFile(file, data, 0o600); err != nil {
		return "", err
	}
	return file, nil
}

type local struct{}

func NewLocal() local {
	return local{}
}

func (g local) Get(src string) (string, error) {
	if _, err := os.Stat(src); err != nil {
		return "", err
	}
	return src, nil
}

func (g local) GetFile(src string) (string, error) {
	return g.Get(src)
}
