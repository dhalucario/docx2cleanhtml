package main

import (
	"fmt"
	"leong/docx2cleanhtml/settingsStorage"
	"leong/docx2cleanhtml/simpleDocxParser"
	"log"
	"os"
)

func main() {

	pgs := programSettings.New(nil)
	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:             "v",
		Long:              "verbose",
		DefaultValue:      false,
		MultipleArguments: false,
		MaxArgumentParam:  0,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("verbose", true)
		},
	})

	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:             "i",
		Long:              "in",
		DefaultValue:      "",
		MultipleArguments: true,
		MaxArgumentParam:  1,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("in", commandLineArgs[0])
		},
	})

	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:             "o",
		Long:              "out",
		DefaultValue:      "",
		MultipleArguments: true,
		MaxArgumentParam:  1,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("out", true)
		},
	})

	// Used to get unique id's in filenames
	pgs.Set("tempcounter", 0)

	args := os.Args[1:len(os.Args)]
	pgs.ReadCommandLineSettings(args)

	doc, err := simpleDocxParser.New(pgs.Get("in").(string), &pgs)
	if err != nil {
		log.Fatal(err)
	}
	doc.ReadRelations()

	//TODO: Run as web service
	if pgs.Get("out") != "" {
		// TODO: Write to file
	} else {
		fmt.Print(doc.HTML())
	}

}
