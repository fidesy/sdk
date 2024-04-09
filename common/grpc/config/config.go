package config

import (
	"context"
	"fmt"
	"github.com/fidesy/sdk/common/logger"
	desc "github.com/fidesy/sdk/services/realtime-configs-service/pkg/realtime-configs-service"
	"strconv"
	"time"
)

type Value struct {
	value string
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

	return Value{value: resp.Value}
}

func (v Value) String() string {
	return v.value
}

func (v Value) MustBool() bool {
	if v.value == "true" {
		return true
	}

	if v.value == "false" {
		return false
	}

	panic("type is not boolean")
}

func (v Value) MustDuration() time.Duration {
	duration, err := time.ParseDuration(v.value)
	if err != nil {
		panic(fmt.Errorf("time.ParseDuration: %v", err))
	}

	return duration
}

func (v Value) MustInt() int64 {
	value, err := strconv.ParseInt(v.value, 10, 64)
	if err != nil {
		panic(fmt.Errorf("strconv.ParseInt: %v", err))
	}

	return value
}
