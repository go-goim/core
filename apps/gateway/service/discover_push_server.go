package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/yusank/goim/apps/gateway/app"
)

func LoadMatchedPushServer(ctx context.Context) (string, error) {
	reg := app.GetRegister()
	if reg == nil {
		return "", fmt.Errorf("not init register")
	}

	// todo read service name from config
	list, err := reg.GetService(ctx, "goim.push.service")
	if err != nil {
		return "", err
	}

	if len(list) == 0 {
		return "", fmt.Errorf("service not found")
	}

	for _, instance := range list {
		if instance.Metadata["ready"] != "true" {
			continue
		}

		if len(instance.Endpoints) == 0 {
			continue
		}

		return instance.Endpoints[rand.Int()%len(instance.Endpoints)], nil
	}

	return "", fmt.Errorf("no matched service")
}
