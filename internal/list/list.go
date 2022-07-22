package list

import (
	"archive/tar"
	"context"
	"errors"
	"io"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

// List files in a container image and write them to the provided io.Writer
func List(ctx context.Context, img v1.Image, w io.Writer, opts ...Option) error {
	o := makeOptions(opts...)

	rc := mutate.Extract(img)
	defer rc.Close()

	tr := tar.NewReader(rc)

	out := o.OutputFactory(w)
	defer out.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		header, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}

		out.Write(header)
	}

	return nil
}

// Option is a functional option for configuring the behaviour of List
type Option func(*options)

type options struct {
	OutputFactory OutputFactory
}

func makeOptions(opts ...Option) *options {
	o := &options{
		OutputFactory: ShortOutput,
	}
	for _, option := range opts {
		option(o)
	}

	return o
}

// WithLongOutput is a function option that logs extended file information
func WithLongOutput() Option {
	return func(o *options) {
		o.OutputFactory = LongOutput
	}
}
