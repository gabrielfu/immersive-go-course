package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	pb "github.com/CodeYourFuture/immersive-go-course/grpc-client-server/prober"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement prober.ProberServer.
type server struct {
	pb.UnimplementedProberServer
}

var opsLatency = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "prober_ops_latency_msecs",
	Help: "The latency of the probes in milliseconds.",
}, []string{"endpoint"})

func (s *server) DoProbes(ctx context.Context, in *pb.ProbeRequest) (*pb.ProbeReply, error) {
	numRequests := int(in.GetNumRequests())
	elapsedMsecs := float32(0)
	for i := 0; i < numRequests; i++ {
		start := time.Now()
		_, err := http.Get(in.GetEndpoint())
		if err != nil {
			return nil, err
		}
		elapsed := time.Since(start)
		elapsedMsecs += float32(elapsed / time.Millisecond)
	}
	elapsedMsecs = elapsedMsecs / float32(numRequests)
	opsLatency.With(prometheus.Labels{
		"endpoint": in.GetEndpoint(),
	}).Set(float64(elapsedMsecs))
	return &pb.ProbeReply{MeanLatencyMsecs: elapsedMsecs}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// add prometheus handler
	prometheus.MustRegister(opsLatency)
	httpServer := &http.Server{
		Handler: promhttp.Handler(),
		Addr:    fmt.Sprintf(":%d", 2112),
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	// grpc server
	s := grpc.NewServer()
	pb.RegisterProberServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
