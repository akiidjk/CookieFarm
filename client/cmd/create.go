package cmd

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/filesystem"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate a new exploit template",
	Long: `Generate a new exploit template for the CookieFarm client.
	This command initializes a structured exploit template file in your specified directory with all necessary components for immediate use.`,
	Run: Create,
}

var name string // Name of the exploit template

func init() {
	RootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the exploit template")
	createCmd.MarkFlagRequired("name")
}

func Create(cmd *cobra.Command, args []string) {
	path, err := filesystem.ExpandTilde(config.DefaultConfigPath)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error expanding path")
		return
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.Log.Warn().Msg("Default exploit path not exists... Creating it")
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger.Log.Error().Err(err).Msg("Error creating exploit path")
			return
		}
	}

	logger.Log.Debug().Str("Exploit name", name).Msg("Creating exploit template")

	name, err = filesystem.NormalizeNamePathExploit(name)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error normalizing exploit name")
		return
	}

	if filesystem.IsPath(name) {
		path = name
	} else {
		path = filepath.Join(path, name)
	}

	exploitFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o777)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Error creating exploit file")
		return
	}
	exploitFile.Write(config.ExploitTemplate)
	defer exploitFile.Close()

	logger.Log.Info().Str("Exploit path", path).Msg("File created successfully")
}
