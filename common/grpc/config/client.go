package config

import (
	realtime_configs_service "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
	"os"
)

var appName = os.Getenv("APP_NAME")

var (
	realtimeConfigsServiceClient realtime_configs_service.RealtimeConfigsServiceClient
)

func Init(client realtime_configs_service.RealtimeConfigsServiceClient) {
	realtimeConfigsServiceClient = client
}
