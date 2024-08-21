package cmd

import (
	"github.com/shirou/gopsutil/v3/process"
	"github.com/spf13/cobra"
	"strings"
)

// killChromeCmd represents the killChrome command
var killChromeCmd = &cobra.Command{
	Use:   "kill_chrome",
	Short: "kill all chrome process(chromium)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		KillChromeProcesses(all)
	},
}

var (
	all bool
)

func init() {
	rootCmd.AddCommand(killChromeCmd)
	killChromeCmd.Flags().BoolVarP(&all, "all", "a", false, "all is contains chrome")
}

func KillChromeProcesses(all bool) {
	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		if isChromeProcess(p, all) {
			_ = p.Kill()
		}
	}
}

// isChromeProcess checks if a process is chrome/chromium
func isChromeProcess(process *process.Process, all bool) bool {
	name, _ := process.Name()
	if name == "" {
		return false
	}
	lowerName := strings.ToLower(name)
	if all {
		return strings.Contains(lowerName, "chromium") || strings.Contains(lowerName, "chrome")
	}
	return strings.Contains(lowerName, "chromium")
}
