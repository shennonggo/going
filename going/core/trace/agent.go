package trace

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"

	"github.com/shennonggo/going/pkg/log"
)

/*
初始化不同的export的设置
*/

const (
	kindJaeger = "jaeger"
	kindZipkin = "zipkin"
)

var (
	//java中的set概念 运用到go中 , struct 空结构体不占内存， zerobase
	agents = make(map[string]struct{})
	lock   sync.Mutex
)

func InitAgent(o Options) {
	lock.Lock()
	defer lock.Unlock()
	//如果endpoint为空  就什么都不做
	_, ok := agents[o.Endpoint]
	if ok {
		return
	}
	err := startAgent(o)
	if err != nil {
		return
	}
	agents[o.Endpoint] = struct{}{}
}

func startAgent(o Options) error {
	var sexp trace.SpanExporter
	var err error

	opts := []trace.TracerProviderOption{
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(o.Sampler))),
		trace.WithResource(resource.NewSchemaless(semconv.ServiceNameKey.String(o.Name))),
	}
	//使用的是什么 jaeger？zipkin？
	if len(o.Endpoint) > 0 {
		switch o.Batcher {
		case kindJaeger:
			sexp, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(o.Endpoint)))
			if err != nil {
				return err
			}
		case kindZipkin:
			sexp, err = zipkin.New(o.Endpoint)
			if err != nil {
				return err
			}
		}
		opts = append(opts, trace.WithBatcher(sexp))
	}

	tp := trace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	//设置传播提取器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	//错误时打印日志
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		log.Errorf("[otel] error: %v", err)
	}))
	return nil
}
