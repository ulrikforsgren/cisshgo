package handlers

import "github.com/tbotnz/cisshgo/fakedevices"
import "github.com/gliderlabs/ssh"
import "github.com/tbotnz/cisshgo/utils"

// PlatformHandler defines a default type for all platform handlers
type PlatformHandler func(*utils.CmdlineArguments, *fakedevices.FakeDevice, ssh.Session)
