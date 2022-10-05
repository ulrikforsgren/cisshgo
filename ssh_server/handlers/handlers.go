package handlers

import "github.com/tbotnz/cisshgo/fakedevices"
import "github.com/gliderlabs/ssh"

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*fakedevices.FakeDevice, ssh.Session)
