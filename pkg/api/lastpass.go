package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"

	"github.com/ansd/lastpass-go"
)

type Error string

const (
	ErrBadCreds = Error("invalid username or password")
	ErrBadOTP   = Error("invalid OTP")
)

func (e Error) Error() string { return string(e) }

func Load(ctx context.Context, username, password, otp string) ([]store.Entry, error) {
	client, err := lastpass.NewClient(context.Background(), username, password, lastpass.WithOneTimePassword(otp))
	if err != nil {
		if authErr, ok := err.(*lastpass.AuthenticationError); ok {
			if strings.Contains(authErr.Error(), "password_invalid") {
				return nil, ErrBadCreds
			}
		}
		return nil, fmt.Errorf("failed to connect to lastpass API: %w", err)
	}

	accounts, err := client.Accounts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	var entries []store.Entry
	for _, acc := range accounts {
		entry := store.Entry{
			ID:       acc.ID,
			Name:     acc.Name,
			Password: acc.Password,
			Notes:    acc.Notes,
			Group:    acc.Group,
			Username: acc.Username,
			URL:      acc.URL,
		}
		if entry.Group == "" {
			entry.Group = acc.Share
		}
		entries = append(entries, entry)
	}
	return entries, nil
}
