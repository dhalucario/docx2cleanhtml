package main

import (
	"baliance.com/gooxml/document"
	"fmt"
	"leong/docx2clearhtml/htmlAliasFormatter"
	programSettings "leong/docx2clearhtml/settingsStorage"
	"log"
	"os"
	"strings"
)

func main() {

	htmlElementAliases := map[string]string {
		"title": "<h1>%s</h1>",
		"heading 1": "<h2>%s</h2>",
		"heading 2": "<h3>%s</h3>",
		"heading 3": "<h4>%s</h4>",
		"heading 4": "<h5>%s</h5>",
	}
	haf := htmlAliasFormatter.New(htmlElementAliases)

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

	doc, err := document.Open(pgs.Get("in").(string))
	if err != nil {
		log.Fatal(err)
	}

	styleIdNames := make(map[string]string)
	for _, s := range doc.Styles.Styles() {
		pgs.VerbosePrintf("%s (ID: %s)\n", s.Name(), s.StyleID())
		styleIdNames[s.StyleID()] = s.Name()
	}

	para := doc.Paragraphs()
	paraBuffer := ""
	for _, p := range para {
		style := p.Style()
		paraBuffer = ""
		for _, r := range p.Runs() {
			paraBuffer = paraBuffer + r.Text()
		}

		fmt.Println(haf.ConvertToHtml(paraBuffer, strings.ToLower(styleIdNames[style])))
	}
}


