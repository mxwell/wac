package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mxwell/wac/model"
	"github.com/mxwell/wac/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ExecMethod struct {
	command string
}

type Outcome struct {
	exec_time      time.Duration
	output_differs bool
}

var ExecMethodByName = map[string]*ExecMethod{}
var ExecMethodName string
var TheMethod *ExecMethod
var SolutionName string
var UseStdStreams bool

func readExecConfig() {
	methods := viper.GetStringMap("RunMethods")
	for name, _ := range methods {
		subtree := viper.Sub("RunMethods." + name)
		command := subtree.GetString("Command")
		ExecMethodByName[name] = &ExecMethod{command}
	}
}

func getSolutionCommand() *exec.Cmd {
	commandLine := strings.Replace(TheMethod.command, "$OUTPUT", SolutionName, -1)
	tokens := strings.Split(commandLine, " ")
	if len(tokens) == 1 {
		return exec.Command(tokens[0])
	} else if len(tokens) > 1 {
		return exec.Command(tokens[0], tokens[1:]...)
	} else {
		panic("Empty command to run solution!")
	}
}

func doRun(inputPath string, resultPath string) (error, time.Duration) {
	var stderr bytes.Buffer

	command := getSolutionCommand()
	if len(inputPath) > 0 {
		inputReader, err := os.Open(inputPath)
		if err != nil {
			return fmt.Errorf("failed to open test input: %s", err), 0
		}
		defer inputReader.Close()
		command.Stdin = inputReader
	} else {
		command.Stdin = os.Stdin
	}
	if len(resultPath) > 0 {
		resultWriter, err := os.OpenFile(resultPath, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("failed to open file to write output: %s", err), 0
		}
		defer resultWriter.Close()
		command.Stdout = resultWriter
	} else {
		command.Stdout = os.Stdout
	}
	command.Stderr = &stderr

	start := time.Now()
	err := command.Run()
	elapsed := time.Since(start)

	if stderr.Len() > 0 {
		fmt.Println("<stderr>")
		stderr.WriteTo(os.Stdout)
	}

	return err, elapsed
}

func readWholeLine(r *bufio.Reader) (string, error) {
	result := make([]byte, 0)
	for {
		if bytes, isPrefix, err := r.ReadLine(); err == nil {
			result = append(result, bytes...)
			if !isPrefix {
				break
			}
		} else {
			return "", err
		}
	}
	return strings.TrimRight(string(result), " "), nil
}

/* true if result is different from what was expected */
func checkOutput(expectedPath string, resultPath string) (bool, error) {
	a0, err := os.Open(expectedPath)
	if err != nil {
		return true, fmt.Errorf("failed to open file with expected output: %s", err)
	}
	defer a0.Close()
	a := bufio.NewReader(a0)

	b0, err := os.Open(resultPath)
	if err != nil {
		return true, fmt.Errorf("failed to open file with solution output: %s", err)
	}
	defer b0.Close()
	b := bufio.NewReader(b0)

	adone := false
	bdone := false
	for !adone && !bdone {
		aline, aerr := readWholeLine(a)
		bline, berr := readWholeLine(b)
		if aerr == io.EOF {
			adone = true
			aerr = nil
		}
		if berr == io.EOF {
			bdone = true
			berr = nil
		}
		if aerr != nil {
			return true, fmt.Errorf("error when reading from '%s': %s", expectedPath, aerr)
		}
		if berr != nil {
			return true, fmt.Errorf("error when reading from '%s': %s", resultPath, berr)
		}
		if aline != bline {
			return true, nil
		}
	}
	return false, nil
}

func runSingleTest(taskDir string, testToken string) (*Outcome, error) {
	testPathPrefix := filepath.Join(taskDir, testToken)
	inputPath := testPathPrefix + ".in"
	outputPath := testPathPrefix + ".out"
	resultPath := testPathPrefix + ".result"

	err, elapsed := doRun(inputPath, resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to run solution: %s", err)
	}

	diff, err := checkOutput(outputPath, resultPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check output: %s", err)
	}
	return &Outcome{elapsed, diff}, nil
}

var runCmd = &cobra.Command{
	Use:   "run [TOKEN1 TOKEN2 ...]",
	Short: "Run built solution on test cases",
	Long: `Run built solution on test cases. Set and order of test cases could be specified in command arguments as test tokens separated by spaces. If no arguments are given, then all available tests are used.`,
	Run: func(cmd *cobra.Command, args []string) {
		readExecConfig()
		if len(ExecMethodName) == 0 {
			ExecMethodName = viper.GetString("DefaultRunMethod")
		}
		var ok bool
		if TheMethod, ok = ExecMethodByName[ExecMethodName]; !ok {
			log.Fatalf("ERROR exec method '%s' not found in config\n", ExecMethodName)
		}
		if len(SolutionName) == 0 {
			SolutionName = viper.GetString("SolutionName")
		}
		contest, err := model.LocateContest()
		if err != nil {
			log.Fatalf("ERROR %s\n", err)
		}
		taskToken, err := model.DetermineCurrentTask(contest)
		if err != nil {
			log.Fatalf("ERROR can't determine current task: %s\n", err)
		}
		task, ok := contest.Tasks[taskToken]
		if !ok {
			log.Fatalf("ERROR contest have no task for the working directory")
		}
		if UseStdStreams {
			if len(args) > 0 {
				log.Fatalf("ERROR test tokens are now allowed when stdin/stdout are used")
			}
			err, _ := doRun("", "")
			if err != nil {
				log.Fatalf("ERROR failed to run solution: %s", err)
			}
			return
		}
		if len(task.TestTokens) == 0 {
			fmt.Println("No tests.")
			return
		}
		for _, testToken := range args {
			if !util.ContainsString(&task.TestTokens, testToken) {
				log.Fatalf("ERROR test with token '%s' not found", testToken)
			}
		}
		var selection *[]string
		if len(args) > 0 {
			selection = &args
		} else {
			selection = &task.TestTokens
		}
		for _, testToken := range *selection {
			fmt.Printf("[%s] ... ", testToken)
			outc, err := runSingleTest(filepath.Join(contest.RootDir, taskToken), testToken)
			if err != nil {
				log.Fatalf("ERROR failed to run test '%s': %s", testToken, err)
			}
			var msg string
			if outc.output_differs {
				msg = "Differs"
			} else {
				msg = "Ok"
			}
			fmt.Printf("%s -- %dms\n", msg, int(outc.exec_time/1000000))
			if outc.output_differs {
				break
			}
		}
	},
}

func init() {
	runCmd.Flags().StringVarP(&ExecMethodName, "with", "w", "", "Execution method name, like gcc (default is set in config under DefaultRunMethod)")
	runCmd.Flags().StringVarP(&SolutionName, "solution", "s", "", "Built solution name, like 'main' (default is set in config under SolutionName)")
	runCmd.Flags().BoolVarP(&UseStdStreams, "interactive", "i", false, "Interactive mode: use stdin and stdout instead of files")
	RootCmd.AddCommand(runCmd)
}
