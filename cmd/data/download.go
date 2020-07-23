package data

import (
	"github.com/prometheus/common/log"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

var downloadCmd = &cobra.Command{
	Use:   "download [selector]",
	Short: "Download the data of a certain commit",
	Long:  `Dowloads the actual data of a commit into the specified folder with selector on filename`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return errNotEnoughArgs
		}
		if len(args) == 0 || isValidSelector(args[0]) {
			return nil
		}
		return errInvalidArgs(args[0])
	},
	Run: downloadRun,
}

func downloadRun(cmd *cobra.Command, args []string) {
	log.Infoln("`iterum data download is supposed to download data from the daemon at request`")
}
