package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"runtime/debug"

	"github.com/ArchiveNetwork/wgcf-cli/constant"
)

type Actions struct {
	Help      bool
	Register  bool
	Version   bool
	TeamToken string
	Bind      bool
	UnBind    bool
	Cancel    bool
	License   string
	FileName  string
	Name      string
	Generate  string
	Plus      bool
	Update    bool
	Convert   bool
}

func ParseCommandLine() Actions {
	var action Actions

	flag.BoolVar(&action.Help, "h", false, "")
	flag.BoolVar(&action.Help, "help", false, "")

	flag.BoolVar(&action.Register, "r", false, "")
	flag.BoolVar(&action.Register, "register", false, "")

	flag.BoolVar(&action.Version, "V", false, "")
	flag.BoolVar(&action.Version, "version", false, "")

	flag.StringVar(&action.TeamToken, "t", "", "")
	flag.StringVar(&action.TeamToken, "token", "", "")

	flag.BoolVar(&action.Bind, "b", false, "")
	flag.BoolVar(&action.Bind, "bind", false, "")

	flag.BoolVar(&action.UnBind, "U", false, "")
	flag.BoolVar(&action.UnBind, "unbind", false, "")

	flag.BoolVar(&action.Cancel, "C", false, "")
	flag.BoolVar(&action.Cancel, "cancel", false, "")

	flag.StringVar(&action.FileName, "f", "", "")
	flag.StringVar(&action.FileName, "file", "", "")

	flag.StringVar(&action.License, "l", "", "")
	flag.StringVar(&action.License, "license", "", "")

	flag.StringVar(&action.Name, "n", "", "")
	flag.StringVar(&action.Name, "name", "", "")

	flag.StringVar(&action.Generate, "g", "", "")
	flag.StringVar(&action.Generate, "generate", "", "")

	flag.BoolVar(&action.Plus, "p", false, "")
	flag.BoolVar(&action.Plus, "plus", false, "")

	flag.BoolVar(&action.Update, "u", false, "")
	flag.BoolVar(&action.Update, "update", false, "")

	flag.BoolVar(&action.Convert, "c", false, "")
	flag.BoolVar(&action.Convert, "covert", false, "")

	flag.Usage = func() {
		Help()
	}

	flag.Parse()

	if action.License != "" {
		expectedPattern := `^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$`

		match, err := regexp.MatchString(expectedPattern, action.License)
		if err != nil {
			panic(err)
		}

		if !match {
			panic("License should be something matchs: ^[0-9A-Za-z]{8}-[0-9A-Za-z]{8}-[0-9A-Za-z]{8}$")
		}
	}

	return action
}

func Help() {
	Version()
	fmt.Fprintf(os.Stderr, `Usage:	%s [Options]
Options:    -h/--help                             Show help
            -V/--version                          Show version
            -f/--file [string]                    Configuration file (default "wgcf.ini")
            -r/--register                         Register an account
            -t/--token [string]                   Team token (must be used with -r/--register)
            -b/--bind                             Get the account binding devices
            -n/--name [string]                    Change the device name
            -l/--license [string]                 Change the license
            -U/--unbind                           Unbind a device from the account
            -C/--cancel                           Cancel the account
            -g/--generate [sing-box/wg/xray]      Generate a [sing-box/wg/xray] configuration file
            -p/--plus                             Recharge your account indefinitely
            -u/--update                           Update the configuration file
            -c/--convert                          Convert the configuration file
`, os.Args[0])
}

func Version() {
	var revision string
	debugInfo, loaded := debug.ReadBuildInfo()

	if loaded {
		for _, setting := range debugInfo.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			}
		}
	}
	if revision == "" {
		revision = "unknow"
	}
	fmt.Fprintf(os.Stderr, `wgcf-cli version `+"\033[1;35m"+constant.Version+"\033[0m"+` 
Environment: %s %s/%s
`+"Revision: \033[1;36m"+revision+"\033[0m"+`
`, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
