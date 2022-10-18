package main

import (
	"github.com/iktech/demo-service/internal/adapters/left/customer"
	"github.com/iktech/demo-service/internal/adapters/left/customers"
	"github.com/iktech/demo-service/internal/adapters/left/http"
	"github.com/iktech/demo-service/internal/adapters/left/version"
	"github.com/spf13/viper"
	"log"
	"net"
)

func main() {
	viper.SetDefault("server.bindAddress", "")
	viper.SetDefault("server.port", 5000)
	viper.SetConfigName("config")              // name of config file (without extension)
	viper.AddConfigPath("/etc/demo-service/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.demo-service") // call multiple times to add many search paths
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("application config not found, using defaults")
		} else {
			log.Fatalf("cannot read application config: %v", err)
		}
	}
	_ = viper.BindEnv("server.bindAddress", "DEMO_SERVER_BIND_ADDRESS")
	_ = viper.BindEnv("server.port", "DEMO_HTTP_PORT")

	ipAddress := viper.GetString("server.bindAddress")
	ip, err := net.ResolveIPAddr("ip", ipAddress)
	if err != nil {
		log.Fatalf("ip address %v is incorrect", ipAddress)
	}

	server := http.NewServerBuilder(ip, viper.GetInt("server.port")).
		AddRoute("GET", "/customers/:id", customer.NewHandler()).
		AddRoute("GET", "/version", version.NewHandler()).
		AddRoute("GET", "/customers", customers.NewHandler()).
		Build()
	err = server.Serve()
	if err != nil {
		log.Fatalf("cannot start server at %v:%d", ipAddress, viper.GetInt("server.port"))
	}
}
