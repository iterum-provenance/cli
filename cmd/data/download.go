package data

import (
	"os"
	"regexp"

	"github.com/iterum-provenance/cli/idv"

	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().BoolVarP(&NoPrompt, "no-prompt", "Y", false, "Do not ask to reassure before downloading the actual files")
	downloadCmd.PersistentFlags().BoolVarP(&ShowFiles, "show-files", "s", false, "Show a list of to be downloaded")

}

// Gather more detailed status information for a specific pipeline
var downloadCmd = &cobra.Command{
	Use:   "download [selector] [dirpath]",
	Short: "Download (a) file(s) from a commit of a dataset",
	Long:  `Downloads (a) file(s) from the set of files of a commit stored in a dataset`,
	Args:  cobra.ExactArgs(2),
	Run:   downloadRun,
}

func downloadRun(cmd *cobra.Command, args []string) {
	info, err := os.Stat(args[1])
	if err != nil {
		log.Fatalln(err)
	}
	if !info.IsDir() {
		log.Fatalln("Passed dirpath arg is not a directory")
	}
	selector, err := regexp.Compile(args[0])
	if err != nil {
		log.Fatalln(err)
	}
	err = idv.Download(selector, args[1], ShowFiles, !NoPrompt)
	if err != nil {
		log.Fatalln(err)
	}
}
