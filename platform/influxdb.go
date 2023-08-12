package platform

import (
	"context"
	"time"

	"github.com/codern-org/codern/domain"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type influxDb struct {
	client   influxdb2.Client
	writeApi api.WriteAPIBlocking
	queryApi api.QueryAPI
}

func NewInfluxDb(url string, token string, org string, bucket string) (domain.InfluxDb, error) {
	client := influxdb2.NewClient(url, token)
	writeApi := client.WriteAPIBlocking(org, bucket)

	if ok, err := client.Ping(context.Background()); !ok {
		return nil, err
	}

	return &influxDb{
		client:   client,
		writeApi: writeApi,
		queryApi: client.QueryAPI(org),
	}, nil
}

func (db *influxDb) WritePoint(
	measurement string,
	tags map[string]string,
	fields map[string]interface{},
) error {
	point := influxdb2.NewPoint(measurement, tags, fields, time.Now())
	return db.writeApi.WritePoint(context.Background(), point)
}

func (db *influxDb) Close() {
	db.client.Close()
}
