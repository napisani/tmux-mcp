package resource

import (
	"errors"
	"strings"
)

func extractParams(uri string) ([]string, error) {
	parts := strings.Split(uri, "://")
	if len(parts) < 2 {
		return nil, errors.New("URI has no parameter parts")
	}

	dest := make([]string, 0, len(parts)-1)
	copy(dest, parts[1:])
	return dest, nil
}
