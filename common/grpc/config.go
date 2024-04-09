package grpc

import (
	"context"
	"encoding/json"
	"github.com/fidesy/sdk/common/logger"
	desc "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
)

type Value struct {
	value []byte
}

func (v *Value) MustBool() bool {
	var data bool
	err := json.Unmarshal(v.value, &data)
	if err != nil {
		panic(err)
	}

	return data
}

func GetValue(ctx context.Context, key string) Value {
	resp, err := realtimeConfigsServiceClient.GetValue(ctx, &desc.GetValueRequest{
		Key:         key,
		ServiceName: appName,
	})
	if err != nil {
		logger.Errorf("realtimeConfigsServiceClient.GetValue: %v", err)
		return Value{}
	}

	data, err := json.Marshal(resp.Value)
	if err != nil {
		logger.Errorf("json.Marshal: %v", err)
		return Value{}
	}

	return Value{value: data}
}
