package main

import (
	"fmt"
	"leong/docx2cleanhtml/settingsStorage"
	"leong/docx2cleanhtml/simpleDocxParser"
	"os"
)

func main() {

	/*htmlElementAliases := map[string]string {
		"title": "<h1>%s</h1>",
		"heading 1": "<h2>%s</h2>",
		"heading 2": "<h3>%s</h3>",
		"heading 3": "<h4>%s</h4>",
		"heading 4": "<h5>%s</h5>",
	}
	haf := htmlAliasFormatter.New(htmlElementAliases)*/

	pgs := programSettings.New(nil)
	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:            "v",
		Long:             "verbose",
		DefaultValue:     false,
		MultipleArguments: false,
		MaxArgumentParam: 0,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("verbose", true)
		},
	})

	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:            "i",
		Long:             "in",
		DefaultValue:     "",
		MultipleArguments: true,
		MaxArgumentParam: 1,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("in", commandLineArgs[0])
		},
	})

	args := os.Args[1:len(os.Args)]
	pgs.ReadCommandLineSettings(args)

	doc, err := simpleDocxParser.New(pgs.Get("in").(string))

	if err != nil {
		fmt.Print(err)
	}

	doc.GetHTML()

}
