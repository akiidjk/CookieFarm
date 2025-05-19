package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ByteTheCookies/cookieclient/internal/config"
	"github.com/ByteTheCookies/cookieclient/internal/logger"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generate a new exploit template",
	Long:  `Generate a new exploit template for the CookieFarm client. This command initializes a structured exploit template file in your specified directory with all necessary components for immediate use.`,
	Run:   Create,
}

var (
	name string // Name of the exploit template
	path string // Path to save the exploit template
)

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the exploit template")
	createCmd.Flags().StringVarP(&path, "path", "p", "", "Path to save the exploit template")
	createCmd.MarkFlagRequired("name")
	createCmd.MarkFlagRequired("path")
}

func Create(cmd *cobra.Command, args []string) {
	logger.Log.Debug().Str("Exploit name", name).Str("Exploit path", path).Msg("Creating exploit template")
	if !strings.HasSuffix(name, ".py") {
		name = name + ".py"
	}
	final_path := filepath.Join(path, name)
	exploit_file, err := os.OpenFile(final_path, os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0o777)
	if err != nil {
		fmt.Println("Errore durante la creazione del file:", err.Error())
		return
	}
	exploit_file.Write(config.ExploitTemplate)
	defer exploit_file.Close()

	logger.Log.Info().Str("Exploit path", final_path).Msg("File creato con successo")
}
