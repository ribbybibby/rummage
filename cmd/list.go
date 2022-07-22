package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/ribbybibby/rummage/internal/image"
	"github.com/ribbybibby/rummage/internal/list"
	"github.com/spf13/cobra"
)

type imageOptions struct {
	Cache  bool
	Source string
}

func (o *imageOptions) AddFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&o.Cache, "cache", true, "cache images when using the remote source")
	cmd.Flags().StringVarP(&o.Source, "source", "s", "remote", "image source. Can be one of remote, daemon (Docker) or tarball")
}

func (o *imageOptions) Options() ([]image.Option, error) {
	opts := []image.Option{}

	if o.Source != "" {
		opts = append(opts, image.WithSource(o.Source))
	}

	if o.Cache {
		path, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		path = filepath.Join(path, "rummage")
		opts = append(opts, image.WithCache(path))
	}

	return opts, nil
}

type listOptions struct {
	Long  bool
	Image imageOptions
}

func (o *listOptions) AddFlags(cmd *cobra.Command) {
	o.Image.AddFlags(cmd)

	cmd.Flags().BoolVarP(&o.Long, "long", "l", false, "print extended file information, like ls -l or tar -tv")
}

func (o *listOptions) Options() []list.Option {
	var opts []list.Option

	if o.Long {
		opts = append(opts, list.WithLongOutput())
	}

	return opts
}

func List() *cobra.Command {
	o := &listOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files in a container image.",
		Long:  `List files in a container image.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			iOpts, err := o.Image.Options()
			if err != nil {
				return err
			}

			img, err := image.Get(ctx, args[0], iOpts...)
			if err != nil {
				return err
			}

			return list.List(ctx, img, os.Stdout, o.Options()...)
		},
	}

	o.AddFlags(cmd)

	return cmd
}

func init() {
	rootCmd.AddCommand(List())
}
