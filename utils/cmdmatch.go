package utils

import (
	// "fmt"
    "regexp"
	"strings"
)

type CommandPattern struct {
    Pattern []interface{}
    File string
}

type MatchCommands []*CommandPattern

func CompileCommands(supportedCommands map[string]string) (*MatchCommands, error) {
    fieldsMap := make(MatchCommands, 0)

	for cmd, file := range supportedCommands {
        fields := strings.Fields(strings.ToLower(cmd))
        comp_fields := make([]interface{}, len(fields))
        for n,f := range fields {
            if  strings.HasPrefix(f, "<r>") {
                comp_fields[n] =regexp.MustCompile("^"+f[3:]+"$")
            } else {
                comp_fields[n] = f
            }
        }
        fieldsMap = append(fieldsMap, &CommandPattern{Pattern: comp_fields, File: file})
	}

    return &fieldsMap, nil
}


// CmdMatch searches the provided supportedCommands to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//  multipleMatches: bool
//	error
// func CmdMatch(userInput string, supportedCommands *MatchCommands) (bool, string, bool, error) {
// 
// 	// Setup our return variables
// 	match := false
// 	matchedCmd := ""
// 	multipleMatches := false
// 
// 	// Setup a Map to hold any possibleMatches as keys, and the string.Fields as values
// 	possibleMatches := make(map[string][]string)
// 
// 	// Turn our input string into fields
// 	// fmt.Printf("userInput: %s\n", userInput)
// 	userInput = strings.ToLower(userInput) // Lowercase the user input
// 	userInputFields := strings.Fields(userInput)
// 
// 	// Iterate through all the commands in the supportedCommands map
// 	for supportedCommand := range supportedCommands {
// 		supportedCommand := strings.ToLower(supportedCommand) // Lowercase our supported command
// 		commandFields := strings.Fields(supportedCommand)
// 
// 		// Match against the 1st field in each command,
// 		// and that the number of fields is the same,
// 		// to find any possibleMatches.
// 		if strings.HasPrefix(commandFields[0], userInputFields[0]) &&
// 			(len(commandFields) == len(userInputFields)) {
// 			// fmt.Printf("supportedCommand: %s\n", k)
// 			possibleMatches[supportedCommand] = commandFields
// 		}
// 	}
// 
// 	// Setup a map to hold our closestMatch(es)
// 	closestMatch := make(map[string]struct{})
// 
// 	// Iterate through all possibleMatches to find the best match
// 	// fmt.Printf("possibleMatches: %+v\n", possibleMatches)
// 	for possibleMatch := range possibleMatches {
// 
// 		// First evaluate if we have an exact string match and break/return that
// 		if userInput == possibleMatch {
// 			closestMatch[possibleMatch] = struct{}{}
// 			break
// 		}
// 
// 		// Next, test if the entire input is contained within one of our commands
// 		if strings.HasPrefix(possibleMatch, userInput) {
// 			closestMatch[possibleMatch] = struct{}{}
// 			break
// 		}
// 
// 		// Next delve into the fields and find best match
// 		for p, possibleMatchField := range possibleMatches[possibleMatch] {
// 			// fmt.Printf("possibleMatchField: %s\n", possibleMatchField)
// 			if !strings.HasPrefix(possibleMatchField, userInputFields[p]) {
// 				// We did not get a match on this field, break
// 				break
// 			}
// 			// fmt.Printf("%d\n", p)
// 			// fmt.Printf("length of possibleMatch fields: %d\n", len(possibleMatches[possibleMatch]))
// 			if p == (len(possibleMatches[possibleMatch]) - 1) {
// 				closestMatch[possibleMatch] = struct{}{}
// 			}
// 		}
// 	}
// 
// 	// Evaluate our closestMatch(es)
// 	if len(closestMatch) > 1 {
// 		// We had more than two matches to all conditions, return no match!
// 		fmt.Printf("multiple matchedCmds: %s\n", closestMatch)
// 		match = true
// 		matchedCmd = ""
// 		multipleMatches = true
// 	} else if len(closestMatch) < 1 {
// 		// We had _NO_ matches to any conditions, return no match!
// 		// fmt.Printf("no matchedCmds\n")
// 		match = false
// 		matchedCmd = ""
// 	} else {
// 		match = true
// 		for k := range closestMatch {
// 			matchedCmd = k
// 		}
// 
// 	}
// 
// 	// fmt.Printf("matchedCmd: %s\n\n", matchedCmd)
// 	return match, matchedCmd, multipleMatches, nil
// }



// CommandMatch searches the provided supportedContexts to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//	error
func CommandMatch(userInput string, supportedCommands *MatchCommands) (bool, string, bool, error) {

	// Setup our return variables
	match := false
    var matchedCommand string

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedContexts map and create
    // regexp.
	for _, cmd := range *supportedCommands {
	    file := cmd.File
        contextFields := cmd.Pattern
        if len(userInputFields) == len(contextFields) {
	        match = true
            fieldsLoop: for n,f := range contextFields {
                switch f.(type) {
                case string:
                    //fmt.Println("COMP string", f, userInputFields[n])
                    if !strings.HasPrefix(f.(string), userInputFields[n]) {
                        match = false
                        //fmt.Println("NO MATCH!")
                        break fieldsLoop
                    }
                    //fmt.Println("MATCH!")
                default: // *regexp.Regexp
                    //fmt.Println("COMP regexp", f, userInputFields[n])
                    if f.(*regexp.Regexp).MatchString(userInputFields[n]) == false {
                        match = false
                        //fmt.Println("NO MATCH!")
                        break fieldsLoop
                    }
                    //fmt.Println("MATCH!")
                }
            }
            if match {
                matchedCommand = file
                break
            }
        }
	}

	return match, matchedCommand, false, nil
}



type Pattern interface{
    Match(string)(int)
}


type StringMatch struct {
    s string
}

func (s StringMatch)Match(cmd string)(int) {
    // fmt.Print("StringMatch")
    if strings.HasPrefix(s.s, cmd) {
        // fmt.Println(" -> MATCH")
        return 1
    }
    // fmt.Println(" -> NOMATCH")
    return 0
}

type RegexpMatch struct {
    r *regexp.Regexp
}

func (r RegexpMatch)Match(cmd string)(int) {
    // fmt.Print("RegexpMatch")
    if r.r.MatchString(cmd) {
        // fmt.Println(" -> MATCH")
        return 1
    }
    // fmt.Println(" -> NOMATCH")
    return 0
}

type JoinKeyMatch struct {
    s string
    r *regexp.Regexp
}

func (j JoinKeyMatch)Match(cmd string)(int) {
    // fmt.Print("JoinKeyMatch")
    if strings.HasPrefix(j.s, cmd) {
        // fmt.Println(" -> s MATCH")
        return 1
    }
    if j.r.MatchString(cmd) {
        // fmt.Println(" -> r MATCH")
        return 2
    }
    // fmt.Println(" -> NOMATCH")
    return 0
}



type ContextPattern struct {
    Pattern []Pattern
    Context *TranscriptMapContext
    Commands MatchContexts
}

type MatchContexts []*ContextPattern

type ContextHierarchy map[uint]*ContextPattern


func CompileMatches(supportedContexts []*TranscriptMapContext) (*MatchContexts, *ContextHierarchy, error) {
    flatMap := make(MatchContexts, 0)
    contextHierarchy := make(ContextHierarchy)

	for _, ctx := range supportedContexts {
        fields := strings.Fields(strings.ToLower(ctx.Cmd))
        comp_fields := make([]Pattern, len(fields))
        for n,f := range fields {
            if strings.HasPrefix(f, "<r>") {
                comp_fields[n] = RegexpMatch{regexp.MustCompile("^"+f[3:]+"$")}
            } else if strings.HasPrefix(f, "<j>") {
                r := fields[n+1]
                // TODO: Check statement existence + type <r>
                m := JoinKeyMatch{
                        s: f[3:],
                        r: regexp.MustCompile("^"+f[3:]+r[3:]+"$"),
                }
                comp_fields[n] = m
            } else {
                comp_fields[n] = StringMatch{f}
            }
        }
        flatMap = append(flatMap, &ContextPattern{ Pattern: comp_fields, Context: ctx})
	}

    // Build hierarchy index
    for _, mode := range flatMap {
        contextHierarchy[mode.Context.Id] = mode
    }

    // Build hierarchy
    for _, mode := range flatMap {
        if mode.Context.Up != 0 {
            contextHierarchy[mode.Context.Up].Commands = append(contextHierarchy[mode.Context.Up].Commands, mode)
        }
    }

    return &flatMap, &contextHierarchy, nil
}


// ContextMatch searches the provided supportedContexts to find a match for the provided userInput
// Returns:
//	match: bool
// 	matchedCommand: string
//	error
func ContextMatch(userInput string, supportedContexts *MatchContexts) (bool, *ContextPattern, bool, error) {

	// Setup our return variables
//    fmt.Println("ContextMatch")
	match := false
    var matchedContext *ContextPattern = nil

	// Turn our input string into fields
	// fmt.Printf("userInput: %s\n", userInput)
	userInput = strings.ToLower(userInput) // Lowercase the user input
	userInputFields := strings.Fields(userInput)

	// Iterate through all the commands in the supportedContexts map and create
    // regexp.
	for _, cmd := range *supportedContexts {
        // fmt.Println("Matching command", cmd.Context.Cmd)
	    supportedContext := cmd
        contextFields := cmd.Pattern

        n_i := 0
        n_c := 0
	    match = true
        for n_i<len(userInputFields) && n_c<len(contextFields) {
            f := contextFields[n_c]
            m_c := f.Match(userInputFields[n_i])
            if m_c == 0 {
                match = false
                // fmt.Println("BREAK")
                break
            }
            n_c += m_c
            n_i += 1
        }
        if match {
            // fmt.Println(n_i, len(userInputFields), n_c, len(contextFields))
            if n_i==len(userInputFields) && n_c==len(contextFields) {
                matchedContext = supportedContext
                break
            }
            // fmt.Println("Len mismatch", match)
        }
	}

	return match, matchedContext, false, nil
}
