package prefix

import (
	"io"
	"strings"
	"testing"
)

type threeByteToABC struct {
	bytesRead int
}

func (r *threeByteToABC) Rewrite(in []byte, n int) ([]byte, bool) {
	var out []byte
	for i := 0; i < n; i++ {
		switch r.bytesRead + i {
		case 0:
			out = append(out, 'A')
		case 1:
			out = append(out, 'B')
		case 2:
			out = append(out, 'C')
		default:
			out = append(out, in[i])
		}
	}
	r.bytesRead += n
	return out, r.bytesRead > 3
}

func TestPrefixReader(t *testing.T) {
	rd, err := NewReader(strings.NewReader("1234567890"), &threeByteToABC{})
	if err != nil {
		t.Fatal(err)
	}
	bs, err := io.ReadAll(rd)
	if err != nil {
		t.Fatal(err)
	}
	if want, got := "ABC4567890", string(bs); want != got {
		t.Fatalf("want `%s` but got `%s`", want, got)
	}
}
