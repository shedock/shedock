package insights

type InsightTopic int

const (
	// DependencyNotFound is the topic for when a dependency is not found
	DependencyNotFound InsightTopic = iota
	// DependencyNotSupported is the topic for when a dependency is not supported
	DependencyNotSupported
	// DependencyAlternative is the topic for when a dependency has an alternative
	DependencyAlternative
)

func (i InsightTopic) String() string {
	switch i {
	case DependencyNotFound:
		return "dependency-not-found"
	case DependencyNotSupported:
		return "dependency-not-supported"
	case DependencyAlternative:
		return "dependency-alternative"
	default:
		return "unknown"
	}
}

var NOT_SUPPORTED_COMMANDS []string = []string{
	"xdg-open",
	"notify-send",
	"powershell.exe",
	"start",
	"terminal-notifier",
	"growlnotify",
	"kdialog",
	"notify",
}

var SPECIAL_COMMANDS map[string]string = map[string]string{
	"tput": "ncurses",
	"gawk": "coreutils",
}

var ALTERNATIVE_COMMANDS map[string]string = map[string]string{}
