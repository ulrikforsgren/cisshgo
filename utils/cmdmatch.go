package utils

import (
	"fmt"
    "regexp"
	"strings"
)

// CmdMatch searches the provided supportedCommands to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//  multipleMatches: bool
//	error
func CmdMatch(userInput string, supportedCommands map[string]string) (bool, string, bool, error) {

	// Setup our return variables
	match := false
	matchedCmd := ""
	multipleMatches := false

	// Setup a Map to hold any possibleMatches as keys, and the string.Fields as values
	possibleMatches := make(map[string][]string)

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedCommands map
	for supportedCommand := range supportedCommands {
		supportedCommand := strings.ToLower(supportedCommand) // Lowercase our supported command
		commandFields := strings.Fields(supportedCommand)

		// Match against the 1st field in each command,
		// and that the number of fields is the same,
		// to find any possibleMatches.
		if strings.HasPrefix(commandFields[0], userInputFields[0]) &&
			(len(commandFields) == len(userInputFields)) {
			// fmt.Printf("supportedCommand: %s\n", k)
			possibleMatches[supportedCommand] = commandFields
		}
	}

	// Setup a map to hold our closestMatch(es)
	closestMatch := make(map[string]struct{})

	// Iterate through all possibleMatches to find the best match
	// fmt.Printf("possibleMatches: %+v\n", possibleMatches)
	for possibleMatch := range possibleMatches {

		// First evaluate if we have an exact string match and break/return that
		if userInput == possibleMatch {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		// Next, test if the entire input is contained within one of our commands
		if strings.HasPrefix(possibleMatch, userInput) {
			closestMatch[possibleMatch] = struct{}{}
			break
		}

		// Next delve into the fields and find best match
		for p, possibleMatchField := range possibleMatches[possibleMatch] {
			// fmt.Printf("possibleMatchField: %s\n", possibleMatchField)
			if !strings.HasPrefix(possibleMatchField, userInputFields[p]) {
				// We did not get a match on this field, break
				break
			}
			// fmt.Printf("%d\n", p)
			// fmt.Printf("length of possibleMatch fields: %d\n", len(possibleMatches[possibleMatch]))
			if p == (len(possibleMatches[possibleMatch]) - 1) {
				closestMatch[possibleMatch] = struct{}{}
			}
		}
	}

	// Evaluate our closestMatch(es)
	if len(closestMatch) > 1 {
		// We had more than two matches to all conditions, return no match!
		fmt.Printf("multiple matchedCmds: %s\n", closestMatch)
		match = true
		matchedCmd = ""
		multipleMatches = true
	} else if len(closestMatch) < 1 {
		// We had _NO_ matches to any conditions, return no match!
		// fmt.Printf("no matchedCmds\n")
		match = false
		matchedCmd = ""
	} else {
		match = true
		for k := range closestMatch {
			matchedCmd = k
		}

	}

	// fmt.Printf("matchedCmd: %s\n\n", matchedCmd)
	return match, matchedCmd, multipleMatches, nil
}



// ContextMatch searches the provided supportedContexts to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//	error
func ContextMatch1(userInput string, supportedContexts map[string]string) (bool, string, bool, error) {

	// Setup our return variables
	match := false
    matchedContext := ""

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedContexts map and create
    // regexp.
	for supportedContext := range supportedContexts {
		contextFields := strings.Fields(strings.ToLower(supportedContext))

        if len(userInputFields) == len(contextFields) {
	        match = true
            for n,f := range contextFields {
                // Compilation of regexps should be done one time at startup!
                r := regexp.MustCompile("^"+f+"$")
                if r.MatchString(userInputFields[n]) == false {
                    match = false
                    break
                }
            }
            if match {
                matchedContext = supportedContext
                break
            }
        }
	}

	return match, matchedContext, false, nil
}



// ContextMatch searches the provided supportedContexts to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//	error
func ContextMatch(userInput string, supportedContexts *MatchContexts) (bool, string, bool, error) {

	// Setup our return variables
	match := false
    matchedContext := ""

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedContexts map and create
    // regexp.
	for supportedContext, contextFields := range *supportedContexts {
        if len(userInputFields) == len(contextFields) {
	        match = true
            fieldsLoop: for n,f := range contextFields {
                switch f.(type) {
                case string:
                   // fmt.Println("COMP string", f, userInputFields[n])
                    if !strings.HasPrefix(f.(string), userInputFields[n]) {
                        match = false
                        break fieldsLoop
                    }
                    //fmt.Println("MATCH!")
                default: // *regexp.Regexp
                    if f.(*regexp.Regexp).MatchString(userInputFields[n]) == false {
                        match = false
                        break fieldsLoop
                    }
                }
            }
            if match {
                matchedContext = supportedContext
                break
            }
        }
	}

	return match, matchedContext, false, nil
}

type MatchContexts map[string][]interface{}

func CompileMatches(supportedContexts map[string]string) (*MatchContexts, error) {
    fieldsMap := make(MatchContexts)

	for ctx := range supportedContexts {
        fields := strings.Fields(strings.ToLower(ctx))
        comp_fields := make([]interface{}, len(fields))
        for n,f := range fields {
            if  strings.HasPrefix(f, "<r>") {
                comp_fields[n] =regexp.MustCompile("^"+f[3:]+"$")
            } else {
                comp_fields[n] = f
            }
        }
        fieldsMap[ctx] = comp_fields
	}

    return &fieldsMap, nil
}
