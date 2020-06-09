package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/lunarway/shuttle/pkg/markdown"
	"github.com/spf13/cobra"
)

var showTitles bool

var (
	planHelpCmd = &cobra.Command{
		Use:   "help",
		Short: "Show help for the plan",
		Run: func(cmd *cobra.Command, args []string) {
			context := getProjectContext()
			readme := path.Join(context.LocalPlanPath, "README.md")

			if showTitles {
				titles, err := markdown.GetTitles(readme, context.LocalPlanPath)
				if err != nil {
					fmt.Printf("Error: %s", err)
					os.Exit(1)
				}
				//fmt.Printf("%+v", titles)
				printTitles(titles, 0)
			} else {
				rootCmd.Help()
				// os.Exit(1)
				// output, err := markdown.RenderFile(readme, context.LocalPlanPath)
				// if err != nil {
				// 	fmt.Printf("Error: %s", err)
				// 	os.Exit(1)
				// }
				// fmt.Println(output)
			}

		},
	}
)

func init() {
	planHelpCmd.Flags().BoolVar(&showTitles, "titles", false, "")
	rootCmd.SetHelpCommand(planHelpCmd)
}

func printTitles(titles []markdown.Title, level int) {
	for _, title := range titles {
		fmt.Printf("%s • %s (%v)\n", strings.Repeat(" ", level*2), title.Name, title.Level)
		printTitles(title.SubTitle, level+1)
	}
}
