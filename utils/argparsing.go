package utils

/*
 TODO:
  - Separate reading config and creating structures to be use in the device.
    - Read config
    - Compile regexps and ceate hierarchy map
    - Create hierarchical context list
*/

import (
	"flag"
//    "fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type TranscriptMapContext struct {
    Id                 uint          `yaml:"id" json:"id"`
    Up                 uint          `yaml:"up" json:"up"`
    Mode               string       `yaml:"mode" json:"mode"`
    ExitCmd            string       `yaml:"exit-cmd" json:"exit-cmd"`
}

// TranscriptMapPlatform struct for use inside of a TranscriptMap struct
type TranscriptMapPlatform struct {
	Vendor             string            `yaml:"vendor" json:"vendor"`
	Hostname           string            `yaml:"hostname" json:"hostname"`
	Password           string            `yaml:"password" json:"password"`
	CommandTranscripts map[string]string `yaml:"command_transcripts" json:"command_transcripts"`
	ContextSearch      map[string]*TranscriptMapContext `yaml:"context_search" json:"context_search"`
}

// TranscriptMap Struct for modeling the TranscriptMap YAML
type TranscriptMap struct {
	Platforms []map[string]TranscriptMapPlatform `yaml:"platforms" json:"platforms"`
}

type Transcript struct {
	Vendor             string
	Hostname           string
	Password           string
	CommandSearch      *MatchCommands
	ContextSearch      *MatchContexts
    ContextHierarchy   *ContextHierarchy
}

// ParseArgs parses command line arguments for cisshgo
func ParseArgs() (*string, *string, int, *int, *Transcript) {
	// Gather command line arguments and parse them
	vendor := flag.String("vendor", "cisco", "Device vendor")
	platform := flag.String("platform", "csr1000v", "Device platform")
	listenersPtr := flag.Int("listeners", 50, "How many listeners do you wish to spawn?")
	startingPortPtr := flag.Int("startingPort", 10000, "What port do you want to start at?")
	transcriptMapPtr := flag.String(
		"transcriptMap",
		"transcripts/transcript_map.yaml",
		"What file contains the map of commands to transcribed output?",
	)
	flag.Parse()

	// How many total listeners will we have?
	numListeners := *startingPortPtr + *listenersPtr

	// Gather the command transcripts and create a map of vendor/platform/command
	transcriptMapRaw, err := ioutil.ReadFile(*transcriptMapPtr)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// fmt.Printf("Raw Transcript Map from file:\n\n%s\n", transcriptMapRaw)

	myTranscriptMap := TranscriptMap{}
	err = yaml.UnmarshalStrict([]byte(transcriptMapRaw), &myTranscriptMap)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

    tMap := myTranscriptMap.Platforms[0][*platform]
	//fmt.Printf("YAML Parsed Transcript Map:\n\n%+v\n", myTranscriptMap)
    // Create context hierarchy
    compiledSupportedCommands, _ := CompileCommands(tMap.CommandTranscripts)
    compiledContextSearch, contextHierarchy, _ := CompileMatches(tMap.ContextSearch)

    transcript := Transcript{
	                  Vendor:           tMap.Vendor,
	                  Hostname:         tMap.Hostname,
	                  Password:         tMap.Password,
	                  CommandSearch:    compiledSupportedCommands,
	                  ContextSearch:    compiledContextSearch,
                      ContextHierarchy: contextHierarchy,
                  }

	return vendor, platform, numListeners, startingPortPtr, &transcript
}
