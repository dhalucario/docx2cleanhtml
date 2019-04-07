package main

import (
	"fmt"
	"log"
	"os"

	"leong/docx2clearhtml/settingsStorage"
	"leong/docx2clearhtml/htmlAliasFormatter"

	"baliance.com/gooxml/document"
)

func main() {

	htmlElementAliases := map[string]string {
		"heading 1": "<h1>%s</h1>",
		"heading 2": "<h2>%s</h2>",
		"heading 3": "<h3>%s</h3>",
		"heading 4": "<h4>%s</h4>",
		"heading 5": "<h5>%s</h5>",
	}
	haf := htmlAliasFormatter.New(htmlElementAliases)

	pgs := programSettings.New(nil)
	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:            "v",
		Long:             "verbose",
		DefaultValue:     false,
		MultipleArguments: true,
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
		MaxArgumentParam: 0,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("verbose", true)
		},
	})

	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short:            "o",
		Long:             "out",
		DefaultValue:     "",
		MultipleArguments: true,
		MaxArgumentParam: 1,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			ps.Set("verbose", true)
		},
	})

	args := os.Args[1:len(os.Args)]
	wg := pgs.ReadCommandLineSettings(args)
	(*wg).Wait()

	logVerbose("Loading doc file...", &pgs)
	doc, err := document.Open("test.docx")

	if err != nil {
		log.Fatal(err)
	}

	logVerbose("Loaded doc", &pgs)

	styleIdNames := make(map[string]string)
	for _, s := range doc.Styles.Styles() {
		styleIdNames[s.StyleID()] = s.Name()
	}

	para := doc.Paragraphs()
	for _, p := range para {
		style := p.Style()
		for _, r := range p.Runs() {
			fmt.Println(haf.ConvertToHtml(r.Text(), styleIdNames[style]))
		}
	}

}

func logVerbose(output string, pgs *programSettings.ProgramSettings) {
	if (*pgs).Get("verbose").(bool) {
		fmt.Println(output)
	}
}
