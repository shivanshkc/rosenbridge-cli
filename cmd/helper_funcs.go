package cmd

import (
	"context"
	"os"
	"time"

	"github.com/shivanshkc/rosenbridge-cli/lib"

	"github.com/fatih/color"
)

// exitWithPrintf prints the provided message in Printf style and then calls os.Exit with provided code.
func exitWithPrintf(code int, format string, a ...interface{}) {
	switch code {
	case 0:
		color.Green(format+"\n", a...)
	default:
		color.Red(format+"\n", a...)
	}
	os.Exit(code)
}

// printMessage prints the provided message in appropriate format.
// If the provided message is nil, it prints that the message is ill-formatted.
func printMessage(ctx context.Context, inMessage *lib.IncomingMessage, err error) {
	if err != nil {
		color.Red(">> [%s] Error while reading the message: %s\n", err.Error())
		return
	}
	if inMessage == nil {
		color.Red(">> [%s] Message is of unrecognized format.\n", time.Now().Format(time.Kitchen))
		return
	}
	color.Yellow(">> [%s] %s: %s\n", time.Now().Format(time.Kitchen), inMessage.SenderID, inMessage.Message)
}
