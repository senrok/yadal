package fs

import (
	"context"
	"github.com/senrok/yadal/errors"
	"github.com/senrok/yadal/interfaces"
	"github.com/senrok/yadal/object"
	"github.com/senrok/yadal/options"
	"github.com/senrok/yadal/providers"
	"github.com/senrok/yadal/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Driver struct {
	root string
}

func (d Driver) fsMetadata(absPath string) (os.FileInfo, error) {
	info, err := os.Stat(absPath)
	if info.IsDir() != strings.HasSuffix(absPath, "/") {
		return nil, os.ErrNotExist
	}
	return info, err
}

func (d Driver) Metadata() interfaces.Metadata {
	return providers.NewMetadata(
		interfaces.Fs,
		d.root,
		"",
		interfaces.Read|interfaces.Write|interfaces.List,
	)
}

func (d Driver) Create(ctx context.Context, path string, args options.CreateOptions) error {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return errors.ParseFsError(errors.ErrCreateFailed, err, path)
	}
	parent := filepath.Dir(p)
	err = os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return errors.ParseFsError(errors.ErrCreateFailed, err, path)
	}
	switch interfaces.ObjectMode(args.Mode) {
	case interfaces.DIR:
		err = os.MkdirAll(p, os.ModePerm)
		if err != nil {
			return errors.ParseFsError(errors.ErrCreateFailed, err, p)
		}
	case interfaces.FILE:
		var file *os.File
		file, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE, os.ModePerm)
		defer file.Close()
		if err != nil {
			return errors.ParseFsError(errors.ErrCreateFailed, err, p)
		}
	}
	return nil
}

func (d Driver) Read(ctx context.Context, path string, args options.ReadOptions) (io.ReadCloser, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, errors.ParseFsError(errors.ErrReadFailed, err, path)
	}
	file, err := os.OpenFile(p, os.O_RDONLY, os.ModePerm)
	if err != nil {
		_ = file.Close()
		return nil, errors.ParseFsError(errors.ErrReadFailed, err, path)
	}
	if args.Offset != nil {
		_, err = file.Seek(int64(*args.Offset), 0)
		if err != nil {
			_ = file.Close()
			return nil, errors.ParseFsError(errors.ErrReadFailed, err, path)
		}
	}
	if args.Size != nil {
		return utils.NewFileLimitReader(file, int64(*args.Size)), nil
	}
	return file, nil
}

func (d Driver) Write(ctx context.Context, path string, args options.WriteOptions, reader io.ReadSeeker) (uint64, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return 0, errors.ParseFsError(errors.ErrWriteFailed, err, path)
	}
	parent := filepath.Dir(p)
	err = os.MkdirAll(parent, os.ModePerm)
	if err != nil {
		return 0, errors.ParseFsError(errors.ErrWriteFailed, err, path)
	}
	var file *os.File
	file, err = os.OpenFile(p, os.O_RDONLY|os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.ModePerm)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		return 0, errors.ParseFsError(errors.ErrWriteFailed, err, path)
	}
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return 0, errors.ParseFsError(errors.ErrWriteFailed, err, path)
	}
	written, err := file.Write(bytes)
	if err != nil {
		return 0, errors.ParseFsError(errors.ErrWriteFailed, err, path)
	}
	return uint64(written), nil
}

func (d Driver) Stat(ctx context.Context, path string, args options.StatOptions) (interfaces.ObjectMetadata, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, errors.ParseFsError(errors.ErrStatFailed, err, p)
	}

	info, err := os.Stat(p)
	if err != nil {
		return nil, errors.ParseFsError(errors.ErrStatFailed, err, p)
	}
	return object.NewMetadata(object.SetFromFileInfo(info))
}

func (d Driver) Delete(ctx context.Context, path string, args options.DeleteOptions) error {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return errors.ParseFsError(errors.ErrDeleteFailed, err, path)
	}
	meta, err := d.fsMetadata(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return errors.ParseFsError(errors.ErrDeleteFailed, err, path)
	}
	if meta.IsDir() {
		err = os.RemoveAll(p)
		if err != nil {
			return errors.ParseFsError(errors.ErrDeleteFailed, err, path)
		}
	} else {
		err = os.Remove(p)
		if err != nil {
			return errors.ParseFsError(errors.ErrDeleteFailed, err, path)
		}
	}
	return nil
}

func (d Driver) List(ctx context.Context, path string, args options.ListOptions) (interfaces.ObjectStream, error) {
	p, err := utils.BuildAbsPath(d.root, path)
	if err != nil {
		return nil, errors.ParseFsError(errors.ErrListFailed, err, path)
	}
	list, err := os.ReadDir(p)
	if err != nil {
		return nil, errors.ParseFsError(errors.ErrListFailed, err, path)
	}
	return &DirStream{
		Driver:  &d,
		root:    d.root,
		path:    path,
		entries: list,
	}, nil
}

func (d Driver) PreSign(ctx context.Context, path string, args options.PreSignOptions) (*http.Request, error) {
	return nil, errors.ErrUnsupportedMethod
}

func (d Driver) CreateMultipart(ctx context.Context, path string, args options.CreateMultipart) (string, error) {
	return "", errors.ErrUnsupportedMethod
}

func (d Driver) WriteMultipart(ctx context.Context, path string, args options.WriteMultipart, reader io.ReadSeeker) (interfaces.ObjectPart, error) {
	return nil, errors.ErrUnsupportedMethod
}

func (d Driver) CompleteMultipart(ctx context.Context, path string, args options.CompleteMultipart) error {
	return errors.ErrUnsupportedMethod
}

func (d Driver) AbortMultipart(ctx context.Context, path string, args options.AbortMultipart) error {
	return errors.ErrUnsupportedMethod
}

type Options struct {
	Root string
}

func NewDriver(opt Options) interfaces.Accessor {
	return &Driver{root: utils.NormalizeRoot(opt.Root)}
}
