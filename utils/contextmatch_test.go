package utils

import "testing"

type contextMatch struct {
	match           bool   // boolean of if a match is expected
	matchedCommand  string // string of expected match to this input
	multipleMatches bool   // Were multiple commands matched?
}

func TestContextMatch(t *testing.T) {

	// Create a fake SupportedCommands map
	mySupportedContexts := map[string]string{
		"interface Ethernet":  "Interface Ethernet sub mode",
		"interface Ethernet[0-9]+/[0-9]+":  "Interface Ethernet sub mode",
	}

    CompileMatches(mySupportedContexts)
	inputs := make(map[string]contextMatch)

	inputs["interface ethernet"] = contextMatch{true, "interface Ethernet", false} // Should match "show version"
	inputs["interface ethernet0/0"] = contextMatch{true, "interface Ethernet[0-9]+/[0-9]+", false} // Should match "show version"
	inputs["interface ethernet100/100"] = contextMatch{true, "interface Ethernet[0-9]+/[0-9]+", false} // Should match "show version"
	inputs["s v"] = contextMatch{false, "", false}                            // Should return no match
    inputs["show version made-up"] = contextMatch{false, "", false}         // Should return no match

	for input, expected := range inputs {
		match, matchedCommand, multipleMatches, err := ContextMatch(input, mySupportedContexts)
		if err != nil {
			t.Errorf("Unknown Error: %s", err)
		}
		if match != expected.match ||
			matchedCommand != expected.matchedCommand ||
			multipleMatches != expected.multipleMatches {
			t.Errorf(
				"ContextMatch('%s', %v) = (%t, '%s', %t); want (%t, '%s', %t)",
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
