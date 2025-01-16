package main

import (
	"context"
	"log"
	"math/rand/v2"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func userOrdersHandler(c *gin.Context) {
	ctx, span := tracer.Start(c.Request.Context(), "/users")
	defer span.End()

	span.SetAttributes(attribute.String("go.handler", "userOrdersHandler"))

	delay := time.Second * time.Duration(rand.IntN(3))
	time.Sleep(delay)

	_, span2 := tracer.Start(ctx, "cache")
	span2.SetStatus(codes.Error, "")
	time.Sleep(time.Second*time.Duration(rand.IntN(2)) + time.Second)
	span2.End()

	span.AddEvent("Processed user data", trace.WithTimestamp(time.Now()))

	c.String(200, "users")
}

func main() {
	tracerProvider, err := InitializeJaeger()
	if err != nil {
		log.Fatalf("Jaeger initialization error: %v", err)
	}
	defer func() {
		_ = tracerProvider.Shutdown(context.Background())
	}()

	tracer = otel.Tracer("user-app-tracer")
	r := gin.Default()

	r.Use(otelgin.Middleware("myservice"))

	r.GET("/users", userOrdersHandler)

	log.Fatal(r.Run(":3000"))
}
