package metrics

type Repository interface {
	Println(v ...interface{})
}

type MetricServer struct {
	logger Repository
}

func New(logger Repository) *MetricServer {
	newMetricServer := MetricServer{
		logger: logger,
	}

	return &newMetricServer
}

func (m *MetricServer) Push(report string) error {
	m.logger.Println(report)

	return nil
}
