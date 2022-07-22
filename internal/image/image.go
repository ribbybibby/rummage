package image

import (
	"context"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/cache"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

// Option is a functional option
type Option func(*options)

type options struct {
	Source string
	Cache  string
}

func makeOptions(opts ...Option) *options {
	o := &options{
		Source: "remote",
	}
	for _, opt := range opts {
		opt(o)
	}

	return o
}

// WithSource is a functional option that configures where images will be
// fetched from. Valid options are remote, daemon or tarball.
func WithSource(src string) Option {
	return func(o *options) {
		o.Source = src
	}
}

// WithCache is a functional option that caches layers to the provided path when
// using the remote source.
func WithCache(path string) Option {
	return func(o *options) {
		o.Cache = path
	}
}

// Get the image specified by imgRef
func Get(ctx context.Context, imgRef string, opts ...Option) (v1.Image, error) {
	o := makeOptions(opts...)

	switch o.Source {
	case "remote":
		return remoteImage(ctx, imgRef, o.Cache)
	case "daemon":
		return daemonImage(ctx, imgRef)
	case "tarball":
		return tarballImage(imgRef)
	default:
		return nil, fmt.Errorf("unsupported source: %s", o.Source)
	}
}

func remoteImage(ctx context.Context, imgRef, cachePath string) (v1.Image, error) {
	ref, err := name.ParseReference(imgRef)
	if err != nil {
		return nil, err
	}

	rOpts := []remote.Option{
		remote.WithContext(ctx),
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
	}
	desc, err := remote.Get(ref, rOpts...)
	if err != nil {
		return nil, err
	}

	img, err := desc.Image()
	if err != nil {
		return nil, err
	}

	if cachePath != "" {
		img = cache.Image(img, cache.NewFilesystemCache(cachePath))
	}

	return img, nil

}

func daemonImage(ctx context.Context, imgRef string) (v1.Image, error) {
	ref, err := name.ParseReference(imgRef)
	if err != nil {
		return nil, err
	}

	img, err := daemon.Image(ref, daemon.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func tarballImage(imgRef string) (v1.Image, error) {
	img, err := tarball.ImageFromPath(imgRef, nil)
	if err != nil {
		return nil, err
	}

	return img, nil
}
