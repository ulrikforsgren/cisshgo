package utils

import (
    "testing"
)

type inputMatch struct {
	match           bool   // boolean of if a match is expected
	matchedCommand  string // string of expected match to this input
	multipleMatches bool   // Were multiple commands matched?
}

func TestCmdMatch(t *testing.T) {

	// Create a fake SupportedCommands map
	mySupportedCommands := map[string]string {
		"show version":    "_show_version_",
		"show vlan":       "_show_vlan_",
		"show vlan brief": "_show_vlan_brief_",
		"reboot":          "_reboot_",
	}

    compiledCommands, _ := CompileCommands(mySupportedCommands)
	inputs := make(map[string]inputMatch)

	inputs["show version"] = inputMatch{true, "_show_version_", false}      // Should match "show version"
	inputs["show ver"] = inputMatch{true, "_show_version_", false}          // Should match "show version"
	inputs["sho vlan"] = inputMatch{true, "_show_vlan_", false}             // Should match "show vlan"
//	inputs["s v"] = inputMatch{true, "", true}                            // Should return no match
	inputs["show version made-up"] = inputMatch{false, "", false}         // Should return no match
	inputs["no version"] = inputMatch{false, "", false}                   // Should return no match
	inputs["Sho vLan BrIef"] = inputMatch{true, "_show_vlan_brief_", false} // Should match "show vlan brief"
	inputs["show vlan!"] = inputMatch{false, "", false}                   // Should return no match
	inputs["how ersion"] = inputMatch{false, "", false}                   // Should return no match
	inputs["ho ersio"] = inputMatch{false, "", false}                     // Should return no match

	for input, expected := range inputs {
		match, matchedCommand, multipleMatches, err := CommandMatch(input, compiledCommands)

		if err != nil {
			t.Errorf("Unknown Error: %s", err)
		}
		if match != expected.match ||
			matchedCommand != expected.matchedCommand ||
			multipleMatches != expected.multipleMatches {
			t.Errorf(
				"CmdMatch('%s', %v) = (%t, '%s', %t); want (%t, '%s', %t)",
				input,
				mySupportedCommands,
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
