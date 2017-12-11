package pointcompress

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func TestCompress(t *testing.T) {
	curves := []elliptic.Curve{elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()}

	for _, c := range curves {

		_, x, y, err := elliptic.GenerateKey(c, rand.Reader)
		if err != nil {
			t.Error(err)
			return
		}

		compressed := CompressPoint(c, x, y)
		xx, yy := DecompressPoint(c, compressed)

		if xx == nil {
			t.Error("failed to decompress")
			break
		}
		if xx.Cmp(x) != 0 || yy.Cmp(y) != 0 {
			t.Error("Decompress returned different values")
			break
		}
	}
}
