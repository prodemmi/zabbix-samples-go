package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func sendReq(ctx context.Context, method, url string, data *[]byte) (*http.Response, error) {
	var body *bytes.Buffer
	if method == "POST" && data != nil {
		body = bytes.NewBuffer(*data)
	} else {
		body = bytes.NewBuffer([]byte{})
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	log.Printf("%v", req.Header)

	return client.Do(req)
}

func usersHandler(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "/users")
	defer span.End()

	span.SetAttributes(attribute.String("go.handler", "usersHandler"))

	res, err := sendReq(ctx, "GET", "http://user-app:3000/users", nil)
	if err != nil {
		span.SetStatus(codes.Error, "Failed to fetch users")
		span.RecordError(err)
		c.String(500, err.Error())
		return
	}

	if res != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			span.RecordError(err)
			c.String(500, err.Error())
			return
		}
		c.String(200, string(body))
	}
	defer res.Body.Close()

	span.AddEvent("User data sent to client")
}

func main() {
	tracerProvider, err := InitializeJaeger()
	if err != nil {
		log.Fatalf("Jaeger initialization error: %v", err)
	}
	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
	}()

	tracer = otel.Tracer("api-gateway-tracer")
	r := gin.Default()

	r.Use(otelgin.Middleware("myservice"))

	r.GET("/users", usersHandler)

	log.Fatal(r.Run(":3000"))
}
