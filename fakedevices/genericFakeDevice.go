package fakedevices

import (
//    "fmt"
	"io/ioutil"
	"log"

	"github.com/tbotnz/cisshgo/utils"
)

// SupportedCommands is a map of the commands a FakeDevice supports and it's corresponding output
//type SupportedCommands map[string]string

// FakeDevice Struct for the device we will be simulating
type FakeDevice struct {
	Vendor            string            // Vendor of this fake device
	Platform          string            // Platform of this fake device
	Hostname          string            // Hostname of the fake device
	DefaultHostname   string            // Default Hostname of the fake device (for resetting)
	Password          string            // Password of the fake device
	SupportedCommands *utils.MatchCommands // What commands this fake device supports
	ContextSearch     *utils.MatchContexts // The available CLI prompt/contexts on this fake device
	ContextHierarchy  map[uint]*utils.TranscriptMapContext // The hierarchy of the available contexts
}

// readFile abstracts the standard error handling of opening and reading a file into a string
func readFile(filename string) string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

// InitGeneric builds a FakeDevice struct for use with cisshgo
func InitGeneric(
	vendor string,
	platform string,
	myTranscriptMap *utils.TranscriptMap,
    contextHierarchy map[uint]*utils.TranscriptMapContext,
) *FakeDevice {

	supportedCommands := make(map[string]string)
	contextSearch := make(map[string]*utils.TranscriptMapContext)
	commandTranscriptFiles := make(map[string]string)

	// Find the hostname, password, and other info in the data for this device
	var deviceHostname string
	var devicePassword string
	for _, fakeDevicePlatform := range myTranscriptMap.Platforms {
		// fmt.Printf("\nPlatform Map:\n%+v\n", fakeDevicePlatform)
		for k, v := range fakeDevicePlatform {
			if k == platform {
				// fmt.Printf("\nKey: %+v\n", k)
				// fmt.Printf("Value: %+v\n", v)
				deviceHostname = v.Hostname
				devicePassword = v.Password
				contextSearch = v.ContextSearch
				commandTranscriptFiles = v.CommandTranscripts
			}
		}
	}

	// Iterate through the command transcripts and read their contents into our supported commands
	for k, v := range commandTranscriptFiles {
		supportedCommands[k] = readFile(v)
	}

    compiledSupportedCommands, _ := utils.CompileCommands(supportedCommands)
    compiledContextSearch, _ := utils.CompileMatches(contextSearch)

	// Create our fake device and return it
	myFakeDevice := FakeDevice{
		Vendor:            vendor,
		Platform:          platform,
		Hostname:          deviceHostname,
		DefaultHostname:   deviceHostname,
		Password:          devicePassword,
		SupportedCommands: compiledSupportedCommands,
		ContextSearch: compiledContextSearch,
		ContextHierarchy:  contextHierarchy,
	}

	//fmt.Printf("\n%+v\n", myFakeDevice)
	return &myFakeDevice
}
