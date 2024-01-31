// EXPERIMENTAL
package dockerfile

import (
	"log"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type Dockerfile struct {
	Dockerfile string
}

func (d *Dockerfile) Parse() error {
	result, err := parser.Parse(strings.NewReader(d.Dockerfile))
	if err != nil {
		return err
	}
	log.Printf("AST: %+v\n", result.AST)
	for _, child := range result.AST.Children {
		log.Println("instruction comments: ", child.PrevComment)
		log.Println("instruction: ", child.Value)
		log.Println("original instruction: ", child.Original)
	}
	return nil
}
