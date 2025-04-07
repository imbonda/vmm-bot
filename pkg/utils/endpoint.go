package utils

import "path"

type Endpoint struct {
	BasePath string
}

func NewEndpoint(base string, elem ...string) *Endpoint {
	basePath := path.Join(base, path.Join(elem...))
	return &Endpoint{
		BasePath: basePath,
	}
}

func (e *Endpoint) Join(paths ...string) string {
	return path.Join(e.BasePath, path.Join(paths...))
}
