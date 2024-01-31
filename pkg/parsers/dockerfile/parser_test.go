package dockerfile

import "testing"

func TestParse(t *testing.T) {
	d := &Dockerfile{
		Dockerfile: `
FROM alpine:3.13.5
RUN apk add --no-cache fish
# comments in dockerfile
ENTRYPOINT ["fish"]
`}
	err := d.Parse()

	if err != nil {
		t.Errorf("Error parsing Dockerfile: %v", err)
	}

}
