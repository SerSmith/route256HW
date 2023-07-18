package tracer

import (
	"github.com/uber/jaeger-client-go/config"
)

func InitGlobal(service string) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	if _, err := cfg.InitGlobalTracer(service); err != nil {
		return err
	}
	return nil
}
