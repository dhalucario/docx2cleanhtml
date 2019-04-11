package programSettings

import (
	"fmt"
	"strings"
	"sync"
)

// Structs
type ProgramSettings struct {
	storage              map[string]interface{}
	commandLineArguments map[string]CommandLineArgument
}

type CommandLineArgument struct {
	Short             string
	Long              string
	DefaultValue      interface{}
	MultipleArguments bool
	MaxArgumentParam  int
	CommandHandler    func(arguments []string, ps *ProgramSettings)
}

// Factory Functions
func New(userValues map[string]interface{}) ProgramSettings {
	var valueStorage map[string]interface{}

	if userValues != nil {
		valueStorage = userValues
	} else {
		valueStorage = make(map[string]interface{})
	}

	return ProgramSettings{valueStorage, map[string]CommandLineArgument{}}
}

// Worker Functions

func (ps *ProgramSettings) Set(key string, value interface{}) {
	ps.storage[key] = value
}

func (ps *ProgramSettings) Get(key string) interface{} {
	return ps.storage[key]
}

func (ps *ProgramSettings) All() map[string]interface{} {
	return ps.storage
}

func (ps *ProgramSettings) Reset(key string) {
	ps.storage[key] = ps.commandLineArguments[key].DefaultValue
}

func (ps *ProgramSettings) RegisterCommandLineSetting(cla CommandLineArgument) {
	ps.commandLineArguments[cla.Long] = cla
	ps.storage[cla.Long] = cla.DefaultValue
}

func (ps *ProgramSettings) ReadCommandLineSettings(pSettingsArray []string) {
	var wg sync.WaitGroup
	wg.Add(len(pSettingsArray))
	// TODO: Expand this at some point to search for command line arguments that need n values to work.
	for argCounter := 0; argCounter < len(pSettingsArray); argCounter++ {
		go ps.ReadSetting(pSettingsArray, argCounter, &wg)
	}
	wg.Wait()
}

func (ps *ProgramSettings) VerbosePrintln(output ...interface{}) {
	if (*ps).Get("verbose").(bool) {
		fmt.Println(output)
	}
}

func (ps *ProgramSettings) VerbosePrint(output ...interface{}) {
	if (*ps).Get("verbose").(bool) {
		fmt.Print(output)
	}
}

func (ps *ProgramSettings) VerbosePrintf(format string, a...interface{}) {
	if (*ps).Get("verbose").(bool) {
		fmt.Printf(format, a...)
	}
}

// "Private" Worker Functions
func (ps *ProgramSettings) ReadSetting(settings []string, settingOffset int, wg *sync.WaitGroup) () {
	defer wg.Done()

	if strings.HasPrefix(settings[settingOffset], "-") {
		for _, v := range (*ps).commandLineArguments {
			if "-"+v.Short == settings[settingOffset] || "--"+v.Long == settings[settingOffset] {
				if v.MultipleArguments {
					countOffset := settingOffset + 1
					arguments := make([]string, 0)

					for countOffset < len(settings) {
						if strings.HasPrefix(settings[countOffset], "-") {
							countOffset = len(settings)
						} else {
							arguments = append(arguments, settings[countOffset])
						}
						countOffset++
					}

					v.CommandHandler(arguments, ps)

				} else {
					v.CommandHandler(make([]string, 0), ps)
				}
			}
		}
	}

}
