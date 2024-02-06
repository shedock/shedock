package insights

type Insight struct {
	Topic             InsightTopic
	Type              string
	CliComment        string
	DockerfileComment string
}

func (i Insight) Generate(deps []string) []*Insight {
	var insights []*Insight

	for _, dep := range deps {
		switch dep {
		case "xdg-open":
			insights = append(insights, &Insight{
				Topic:             DependencyNotSupported,
				Type:              DependencyNotFound.String(),
				CliComment:        "",
				DockerfileComment: "",
			})
		case "notify-send":
			insights = append(insights, &Insight{
				Topic:             DependencyNotSupported,
				Type:              DependencyNotFound.String(),
				CliComment:        "",
				DockerfileComment: "",
			})
		}
	}
	return insights
}
