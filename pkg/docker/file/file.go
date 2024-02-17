package file

import "fmt"

type Dockerfile struct {
	Dependencies          Dependencies
	DependenciesToInstall []string
	// Script is the path of the script in the second layer
	Script string
	// ShellPath is the path of the shell in the second layer
	ShellPath string
}

const (
	FirstLayerBaseImage  string = "alpine:latest"
	SecondLayerBaseImage string = "scratch"
	FirstLayerAlias      string = "builder"
)

func (d *Dockerfile) FirstLayer() (string, error) {
	var install string

	base := fmt.Sprintf("FROM %s as %s\n", FirstLayerBaseImage, FirstLayerAlias)

	if len(d.DependenciesToInstall) > 0 {
		install = "RUN apk add --no-cache \\\n"

		for depCount, dep := range d.DependenciesToInstall {
			if depCount == len(d.DependenciesToInstall)-1 {
				install += fmt.Sprintf("    %s\n", dep)
				break
			} else {
				install += fmt.Sprintf("    %s \\\n", dep)
			}
		}
	}

	install += "\n"

	install += fmt.Sprintf("COPY %s .\n", d.Script)
	install += fmt.Sprintf("RUN chmod +x %s && mv %s /usr/bin/\n", d.Script, d.Script)
	return base + install, nil
}

func (d *Dockerfile) SecondLayer() (string, error) {
	base := fmt.Sprintf("\nFROM %s\n", SecondLayerBaseImage)
	base += labels()

	copyInstructionSet := d.generateCopyInstructionSet()
	base += copyInstructionSet
	envs := d.Envs()
	base += envs

	base += "WORKDIR /app\n"
	base += d.Entrypoint()

	return base, nil
}

func (d *Dockerfile) Envs() string {
	var envs string
	if d.ShellPath != "" {
		envs += fmt.Sprintf("\nENV SHELL=%s\n", d.ShellPath)
	}
	return envs
}

func (d *Dockerfile) Entrypoint() string {
	return fmt.Sprintf("\nENTRYPOINT [\"%s\", \"/usr/bin/%s\"]\n", d.ShellPath, d.Script)
}

func (d *Dockerfile) Build() (string, error) {
	firstLayer, err := d.FirstLayer()
	if err != nil {
		return "", err
	}
	secondLayer, err := d.SecondLayer()
	if err != nil {
		return "", err
	}

	return firstLayer + secondLayer, nil
}

func labels() string {
	lables := `
LABEL description="<description>"
# Add your name & email here
LABEL maintainer="<your name> <your email>"
`
	return lables
}

func (d *Dockerfile) generateCopyInstructionSet() string {
	var copyInstructionSet string

	// first copy the user script
	copyInstructionSet += fmt.Sprintf("\nCOPY --from=%s /usr/bin/%s /usr/bin/\n", FirstLayerAlias, d.Script)

	if len(d.Dependencies.Bin) != 0 {
		for _, dep := range d.Dependencies.Bin {
			copyInstructionSet += fmt.Sprintf("## Required By: %s\n", dep.Requiredby)
			copyInstructionSet += fmt.Sprintf("COPY --from=%s %s %s\n", FirstLayerAlias, dep.FromPath, dep.ToPath)
		}
	}

	if len(d.Dependencies.Lib) != 0 {
		for _, dep := range d.Dependencies.Lib {
			copyInstructionSet += fmt.Sprintf("## Required By: %s\n", dep.Requiredby)
			copyInstructionSet += fmt.Sprintf("COPY --from=%s %s %s\n", FirstLayerAlias, dep.FromPath, dep.ToPath)
		}
	}

	return copyInstructionSet
}
