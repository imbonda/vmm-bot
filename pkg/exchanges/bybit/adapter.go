package bybit

import (
	"fmt"
	"strings"
)

type requestInfo struct {
	method   string
	endpoint string
}

// This is a patch for bybit sdk as it does not return indicative errors.
func wrapResponseErrors(req *requestInfo, err error) error {
	if err == nil {
		return err
	}
	switch {
	// This error is thrown when a timeout occurs.
	// The sdk ignores errors and attempts to do "jsoniter" unmarshall on an empty bytes slice.
	case strings.Contains(err.Error(), "readObjectStart"):
		return fmt.Errorf("%s \"%s\": context deadline exceeded", req.method, req.endpoint)
	default:
		return err
	}
}
