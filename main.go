package main

import (
	"fmt"
	"leong/docx2cleanhtml/simpleDocxParser"
	"log"
	"os"
	"runtime/debug"
	"strings"

	"leong/docx2cleanhtml/settingsStorage"
	"leong/docx2cleanhtml/webHandler"
)


func main() {

	pgs := programSettings.New(nil)

	// Used to get unique id's in filenames. Has to be set for the parser
	pgs.Set("tempcounter", 0)

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

	pgs.RegisterCommandLineSetting(programSettings.CommandLineArgument{
		Short: "wsrv",
		Long: "webserver",
		DefaultValue: nil,
		MultipleArguments: true,
		MaxArgumentParam: 2,
		CommandHandler: func(commandLineArgs []string, ps *programSettings.ProgramSettings) {
			fmt.Println("Starting Webserver")

			wsrvConfig := webHandler.WServerSettings{}

			for _, setting := range commandLineArgs {
				splitSetting := strings.Split(setting, ":")

				if len(splitSetting) != 2 {
					log.Panic("Couldn't parse:" + setting)
				}

				switch splitSetting[0] {
				case "ip":
					wsrvConfig.Ip = splitSetting[1]
				case "port":
					wsrvConfig.Port = splitSetting[1]
				}
			}

			err := os.Mkdir("./uploads", 775)
			if err != nil {
				if !os.IsExist(err) {
					fmt.Println(err)
					debug.PrintStack()
				}
			}

			err = os.MkdirAll("./public/output", 775)
			if err != nil {
				if !os.IsExist(err) {
					fmt.Println(err)
					debug.PrintStack()
				}
			}

			ps.Set("wsrv", &wsrvConfig)

		},
	})

	args := os.Args[1:len(os.Args)]
	pgs.ReadCommandLineSettings(args)

	webServerResult := pgs.Get("wsrv")

	if webServerResult == nil {
		convertSingleFile(&pgs)
	} else {
		webServerConfig := webServerResult.(*webHandler.WServerSettings)
		webServerConfig.AutocompleteEmpty()
		docServer, err := webHandler.NewDocServer(*webServerConfig, 5, "", "")
		if err != nil {
			panic(err)
		}

		err = docServer.Run()
		if err != nil {
			panic(err)
		}
	}


}
// Sets empty values to defaults.

func convertSingleFile(pgs *programSettings.ProgramSettings) {
	doc, err := simpleDocxParser.New(pgs.Get("in").(string))
	if err != nil {
		log.Fatal(err)
	}

	err = doc.ReadRelations()
	if err != nil {
		log.Fatal(err)
	}

	if pgs.Get("out") != "" {
		// TODO: Write to file. If you want to get it into a file just pipe it.
	} else {
		fmt.Print(doc.HTML())
	}
}

