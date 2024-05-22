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
		killChromeProcesses()
	},
}

func init() {
	rootCmd.AddCommand(killChromeCmd)
}

func killChromeProcesses() {
	processes, err := process.Processes()
	if err != nil {
		panic(err)
	}
	for _, p := range processes {
		if isChromeProcess(p) {
			_ = p.Kill()
		}
	}
}

// isChromeProcess checks if a process is chrome/chromium
func isChromeProcess(process *process.Process) bool {
	name, _ := process.Name()
	if name == "" {
		return false
	}
	return strings.HasPrefix(strings.ToLower(name), "chromium")
}
