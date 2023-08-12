package domain

type FiberServer interface {
	Start()
}

type InfluxDb interface {
	WritePoint(measurement string, tags map[string]string, fields map[string]interface{}) error
	Close()
}
