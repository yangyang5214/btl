package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yangyang5214/btl/pkg"
)

// cdzCmd represents the cdz command
var cdzCmd = &cobra.Command{
	Use:   "cdz",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		charging := pkg.NewChargingCDZ(orderCsv)
		err := charging.Run()
		if err != nil {
			panic(err)
		}
	},
}

var (
	orderCsv string
)

func init() {
	rootCmd.AddCommand(cdzCmd)
	cdzCmd.Flags().StringVarP(&orderCsv, "order_csv", "", "", "order csv file path")
}
