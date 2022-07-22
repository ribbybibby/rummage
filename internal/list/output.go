package list

import (
	"archive/tar"
	"fmt"
	"io"
	"text/tabwriter"
)

// OutputFactory is a function that creates a new Output
type OutputFactory func(w io.Writer) Output

// Output writes the output for List
type Output interface {
	io.Closer
	Write(*tar.Header)
}

// ShortOutput writes the file name, separated by a newline.
func ShortOutput(w io.Writer) Output {
	return &shortOutput{w}
}

type shortOutput struct {
	w io.Writer
}

func (o *shortOutput) Write(header *tar.Header) {
	fmt.Fprintf(o.w, fmt.Sprintf("%s\n", header.Name))
}

func (o *shortOutput) Close() error {
	return nil
}

// LongOutput writes extended file information, similar to `tar -tv` or `ls -l`
func LongOutput(w io.Writer) Output {
	tw := tabwriter.NewWriter(w, 0, 1, 1, ' ', 0)
	return &longOutput{tw, tw.Flush}
}

type longOutput struct {
	w     io.Writer
	flush func() error
}

func (o *longOutput) Write(header *tar.Header) {
	if header.Typeflag == tar.TypeSymlink && header.Linkname != "" {
		fmt.Fprintf(o.w, "%s\t%d\t%d\t%d\t%s\t%s -> %s\n",
			header.FileInfo().Mode(),
			header.Uid,
			header.Gid,
			header.Size,
			header.ModTime.Format("2006-01-02 15:04:05"),
			header.Name,
			header.Linkname,
		)
	} else {
		fmt.Fprintf(o.w, "%s\t%d\t%d\t%d\t%s\t%s\n",
			header.FileInfo().Mode(),
			header.Uid,
			header.Gid,
			header.Size,
			header.ModTime.Format("2006-01-02 15:04:05"),
			header.Name,
		)
	}
}

func (o *longOutput) Close() error {
	return o.flush()
}
