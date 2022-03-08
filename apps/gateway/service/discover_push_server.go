package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/yusank/goim/pkg/registry"
)

func LoadMatchedPushServer(ctx context.Context) (string, error) {

	// todo read service name from config
	list, err := registry.GetService(ctx, "goim.push.service")
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
