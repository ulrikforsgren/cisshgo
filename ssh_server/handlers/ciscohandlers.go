// Package handlers contains SSH Handlers for specific device types
// in order to best emulate their actual behavior.
package handlers

import (
//    "fmt"
	"log"
    "net"
    "regexp"
    "strconv"
	"strings"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/tbotnz/cisshgo/fakedevices"
	"github.com/tbotnz/cisshgo/utils"
)

// GenericCiscoHandler function handles generic Cisco style sessions
func GenericCiscoHandler(myFakeDevice *fakedevices.FakeDevice, s ssh.Session) {

	// Prepare the "ssh.DefaultHandler", this houses our device specific functionality

		log.Printf("%s: terminal connected\n", s.LocalAddr())
		var Port = uint(s.LocalAddr().(*net.TCPAddr).Port)

		// Setup our initial "context" or prompt
		ContextState := (*myFakeDevice.ContextHierarchy)[1] // base

		// Setup a terminal with the hostname + initial context state as a prompt
		term := terminal.NewTerminal(s, myFakeDevice.Hostname+ContextState.Context.Mode)

        re_term_width := regexp.MustCompile(`^terminal +width +(\d+)\s*$`)

		// Iterate over any user input that is provided at the terminal
		for {
			userInput, err := term.ReadLine()
			if err != nil {
				break
			}
		    log.Printf("%s: %s %v\n", s.LocalAddr(), ContextState.Context.Mode, userInput)

			// Handle any empty input (assumed to just be a carriage return)
			if userInput == "" {
				// return nothing but a newline if nothing is entered
				term.Write([]byte(""))
				continue
			}
			// Run userInput through the command matcher to look for contextSwitching commands
			matchPrompt, matchedPrompt, multiplePromptMatches, err := utils.ContextMatch(
				userInput, &ContextState.Commands,
			)
			if err != nil {
				log.Println(err)
				break
			}
			// Handle any context switching
			if matchPrompt && !multiplePromptMatches {
				// switch contexts as needed
				ContextState = matchedPrompt
				term.SetPrompt(string(
					myFakeDevice.Hostname+ContextState.Context.Mode,
				))
				continue
			} else if userInput == "exit" || userInput == "end" || strings.HasPrefix(ContextState.Context.ExitCmd, userInput) {
                up := ContextState.Context.Up
                if ContextState.Context.ExitTo != "" {
                    up64, _ := strconv.ParseUint(ContextState.Context.ExitTo, 10 ,0)
                    up = uint(up64)
                }
				// Back out of the lower contexts, i.e. drop from "(config)#" to "#"
				if up == 0 {
					break
				} else {
				    ContextState = (*myFakeDevice.ContextHierarchy)[up]
					term.SetPrompt(string(
						myFakeDevice.Hostname + ContextState.Context.Mode,
					))
					continue
				}
			} else if userInput == "reset state" {
				term.Write(append([]byte("Resetting State..."), '\n'))
				ContextState = (*myFakeDevice.ContextHierarchy)[0] // base
				myFakeDevice.Hostname = myFakeDevice.DefaultHostname
				term.SetPrompt(string(
					myFakeDevice.Hostname + ContextState.Context.Mode,
				))
				continue 
			} else if m := re_term_width.FindAllStringSubmatch(userInput, -1); m != nil {
                width, err := strconv.ParseInt(m[0][1], 10, 0)
                if err == nil {
                    if width == 0 {
                        width = 511 // Maximizing width
                    }
                    _, height := term.GetSize()
		            log.Printf("New terminal width: %d\n", width)
                    term.SetSize(int(width), height)
                } else {
                    term.Write(append([]byte("% Unknown command:  \""+userInput+"\""), '\n'))
                }
                continue
            }

			// Split user input into fields
			userInputFields := strings.Fields(userInput)

			// Handle hostname changes
			if userInputFields[0] == "hostname" && ContextState.Context.Id == 3 {
				// Set the hostname to the values after "hostname" in the userInputFields
				myFakeDevice.Hostname = strings.Join(userInputFields[1:], " ")
				log.Printf("Setting hostname to %s\n", myFakeDevice.Hostname)
				term.SetPrompt(myFakeDevice.Hostname + ContextState.Context.Mode)
				continue
			}

			// Run userInput through the command matcher to look at supportedCommands
			//match, matchedCommand, multipleMatches, err := utils.CmdMatch(userInput, myFakeDevice.SupportedCommands)
			match, matchedCommand, multipleMatches, err := utils.CommandMatch(userInput, myFakeDevice.CommandSearch)
			if err != nil {
				log.Println(err)
				break
			}

			if match && !multipleMatches {
                index := 0
                if matchedCommand.PerDeviceData {
                    index =int(Port)-myFakeDevice.StartingPort
                }
				// Render the matched command output
				output, err := fakedevices.TranscriptReader(
					matchedCommand.CmdData[index],
					myFakeDevice,
				)
				if err != nil {
					log.Fatal(err)
				}

				// Write the output of our matched command
				term.Write(append([]byte(output), '\n'))
				continue
			} else if multipleMatches {
				// Multiple commands were matched, throw ambiguous command
				term.Write(append([]byte("% Ambiguous command:  \""+userInput+"\""), '\n'))
				continue
			} else {
                if ContextState.Context.Id < 3 { // Not in config mode
				    // If all else fails, we did not recognize the input!
				    term.Write(append([]byte("% Unknown command:  \""+userInput+"\""), '\n'))
                }
				continue
			}
		}
		log.Printf("%s: terminal closed\n", s.LocalAddr())

}
