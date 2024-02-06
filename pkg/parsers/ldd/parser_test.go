package ldd

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	lddParser := LddParser{
		Data: []byte(`
		/lib/ld-musl-aarch64.so.1 (0xffffbe6b0000)
		libcurl.so.4 => /usr/lib/libcurl.so.4 (0xffffbe5bc000)
		libz.so.1 => /lib/libz.so.1 (0xffffbe58b000)
		libc.musl-aarch64.so.1 => /lib/ld-musl-aarch64.so.1 (0xffffbe6b0000)
		libnghttp2.so.14 => /usr/lib/libnghttp2.so.14 (0xffffbe54a000)
		libidn2.so.0 => /usr/lib/libidn2.so.0 (0xffffbe509000)
		libssl.so.3 => /lib/libssl.so.3 (0xffffbe465000)
		libcrypto.so.3 => /lib/libcrypto.so.3 (0xffffbe0af000)
		libbrotlidec.so.1 => /usr/lib/libbrotlidec.so.1 (0xffffbe08e000)
		libunistring.so.5 => /usr/lib/libunistring.so.5 (0xffffbdecc000)
		libbrotlicommon.so.1 => /usr/lib/libbrotlicommon.so.1 (0xffffbde9b000)
`),
	}
	dependencies := lddParser.Parse()
	fmt.Println(dependencies)
	// if len(dependencies) != 1 {
	// 	t.Errorf("Expected 1 dependency, got %d", len(dependencies))
	// }
	// if dependencies[0].Binary != "/bin/sh" {
	// 	t.Errorf("Expected /bin/sh, got %s", dependencies[0].Binary)
	// }
	// if len(dependencies[0].Dependencies) != 1 {
	// 	t.Errorf("Expected 1 dependency, got %d", len(dependencies[0].Dependencies))
	// }
	// if dependencies[0].Dependencies[0].Name != "libc.musl-aarch64.so.1" {
	// 	t.Errorf("Expected libc.musl-aarch64.so.1, got %s", dependencies[0].Dependencies[0].Name)
	// }
	// if dependencies[0].Dependencies[0].IsSymlink != false {
	// 	t.Errorf("Expected false, got %t", dependencies[0].Dependencies[0].IsSymlink)
	// }

}
