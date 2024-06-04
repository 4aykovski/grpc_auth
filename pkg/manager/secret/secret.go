package secret

import (
	"context"
	"fmt"
	"os"
)

type Manager struct{}

func (m *Manager) GetSecret(ctx context.Context, appID int) (string, error) {
	secret := os.Getenv(fmt.Sprintf("APP%d_SECRET", appID))
	if secret == "" {
		return "", fmt.Errorf("secret not found")
	}

	return secret, nil
}
