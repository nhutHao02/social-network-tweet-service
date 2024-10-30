package startup

import (
	"github.com/nhutHao02/social-network-tweet-service/config"
	"github.com/nhutHao02/social-network-tweet-service/internal/api"
	"github.com/nhutHao02/social-network-tweet-service/internal/api/http"
)

func Start() {
	// load congig
	cfg := config.LoadConfig()

	// database setup

	// server
	http_server := http.NewHTTPServer(cfg)

	server := api.NewSerVer(http_server)
	runServer(server)

}

func runServer(server *api.Server) {

	// run http server
	server.HTTPServer.RunHTTPServer()

	// run grpc server
}
