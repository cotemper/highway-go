package main

import (
	"flag"
	"fmt"

	"github.com/sonr-io/highway-go/internal/config"
	"github.com/sonr-io/highway-go/internal/handler"
	"github.com/sonr-io/highway-go/internal/svc"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/rest"
)
var configFile = flag.String("f", "etc/greet-api.yaml", "the config file")

func main() {
	// c, err := client.NewClient(context.Background(), "http://localhost:26657", "test", "bad-password")
	// if err != nil {
	// 	log.Fatal(err)
	// }


	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
