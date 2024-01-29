package file

import (
	"testing"
)

func TestBuild(t *testing.T) {
	type fields struct {
		DependenciesToInstall []string
		Script                string
		ShellPath             string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "TestBuild",
			fields: fields{
				DependenciesToInstall: []string{"curl", "wget"},
				Script:                "script.sh",
				ShellPath:             "/bin/bash",
			},
			want: `FROM alpine:latest as builder
RUN apk add --no-cache \
    curl \
    wget

COPY script.sh .
RUN chmod +x script.sh && mv script.sh /usr/bin/

FROM scratch

LABEL version="<version>"
LABEL description="<description>"
LABEL maintainer="<your name> <your email>"

COPY --from=builder /usr/bin/script.sh /usr/bin/

ENV SHELL=/bin/bash
WORKDIR /app

ENTRYPOINT ["/bin/bash", "/usr/bin/script.sh"]
`,
		},
	}
	for _, tt := range tests {
		d := &Dockerfile{
			DependenciesToInstall: tt.fields.DependenciesToInstall,
			Script:                tt.fields.Script,
			ShellPath:             tt.fields.ShellPath,
		}
		got, err := d.Build()
		if err != nil {
			t.Errorf("%q. Dockerfile.Build() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
		if (err != nil) != tt.wantErr {
			t.Errorf("%q. Dockerfile.Build() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("%q. Dockerfile.Build() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
