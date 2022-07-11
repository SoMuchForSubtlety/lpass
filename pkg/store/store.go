package store

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/99designs/keyring"
)

const API_CREDENTIALS_ID = "4c3094e6-7337-4728-8658-3026c782f6a6"

type Error string

const ErrKeyNotFound = Error("Key not fount")

func (e Error) Error() string { return string(e) }

type Entry struct {
	ID       string
	Name     string
	Password string
	Notes    string
	Group    string
	Username string
	URL      string
}

func (e Entry) FilterValue() string {
	return e.Name
}

func ring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName:  "lastpass-cli",
		KeychainName: "lastpass",
		AllowedBackends: []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KeychainBackend,
			keyring.KeyCtlBackend,
			keyring.KWalletBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
			keyring.PassBackend,
		},
	})
}

func LoadAPICredentials() (string, string, error) {
	ring, err := ring()
	if err != nil {
		return "", "", fmt.Errorf("failed to open keyring: %w", err)
	}
	entry, err := ring.Get(API_CREDENTIALS_ID)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", "", ErrKeyNotFound
		}
		return "", "", fmt.Errorf("failed to get key: %w", err)
	}
	e, err := unmarshalEntry(entry.Data)
	if err != nil {
		return "", "", fmt.Errorf("failed to get error: %w", err)
	}

	return e.Username, e.Password, nil
}

func StoreAPICrendentials(username, password string) error {
	ring, err := ring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}
	entry := Entry{
		ID:       API_CREDENTIALS_ID,
		Username: username,
		Password: password,
	}

	data, err := entry.marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal API credentials: %w", err)
	}
	err = ring.Set(keyring.Item{
		Key:   API_CREDENTIALS_ID,
		Data:  data,
		Label: "LastPass credentials",
	})
	if err != nil {
		return fmt.Errorf("failed to store API credentials: %w", err)
	}

	return nil
}

func DeleteAPICredentials() error {
	ring, err := ring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}
	return ring.Remove(API_CREDENTIALS_ID)
}

func Store(entries []Entry) error {
	ring, err := ring()
	if err != nil {
		return fmt.Errorf("failed to open keyring: %w", err)
	}

	keys, err := ring.Keys()
	if err != nil {
		return fmt.Errorf("failed to get keys: %w", err)
	}

	for _, key := range keys {
		if key == API_CREDENTIALS_ID {
			continue
		}
		err := ring.Remove(key)
		if err != nil {
			return fmt.Errorf("failed to remove %q: %w", key, err)
		}
	}

	for _, entry := range entries {
		data, err := entry.marshal()
		if err != nil {
			return fmt.Errorf("failed to marshal entry %q: %w", entry.ID, err)
		}
		err = ring.Set(keyring.Item{
			Key:         entry.ID,
			Label:       entry.Name,
			Data:        data,
			Description: entry.Notes,
		})
		if err != nil {
			return fmt.Errorf("failed to safe entry %q: %w", entry.ID, err)
		}
	}

	return nil
}

func Load() ([]Entry, error) {
	ring, err := ring()
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	keys, err := ring.Keys()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	var accounts []Entry
	for _, key := range keys {
		item, err := ring.Get(key)
		if err != nil {
			return nil, fmt.Errorf("failed to load %q: %w", key, err)
		}
		entry, err := unmarshalEntry(item.Data)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal entry %q: %w", key, err)
		}
		accounts = append(accounts, entry)
	}

	return accounts, nil
}

func (e *Entry) marshal() ([]byte, error) {
	return json.Marshal(e)
}

func unmarshalEntry(data []byte) (Entry, error) {
	var e Entry
	err := json.Unmarshal(data, &e)
	return e, err
}
