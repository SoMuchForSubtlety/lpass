package ui

import (
	"fmt"

	"github.com/SoMuchForSubtlety/lpass/pkg/store"
)

func Render(entry store.Entry) {
	var buf string
	if entry.URL != "" {
		buf += fmt.Sprintf("URL:      %s", entry.URL)
	}
	if entry.Username != "" {
		if buf != "" {
			buf += "\n"
		}
		buf += fmt.Sprintf("username: %s", entry.Username)
	}
	if entry.Password != "" {
		if buf != "" {
			buf += "\n"
		}
		buf += fmt.Sprintf("password: %s", entry.Password)
	}
	if entry.Notes != "" {
		if buf != "" {
			buf += "\n\n"
		}
		buf += entry.Notes
	}
	fmt.Println(buf)
}
