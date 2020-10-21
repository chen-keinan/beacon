package cli

import (
	"github.com/chen-keinan/beacon/internal/cli/commands"
	"github.com/chen-keinan/beacon/internal/logger"
	"github.com/chen-keinan/beacon/internal/startup"
	"github.com/chen-keinan/beacon/pkg/utils"
	"github.com/mitchellh/cli"
	"os"
	"strings"
)

var log = logger.GetLog()

//StartCli init beacon cli , folder , templates and etc
func StartCli() {
	err := utils.CreateHomeFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	err = utils.CreateBenchmarkFolderIfNotExist()
	if err != nil {
		panic(err)
	}
	filesData := startup.GenerateK8sBenchmarkFiles()
	err = startup.SaveBenchmarkFilesIfNotExist(filesData)
	if err != nil {
		panic(err)
	}
}

//InitCLI initialize beacon cli
func InitCLI(sa SanitizeArgs) {
	args := sa(os.Args[1:])
	app := cli.NewCLI("beacon", "1.0.0")
	// init cli folder and templates
	StartCli()
	app.Args = append(app.Args, args...)
	app.Commands = map[string]cli.CommandFactory{
		"a": func() (cli.Command, error) {
			return commands.NewK8sAudit(app.Args), nil
		},
	}
	status, err := app.Run()
	if err != nil {
		log.Console(err.Error())
	}
	os.Exit(status)
}

//ArgsSanitizer sanitize CLI arguments
var ArgsSanitizer SanitizeArgs = func(str []string) []string {
	args := make([]string, 0)
	if len(str) == 0 {
		args = append(args, "")
	}
	for _, arg := range str {
		arg = strings.Replace(arg, "--", "", -1)
		arg = strings.Replace(arg, "-", "", -1)
		c := arg[0:1]
		args = append(args, c)
	}
	return args
}

//SanitizeArgs sanitizer func
type SanitizeArgs func(str []string) []string
