package programSettings

import (
	"os"
	"strings"
	"sync"
)

// Structs
type ProgramSettings struct {
	storage              map[string]interface{}
	commandLineArguments map[string]CommandLineArgument
}

type CommandLineArgument struct {
	short             string
	long              string
	defaultValue	  interface{}
	multipleArguments bool
	maxArgumentParam  int
}

// Factory Functions
func MakeProgramSettings(defaultValues *map[string]interface{}) ProgramSettings {
	var valueStorage map[string]interface{}

	if defaultValues != nil {
		valueStorage = *defaultValues
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

func (ps *ProgramSettings) Reset(key string) {
	ps.storage[key] = ps.commandLineArguments[key].defaultValue
}

func (ps *ProgramSettings) RegisterCommandLineSetting(key string, cla CommandLineArgument) {
	ps.commandLineArguments[key] = cla
	ps.storage[key] = cla.defaultValue
}

func (ps *ProgramSettings) ReadCommandLineSettings(pSettingsArray *[]string) sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(len(*pSettingsArray))

	// TODO: Expand this at some point to search for command line arguments that need n values to work.
	for argCounter := 0; argCounter < len(*pSettingsArray); argCounter++ {
		go ps.ReadSetting(pSettingsArray, argCounter, wg)
	}

	return wg
}

// "Private" Worker Functions

func (ps *ProgramSettings) ReadSetting(settings *[]string, settingOffset int, wg sync.WaitGroup) (bool, string) {

	err := ""
	retVal := false

	defer wg.Done()

	if strings.HasPrefix((*settings)[settingOffset], "-") {
		for k, v := range ps.commandLineArguments {
			if v.short == ("-"+(*settings)[settingOffset]) || v.long == "--"+(*settings)[settingOffset] {
				// TODO: Add function to find proper value
				ps.storage[k] = true
				retVal = true
			} else {
				err = "Command line argument not found"
			}
		}
	}

	return retVal, err
}
