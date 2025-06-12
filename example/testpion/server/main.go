package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pion/mdns/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

type GRPCServiceInfo struct {
	InstanceName string
	Port         int
	TXTRecords   map[string]string
}

func startMDNSWithSRVRecords(service *GRPCServiceInfo) error {
	// Resolve mDNS addresses
	addr4, err := net.ResolveUDPAddr("udp4", mdns.DefaultAddressIPv4)
	if err != nil {
		return fmt.Errorf("failed to resolve IPv4 address: %v", err)
	}

	addr6, err := net.ResolveUDPAddr("udp6", mdns.DefaultAddressIPv6)
	if err != nil {
		return fmt.Errorf("failed to resolve IPv6 address: %v", err)
	}

	// Create UDP listeners
	l4, err := net.ListenUDP("udp4", addr4)
	if err != nil {
		return fmt.Errorf("failed to listen on IPv4: %v", err)
	}

	l6, err := net.ListenUDP("udp6", addr6)
	if err != nil {
		l4.Close()
		return fmt.Errorf("failed to listen on IPv6: %v", err)
	}

	// Create mDNS server configuration
	config := &mdns.Config{
		LocalNames:   []string{fmt.Sprintf("%s.local", service.InstanceName)},
		Name:         service.InstanceName,
		LocalAddress: net.ParseIP("192.168.1.201"),
	}

	// Start mDNS server
	_, err = mdns.Server(ipv4.NewPacketConn(l4), ipv6.NewPacketConn(l6), config)
	if err != nil {
		l4.Close()
		l6.Close()
		return fmt.Errorf("failed to start mDNS server: %v", err)
	}

	log.Printf("mDNS server started:")
	log.Printf("  Instance: %s.local", service.InstanceName)
	log.Printf("  Port: %d", service.Port)

	return nil
}

// FindNextAvailablePort finds the next available port starting from the given port
func FindNextAvailablePort(startPort int) (int, error) {
	for port := startPort; port < 65535; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found starting from %d", startPort)
}

func main() {
	// Generate unique service information
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Start the connect/gRPC server
	go func() {
		port, err := FindNextAvailablePort(50051)
		if err != nil {
			log.Fatalf("Failed to find available port: %v", err)
		}

		// Create service info for mDNS
		serviceInfo := &GRPCServiceInfo{
			InstanceName: fmt.Sprintf("grpc-service-%s", hostname),
			Port:         port,
			TXTRecords: map[string]string{
				"version": "1.0",
				"proto":   "grpc",
				"secure":  "false", // Set to "true" if using TLS
				"path":    "/",     // gRPC path if needed
			},
		}

		// Start mDNS server for this gRPC service
		if err := startMDNSWithSRVRecords(serviceInfo); err != nil {
			log.Printf("Failed to start mDNS server: %v", err)
		}

		log.Printf("Starting gRPC server on port %d", port)

		// Your existing gRPC server code
		// Replace svcCtx.mux with your actual mux
		var mux http.Handler // You'll need to replace this with your actual mux

		if err = http.ListenAndServe(
			fmt.Sprintf(":%d", port),
			h2c.NewHandler(mux, &http2.Server{}),
		); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("Server started. Press Ctrl+C to exit...")
	<-c
	log.Println("Shutting down server...")
}
