package domain

import "io"

type FiberServer interface {
	Start()
}

type InfluxDb interface {
	WritePoint(measurement string, tags map[string]string, fields map[string]interface{}) error
	Close()
}

type SeaweedFs interface {
	Upload(content io.Reader, size int, path string) error
}
