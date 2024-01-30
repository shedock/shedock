package insights

type InsightTopic int

const (
	DependencyNotFound InsightTopic = iota
	DependencyLightweight
)

func (i InsightTopic) String() string {
	switch i {
	case DependencyNotFound:
		return "dependency-not-found"
	default:
		return "unknown"
	}
}
