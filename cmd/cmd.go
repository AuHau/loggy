package cmd

import (
	"fmt"
	"github.com/auhau/loggy/store"
	"github.com/auhau/loggy/ui"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Option names
const (
	CONFIG_OPTION_NAME      = "config"
	BUFFER_SIZE_NAME        = "buffer-size"
	PARSE_PATTERN_NAME      = "pattern"
	PARSE_PATTERN_NAME_NAME = "pattern-name"
)

// Default values
const (
	DEFAULT_BUFFER_SIZE = 10000
)

var LONG_DESCRIPTION = fmt.Sprintf(`By default loggy reads from STDIN or you can specify file path to read the logs from specific file.

Configuration
-------------
All options that can be passed using CLI flags can be configured using config file or environment variables.
loggy looks for a config file ".loggy.toml" in current working directory and in $HOME folder. 
Moreover you can use environment variables with "LOGGY_" prefix, where for example flag "--pattern" would be "LOGGY_PATTERN" env. variable.
The order of precedence is: $HOME config > CWD config > --config config > Env. variables > CLI flags.

Parsing pattern
---------------
%s

Pattern names
-------------
In your config file you can create a [patterns] section where you can predefine your patterns using <name>="<pattern>" syntax and then use --pattern-name/-n flag to use it.

Filter
------
%s

Keyboard shortcuts
------------------
%s

%s

%s
`, ui.HelpParsingPatternText, ui.HelpFilterText, ui.HelpHomeText, ui.HelpNavigationText, ui.HelpInputsText)

var cfgFile string

// TODO: Add option to turn off Regex escaping
// TODO: Add option to set filter
// TODO: Add option to save the logs into file (as the logs are discarded when exceeding buffer size)
// TODO: Allow "infinite" buffer size?

// cmd represents the entry command
var cmd = &cobra.Command{
	Use:     "loggy [path to log file]",
	Short:   "Swiss knife for logs",
	Long:    LONG_DESCRIPTION,
	Args:    cobra.MaximumNArgs(1),
	Version: ui.Version,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			inputStream io.Reader
			err         error
		)

		bufferSize := viper.GetInt(BUFFER_SIZE_NAME)
		patter := viper.GetString(PARSE_PATTERN_NAME)
		patterName := viper.GetString(PARSE_PATTERN_NAME_NAME)

		if patterName != "" {
			patter, err = resolvePatterName(patterName)
			cobra.CheckErr(err)
		}

		if len(args) == 1 {
			file, err := os.Open(args[0])
			cobra.CheckErr(err)

			inputStream = file
		} else {
			inputStream = os.Stdin
		}

		uiApp, uiWriter, err := ui.Bootstrap(bufferSize, patter)
		cobra.CheckErr(err)

		go store.StartBuffering(inputStream, uiWriter, bufferSize)

		err = uiApp.Run()
		cobra.CheckErr(err)
	},
}

// resolvePatterName looks into configured `patterns` and if find one by the name it will use it as parse pattern for the logs
func resolvePatterName(name string) (string, error) {
	patterns := viper.GetStringMapString("patterns")
	chosenPattern, exists := patterns[name]

	if !exists {
		return "", fmt.Errorf("pattern with name '%s' does not exist", name)
	}

	return chosenPattern, nil
}

func init() {
	cobra.OnInitialize(initConfig)

	cmd.Flags().StringVar(&cfgFile, CONFIG_OPTION_NAME, "", "config file (default is $HOME/.loggy.yaml)")
	cmd.Flags().IntP(BUFFER_SIZE_NAME, "b", DEFAULT_BUFFER_SIZE, "number of lines that will be buffered")
	cmd.Flags().StringP(PARSE_PATTERN_NAME, "p", "", "parsing pattern see above for details")
	cmd.Flags().StringP(PARSE_PATTERN_NAME_NAME, "n", "", "use predefined pattern in config")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("toml")
	viper.SetConfigName(".loggy")

	// Home directory config
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(home)
	_ = viper.ReadInConfig() // if no config is found we don't care

	// Get current directory config
	viper.AddConfigPath(".")
	_ = viper.MergeInConfig() // if no config is found we don't care

	// Config file from the flag.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		cobra.CheckErr(viper.MergeInConfig())
	}

	viper.SetEnvPrefix("loggy")
	viper.AutomaticEnv() // read in environment variables that match

	cobra.CheckErr(viper.BindPFlags(cmd.Flags()))
}

func Execute() {
	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
