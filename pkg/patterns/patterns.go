package patterns

import (
	"bufio"
	"io"
	"regexp"

	"github.com/giantswarm/micrologger"
)

// FindMatch returns true if the given pattern is found in the input pipe.
func FindMatch(logger micrologger.Logger, input io.Reader, pattern string) (bool, error) {
	r, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		logger.Log("debug", "line to match: "+scanner.Text())
		if r.MatchString(scanner.Text()) {
			return true, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil
}
