package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/util"
	"github.com/spf13/cobra"
)

var addtestCmd = &cobra.Command{
	Use:   "addtest TOKEN",
	Short: "Add test to task",
	Long:  `Add existing test case to current task. The command will register files with names TOKEN.in and TOKEN.out as a test case.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatalf("ERROR single argument is required for the command")
		}
		testToken := args[0]

		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}
		taskToken, err := model.DetermineCurrentTask(contest)
		if err != nil {
			log.Fatalf("ERROR can't determine current task: %s\n", err)
		}
		task, _ := contest.Tasks[taskToken]
		if util.ContainsString(&task.TestTokens, testToken) {
			log.Fatalf("ERROR test '%s' already added", testToken)
		}
		taskDir := filepath.Join(contest.RootDir, taskToken)
		inputPath := filepath.Join(taskDir, testToken+".in")
		outputPath := filepath.Join(taskDir, testToken+".out")
		if !util.PathExists(inputPath) {
			log.Fatalf("ERROR input file should exist at %s", inputPath)
		}
		if !util.PathExists(outputPath) {
			log.Fatalf("ERROR output file should exist at %s", outputPath)
		}
		task.TestTokens = append(task.TestTokens, testToken)
		contest.Tasks[taskToken] = task
		err = model.SaveContest(contest)
		if err != nil {
			log.Fatalf("ERROR failed to save contest metadata.")
		}
		fmt.Printf("Test %s is added to task %s\n", testToken, taskToken)
	},
}

func init() {
	RootCmd.AddCommand(addtestCmd)
}
