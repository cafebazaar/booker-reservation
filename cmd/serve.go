package cmd

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cafebazaar/booker-reservation/api"
	"github.com/cafebazaar/booker-reservation/common"
)

// serveCmd represents the serve command
var (
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Launches the grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			serveService()
		},
	}
	startTime time.Time
)

func init() {
	RootCmd.AddCommand(serveCmd)
}

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveService() {

	printVersion()

	kp, err := keyPair()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load the keyPair")
	}

	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewServerTLSFromCert(kp))}

	grpcServer := grpc.NewServer(opts...)

	api.RegisterServer(grpcServer)

	conn, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on %s", addr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		uptime := time.Now().Sub(startTime)
		w.Write([]byte(fmt.Sprintf("OK\nVersion: %s\nBuild Time: %s\nUptime: %s",
			common.Version, common.BuildTime, uptime)))
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: grpcHandlerFunc(grpcServer, mux),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*kp},
		},
	}

	startTime = time.Now()

	logrus.Infof("Starting at    %s", addr)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

	if err != nil {
		logrus.WithError(err).Fatalf("Stopped serving at     %s", addr)
	}
}
