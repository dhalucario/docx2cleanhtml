package main

import (
	"fmt"
	"log"
	"os"

	"leong/docx2clearhtml/settingsStorage"

	"baliance.com/gooxml/document"
)

func main() {

	var pgs programSettings.ProgramSettings = programSettings.MakeProgramSettings(nil)
	args := os.Args[1:len(os.Args)]

	pgs.RegisterCommandLineSetting("verbose", programSettings.CommandLineArgument{
		Short:            "v",
		Long:             "verbose",
		DefaultValue:     false,
		MaxArgumentParam: 0,
	})
	wg := pgs.ReadCommandLineSettings(args)
	(*wg).Wait()

	logVerbose("Config Loaded", &pgs)

	doc, err := document.Open("test.docx")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(doc.Paragraphs())

}

func logVerbose(output string, pgs *programSettings.ProgramSettings) {
	if (*pgs).Get("verbose").(bool) {
		fmt.Println(output)
	}
}
