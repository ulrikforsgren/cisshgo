package utils

import (
//    "fmt"
    "testing"
)

type contextMatch struct {
	match           bool   // boolean of if a match is expected
	matchedCommand  *TranscriptMapContext
	multipleMatches bool   // Were multiple commands matched?
}

func TestContextMatch(t *testing.T) {

	// Create a fake SupportedCommands map
	mySupportedContexts := map[string]*TranscriptMapContext{
		"interface Ethernet": { Id: 2, Up: 1 },
		"interface <R>Ethernet[0-9]+/[0-9]+":  { Id: 3, Up: 1 },
	}

    compiledContexts, _ := CompileMatches(mySupportedContexts)
	inputs := make(map[string]contextMatch)

	inputs["interface ethernet"] =        contextMatch{true, &TranscriptMapContext{ Id: 2, Up: 1 }, false}
	inputs["interface ethernet0/0"] =     contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["int ethe"] =                  contextMatch{true, &TranscriptMapContext{ Id: 2, Up: 1 }, false}
	inputs["interface ethernet0/0"] =     contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["interface ethernet100/100"] = contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["s v"] =                       contextMatch{false, nil, false}
    inputs["show version made-up"] =      contextMatch{false, nil, false}

	for input, expected := range inputs {
		match, matchedCommand, multipleMatches, err := ContextMatch(input, compiledContexts)
		if err != nil {
			t.Errorf("Unknown Error: %s", err)
		}
		if match != expected.match ||
            !CompareCommands(matchedCommand, expected.matchedCommand) ||
			multipleMatches != expected.multipleMatches {
			t.Errorf(
				"ContextMatch('%s', %v) = (%t, '%v', %t); want (%t, '%v', %t)",
				input,
				mySupportedContexts,
				match,
				matchedCommand,
				multipleMatches,
				expected.match,
				expected.matchedCommand,
				expected.multipleMatches,
			)
		}
	}

}

func CompareCommands(a *TranscriptMapContext, b *TranscriptMapContext)(bool) {
    return a == b || (a.Id == b.Id &&  a.Up == b.Up)
}
