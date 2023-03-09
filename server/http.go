package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fish/ai-tools/config"
	"github.com/fish/ai-tools/controllers"
)

// HTTPConfig http server config
type HTTPConfig struct {
	Name   string `yaml:"name"`
	Host   string `yaml:"host"`
	Port   uint16 `yaml:"port"`
	Mode   string `yaml:"mode"`
	RPCKey string `yaml:"rpc_key"`
}

// HTTPServer http server
type HTTPServer struct {
	host    string
	port    int
	handler *gin.Engine
	srv     *http.Server
}

// NewHTTPServer return a new http server
func NewHTTPServer(conf *config.SectionHTTP) *HTTPServer {
	gin.SetMode(conf.Mode)
	r := gin.Default()
	// FIXME: I don't understand why not effective!
	// r.Delims("[[", "]]")
	r.LoadHTMLGlob("views/**/*")
	return &HTTPServer{
		conf.Address,
		conf.Port,
		r,
		nil,
	}
}

// Run run a http server
func (s *HTTPServer) Run(mount controllers.Mount) {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	mount(s.handler)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}
	s.srv = srv
	go func() {
		log.Printf("Http Server is running on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

// Shutdown shutdown a http server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	err := s.srv.Shutdown(ctx)
	if err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Print("Http Server exiting\n")
	return err
}
