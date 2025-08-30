# Confman

A cross-platform library for managing the configuration files of an application.

# Installation

```
go get -u github.com/voidwyrm-2/confman@latest
```

# Example

```go
package main

import (
	"fmt"
	"strings"

	"github.com/voidwyrm-2/confman"
)

const (
    PUBKEY_FILE = "sshkey.pub"
    USERS_FILE = "users.json"
)

func main() {
	config, path, err := confman.Open("example_app")
	if err != nil {
		panic(err)
	}

	defer config.Close()

	fmt.Println("Config opened at", path)

	config.DefaultString(PUBKEY_FILE, 0o644, "[SSHKEY]")
	config.DefaultString(USERS_FILE, 0o644, `["Zeus", "Socrates", "John Hammond", "Jane Doe", "Seth Lowell", "Akira"]`)

	pubkey, err := config.ReadString(PUBKEY_FILE)
	if err != nil {
		panic(err)
	}

	fmt.Println("pubkey:", pubkey)

	var users []string

	err = config.ReadJson(USERS_FILE, &users)
	if err != nil {
		panic(err)
	}

	fmt.Printf("users:\n%s\n", strings.Join(users, "\n"))
}
```
