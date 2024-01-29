package apk

import (
	"testing"
)

func TestDependsOn(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected []*PackageDependency
	}{
		{
			data: []byte(`
bash-5.2.15-r5 description:
The GNU Bourne Again shell

bash-5.2.15-r5 webpage:
https://www.gnu.org/software/bash/bash.html

bash-5.2.15-r5 installed size:
3328 KiB

bash-5.2.15-r5 depends on:
/bin/sh
so:libc.musl-aarch64.so.1
so:libreadline.so.8

bash-5.2.15-r5 provides:
cmd:bash=5.2.15-r5

bash-5.2.15-r5 has auto-install rule:

bash-5.2.15-r5 license:
GPL-3.0-or-later
`),
			expected: []*PackageDependency{
				{
					Name: "/bin/sh",
					Type: Binary,
				},
				{
					Name: "libc.musl-aarch64.so.1",
					Type: SharedLibrary,
				},
				{
					Name: "libreadline.so.8",
					Type: SharedLibrary,
				},
			},
		},
		{
			data: []byte(`
fish-3.6.1-r2 description:
Modern interactive commandline shell

fish-3.6.1-r2 webpage:
https://fishshell.com/

fish-3.6.1-r2 installed size:
10 MiB

fish-3.6.1-r2 depends on:
bc
/bin/sh
so:libc.musl-aarch64.so.1
so:libgcc_s.so.1
so:libintl.so.8
so:libncursesw.so.6
so:libpcre2-32.so.0
so:libstdc++.so.6

fish-3.6.1-r2 provides:
cmd:fish=3.6.1-r2
cmd:fish_indent=3.6.1-r2
cmd:fish_key_reader=3.6.1-r2

fish-3.6.1-r2 has auto-install rule:

fish-3.6.1-r2 license:
GPL-2.0-only`),
			expected: []*PackageDependency{
				{
					Name: "bc",
					Type: Binary,
				},
				{
					Name: "/bin/sh",
					Type: Binary,
				},
				{
					Name: "libc.musl-aarch64.so.1",
					Type: SharedLibrary,
				},
				{
					Name: "libgcc_s.so.1",
					Type: SharedLibrary,
				},
				{
					Name: "libintl.so.8",
					Type: SharedLibrary,
				},
				{
					Name: "libncursesw.so.6",
					Type: SharedLibrary,
				},
				{
					Name: "libpcre2-32.so.0",
					Type: SharedLibrary,
				},
				{
					Name: "libstdc++.so.6",
					Type: SharedLibrary,
				},
			},
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		parser := ApkParser{Data: tc.data}
		dependencies, err := parser.DependsOn()
		if err != nil {
			t.Fatal(err)
		}

		for i, dependency := range dependencies {
			if dependency.Name != tc.expected[i].Name {
				t.Fatalf("expected %s, got %s", tc.expected[i].Name, dependency.Name)
			}
			if dependency.Type != tc.expected[i].Type {
				t.Fatalf("expected %s, got %s", tc.expected[i].Type, dependency.Type)
			}
		}
	}
}

func TestProvides(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected []*ProviderDependency
	}{
		{
			data: []byte(`
bash-5.2.15-r5 description:
The GNU Bourne Again shell

bash-5.2.15-r5 webpage:
https://www.gnu.org/software/bash/bash.html

bash-5.2.15-r5 installed size:
3328 KiB

bash-5.2.15-r5 depends on:
/bin/sh
so:libc.musl-aarch64.so.1
so:libreadline.so.8

bash-5.2.15-r5 provides:
cmd:bash=5.2.15-r5

bash-5.2.15-r5 has auto-install rule:

bash-5.2.15-r5 license:
GPL-3.0-or-later
`),
			expected: []*ProviderDependency{
				{
					Name:    "bash",
					Version: "5.2.15-r5",
				},
			},
		},
		{
			data: []byte(`
fish-3.6.1-r2 description:
Modern interactive commandline shell

fish-3.6.1-r2 webpage:
https://fishshell.com/

fish-3.6.1-r2 installed size:
10 MiB

fish-3.6.1-r2 depends on:
bc
/bin/sh
so:libc.musl-aarch64.so.1
so:libgcc_s.so.1
so:libintl.so.8
so:libncursesw.so.6
so:libpcre2-32.so.0
so:libstdc++.so.6

fish-3.6.1-r2 provides:
cmd:fish=3.6.1-r2
cmd:fish_indent=3.6.1-r2
cmd:fish_key_reader=3.6.1-r2

fish-3.6.1-r2 has auto-install rule:

fish-3.6.1-r2 license:
GPL-2.0-only
`),
			expected: []*ProviderDependency{
				{
					Name:    "fish",
					Version: "3.6.1-r2",
				},
				{
					Name:    "fish_indent",
					Version: "3.6.1-r2",
				},
				{
					Name:    "fish_key_reader",
					Version: "3.6.1-r2",
				},
			},
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		parser := ApkParser{Data: tc.data}
		providers, err := parser.Provides()
		if err != nil {
			t.Fatal(err)
		}

		for i, provider := range providers {
			if provider.Name != tc.expected[i].Name {
				t.Fatalf("expected %s, got %s", tc.expected[i].Name, provider.Name)
			}
			if provider.Version != tc.expected[i].Version {
				t.Fatalf("expected %s, got %s", tc.expected[i].Version, provider.Version)
			}
		}
	}
}
