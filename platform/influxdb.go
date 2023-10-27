package platform

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	http2 "github.com/influxdata/influxdb-client-go/v2/api/http"
	"go.uber.org/zap"
)

type InfluxDb struct {
	client   influxdb2.Client
	writeApi api.WriteAPI
	queryApi api.QueryAPI
}

func NewInfluxDb(
	url string,
	token string,
	org string,
	bucket string,
	logger *zap.Logger,
) (*InfluxDb, error) {
	client := influxdb2.NewClient(url, token)
	writeApi := client.WriteAPI(org, bucket)

	writeApi.SetWriteFailedCallback(func(batch string, err http2.Error, retryAttempts uint) bool {
		logger.Warn("InfluxDB write failed", zap.String("batch", batch), zap.Error(&err), zap.Int("retry", int(retryAttempts)))
		return true
	})

	if ok, err := client.Ping(context.Background()); !ok {
		return nil, err
	}

	return &InfluxDb{
		client:   client,
		writeApi: writeApi,
		queryApi: client.QueryAPI(org),
	}, nil
}

func (db *InfluxDb) WritePoint(
	measurement string,
	tags map[string]string,
	fields map[string]interface{},
) {
	point := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	db.writeApi.WritePoint(point)
}

func (db *InfluxDb) Close() {
	db.client.Close()
}
