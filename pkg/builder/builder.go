package builder

import "io"

type Builder interface {
	Build(out io.Writer, image, path string, env []string) error
}
