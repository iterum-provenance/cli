package pipeline

import (
	"net/url"
	"os"
	"path"

	"github.com/iterum-provenance/cli/util"

	"github.com/iterum-provenance/cli/pipeline"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.PersistentFlags().StringVarP(&DaemonURL, "url", "u", "http://localhost:3000", "URL at which the daemon can be reached")
}

// Gather more detailed status information for a specific pipeline
var downloadCmd = &cobra.Command{
	Use:   "download [pipeline-hash] [filename] [dirpath]",
	Short: "Download a result file of a dataset",
	Long:  `Downloads a file from the result set of a pipleline execution into the passed folder`,
	Args:  cobra.ExactArgs(3),
	Run:   downloadRun,
}

func downloadRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(DaemonURL)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := os.Stat(args[2])
	if err != nil {
		log.Fatalln(err)
	}
	if !info.IsDir() {
		log.Fatalln("Passed path is not a directory")
	}
	if util.FileExists(path.Join(args[2], args[1])) {
		log.Fatalln("File with that name already exists in that location")
	}
	err = pipeline.Download(args[0], args[1], args[2], url)
	if err != nil {
		log.Fatalln(err)
	}
}
