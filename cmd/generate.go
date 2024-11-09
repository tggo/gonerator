package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "create build of the data",
	Long:  `Make static data for platform`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("generate called %q with %v", cmd.Use, args)
	},
}

func init() {
	generateCmd.AddCommand(courseGenerateCmd)

	RootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&courseID, "course", "c", "", "Course Uniq ID to generate")
}
