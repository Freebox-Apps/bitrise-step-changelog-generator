package git

import (
	"os"

	"github.com/bitrise-io/go-utils/v2/command"
	"github.com/bitrise-io/go-utils/v2/env"
)

// Git represents a Git project.
type Git struct {
	dir string
}

// New creates a new git project.
func New(dir string) (Git, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return Git{}, err
	}
	return Git{dir: dir}, nil
}

func (g *Git) command(args ...string) command.Command {
	factory := command.NewFactory(env.NewRepository())
	opts := &command.Opts{
		Dir: g.dir,
		Env: []string{"GIT_ASKPASS=echo"},
	}
	return factory.Create("git", args, opts)
}
