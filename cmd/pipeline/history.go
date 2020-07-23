package pipeline

import (
	"net/url"

	"github.com/iterum-provenance/cli/manager"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(historyCmd)
	historyCmd.PersistentFlags().StringVarP(&DaemonURL, "url", "u", "http://localhost:3000", "URL at which the daemon can be reached")
}

// Gather more detailed status information for a specific pipeline
var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "List information of current and historic pipeline runs the daemon knows about",
	Long:  `List information of current and historic pipeline runs the daemon knows about`,
	Args:  cobra.NoArgs,
	Run:   historyRun,
}

func historyRun(cmd *cobra.Command, args []string) {
	url, err := url.Parse(DaemonURL)
	if err != nil {
		log.Fatalln(err)
	}
	err = manager.History(url)
	if err != nil {
		log.Fatalln(err)
	}
}
