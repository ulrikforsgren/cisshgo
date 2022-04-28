package fakedevices

import (
    "fmt"
	"io/ioutil"
	"log"
    "strconv"
    "strings"

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
	CommandSearch     *utils.MatchCommands // What commands this fake device supports
	ContextSearch     *utils.MatchContexts // The available CLI prompt/contexts on this fake device
	ContextHierarchy  *utils.ContextHierarchy // The hierarchy of the available contexts
    StartingPort      int
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
	transcript *utils.Transcript,
    numListeners int,
    startingPort int,
) *FakeDevice {

	// Iterate through the command transcripts and read their contents into our supported commands
	for _, v := range *transcript.CommandSearch {
        //fmt.Println(v.File)
        if strings.Contains(v.File.Name, "<PORT>") {
            v.File.PerDeviceData = true
		    v.File.CmdData = map[int]string {}
            for i:=0; i<numListeners; i++ {
                var name = strings.ReplaceAll(v.File.Name, "<PORT>",strconv.Itoa(startingPort+i))
                fmt.Println(name)
		        v.File.CmdData[i] = readFile(name)
            }
        } else {
            v.File.PerDeviceData = false
		    v.File.CmdData = map[int]string {}
		    v.File.CmdData[0] = readFile(v.File.Name)
        }
	}


	// Create our fake device and return it
	myFakeDevice := FakeDevice{
		Vendor:            transcript.Vendor,
		Platform:          "Undefined", // Currently not used...
		Hostname:          transcript.Hostname,
		DefaultHostname:   transcript.Hostname,
		Password:          transcript.Password,
		CommandSearch:     transcript.CommandSearch,
		ContextSearch:     transcript.ContextSearch,
		ContextHierarchy:  transcript.ContextHierarchy,
        StartingPort:      startingPort,
	}

	//fmt.Printf("\n%+v\n", myFakeDevice)
	return &myFakeDevice
}
