package main

import (
    "io"
    "log"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/ssh_server/handlers"
	"github.com/tbotnz/cisshgo/ssh_server/sshlistners"
	"github.com/tbotnz/cisshgo/utils"
)

func main() {

	// Parse the command line arguments
	args, transcript := utils.ParseArgs()

    if args.Silent {
        log.SetOutput(io.Discard)
    }
	// Make a Channel named "done" for handling Goroutines, which expects a bool as return value
	done := make(chan bool, 1)

	// Iterate through the server ports and spawn a Goroutine for each
	for index := 0; index < args.Listeners; index++ {
        // Initialize our fake device
        aFakeDevice := fakedevices.InitGeneric(
            args,
            *args.Vendor,          // Vendor
            *args.Platform,        // Platform
            transcript,       // Transcript map with locations of command output to play back
            args.Listeners,
            args.StartingPort,
        )
		// Today this is just spawning a generic listener.
		// In the future, this is where we could split out listeners/handlers by device type.
		go sshlistners.GenericListener(args, aFakeDevice, args.StartingPort+index, handlers.GenericCiscoHandler, done)
	}

	// Receive all the values from the channel (essentially wait on it to be empty)
	<-done
}
