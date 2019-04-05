package main

import (
	"os"
	"log"
	//"strings"

	"leong/docx2cleanhtml/settingsStorage"


	//_ "baliance.com/gooxml/document"
)

func main() {

	var pgs programSettings.ProgramSettings = programSettings.MakeProgramSettings(nil)
	var log log.Logger = make(log.Logger)

	pgs.RegisterCommandLineSetting("verbose", programSettings.CommandLineArgument{
		"v",
		"verbose",
		false,
		false,
		0,
	})
	pgs.ReadCommandLineSettings(&os.Args)


}
