package insights

type InsightTopic int

const (
	DependencyNotFound InsightTopic = iota
	DependencyNotSupported
)

func (i InsightTopic) String() string {
	switch i {
	case DependencyNotFound:
		return "dependency-not-found"
	case DependencyNotSupported:
		return "dependency-not-supported"
	default:
		return "unknown"
	}
}

var NOT_SUPPORTED_COMMANDS []string = []string{
	"xdg-open",
	"notify-send",
	"powershell.exe",
	"start",
}

var SPECIAL_COMMANDS map[string]string = map[string]string{
	"tput": "ncurses",
	"gawk": "coreutils",
}
