package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/platforms"
	"github.com/spf13/cobra"
)

func determineRootDirectory(args []string) (string, error) {
	if len(args) == 2 {
		if root_dirname, err := filepath.Abs(args[1]); err != nil {
			return "", err
		} else {
			return root_dirname, nil
		}
	}
	if root_dirname, err := os.Getwd(); err != nil {
		return "", err
	} else {
		return root_dirname, nil
	}
}

func initContestDirectory(contest *model.Contest) error {
	/* Compose path for metadata and check it's not taken */
	root_file := model.GetRootFile(contest)
	if _, err := os.Stat(root_file); err == nil {
		return fmt.Errorf("File %s already exists -- please, remove to initialize the directory for the contest\n", root_file)
	}
	/* Create root directory if it doesn't exist */
	if _, err := os.Stat(contest.RootDir); os.IsNotExist(err) {
		err = os.MkdirAll(contest.RootDir, 0777)
		if err != nil {
			return err
		}
	}
	/* Create root file */
	if err := model.SaveContest(contest); err != nil {
		return err
	}
	/* Create subdirs for tasks */
	for token, _ := range contest.Tasks {
		path := filepath.Join(contest.RootDir, token)
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}

var initCmd = &cobra.Command{
	Use:   "init URL [DIRECTORY]",
	Short: "Initialize contest in directory",
	Long: `Initialize DIRECTORY with metadata of the contest specified by URL. Current directory is used when DIRECTORY is omitted. Non-existing directory will be created.

URLs of regular Codeforces rounds are supported.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 || len(args) > 2 {
			fmt.Printf("wrong number of arguments - %d\n", len(args))
			return
		}

		root_dirname, err := determineRootDirectory(args)
		if err != nil {
			fmt.Printf("ERROR can't determine working directory: %s\n", err)
			return
		}

		platform := platforms.FindPlatform(args[0])
		if platform == nil {
			fmt.Printf("ERROR unable to find platform for url %s\n", args[0])
			return
		}
		contest, err := platform.GetContest(args[0], root_dirname)
		if err != nil {
			fmt.Printf("ERROR can't fetch contest: %s\n", err)
			return
		}

		if err := initContestDirectory(contest); err != nil {
			fmt.Printf("ERROR can't init contest directory: %s\n", err)
			return
		}

		fmt.Printf("Root directory: %s\n", contest.RootDir)
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
