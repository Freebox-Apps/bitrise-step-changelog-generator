package git

import (
	"github.com/bitrise-io/go-utils/command"
)

// Log shows the commit logs. The format parameter controls what is shown and how.
func (g *Git) Log(format string, commitStart string, commitEnd string, otherOptions ...string) *command.Model {

	var options []string

	// Handle format
	if len(format) > 0 {
		options = append(options, "--format="+format)
	}

	// Handle commit range
	if len(commitStart) > 0 {
		if len(commitEnd) > 0 {
			options = append(options, commitStart + ".." + commitEnd)
		}else{
			options = append(options, commitStart + "..HEAD")
		}
	}

	// Handle other options
	if len(otherOptions) > 0 {
		options = append(options, otherOptions...)
	}

	log := append([]string{"log"}, options...)
	return g.command(log...)
}

func (g *Git) FetchTags() *command.Model {
	return g.command("fetch", "--tags")
}
