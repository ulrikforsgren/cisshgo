package utils

import (
    // "fmt"
    "testing"
)

type contextMatch struct {
	match           bool   // boolean of if a match is expected
	matchedCommand  *TranscriptMapContext
	multipleMatches bool   // Were multiple commands matched?
}

func TestContextMatch(t *testing.T) {

	// Create a fake SupportedCommands map
	mySupportedContexts := []*TranscriptMapContext{
		{ Cmd: "Dummy Top Node", Id: 1, Up: 0 },
		{ Cmd: "interface", Id: 2, Up: 1 },
		{ Cmd: "interface <J>Ethernet <R>[0-9]+/[0-9]+", Id: 3, Up: 1 },
		{ Cmd: "interface Vlan <R>[0-9]+", Id: 4, Up: 1 },
	}

    _, contextHierarchy, _ := CompileMatches(mySupportedContexts)
	inputs := make(map[string]contextMatch)

	inputs["fail"] =                        contextMatch{false, nil, false}
	inputs["interface"] =                   contextMatch{true, &TranscriptMapContext{ Id: 2, Up: 1 }, false}
	inputs["inter"] =                       contextMatch{true, &TranscriptMapContext{ Id: 2, Up: 1 }, false}
	inputs["interface none"] =              contextMatch{false, nil, false}
	inputs["interface interface none"] =    contextMatch{false, nil, false}
	inputs["interface ethernet 0/0"] =      contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["inter ether 0/0"] =             contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["interface ethernet0/0"] =       contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["inter ethernet0/0"] =           contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["interface ethernet 100/100"] =  contextMatch{true, &TranscriptMapContext{ Id: 3, Up: 1 }, false}
	inputs["interface ethernet 0/0 none"] = contextMatch{false, nil, false}
	inputs["inter ether0/0"] =              contextMatch{false, nil, false}
	inputs["interface ethernet N100/100"] = contextMatch{false, nil, false}
	inputs["interface Vlan 100"] =          contextMatch{true, &TranscriptMapContext{ Id: 4, Up: 1 }, false}
	inputs["interface Vlan 100/1"] =        contextMatch{false, nil, false}
	inputs["int ethe"] =                    contextMatch{false, nil, false}
	inputs["none none"] =                   contextMatch{false, nil, false}

	for input, expected := range inputs {
        // fmt.Println(input, expected)
		match, matchedCommand, multipleMatches, err := ContextMatch(input, &(*contextHierarchy)[1].Commands)
		if err != nil {
			t.Errorf("Unknown Error: %s", err)
		}
		if match == false && match == expected.match && matchedCommand == nil && multipleMatches == expected.multipleMatches {
            // pass
        } else if match == true && matchedCommand != nil && CompareCommands(matchedCommand.Context, expected.matchedCommand) && multipleMatches == expected.multipleMatches {
            // pass
        } else {
            var ctx *TranscriptMapContext = nil
            if matchedCommand != nil {
                ctx = matchedCommand.Context
            }
			t.Errorf(
				"ContextMatch('%s', %v) = (%t, '%v', %t); want (%t, '%v', %t)",
				input,
				mySupportedContexts,
				match,
				ctx,
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
