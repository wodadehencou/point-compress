package pointcompress

import (
	"crypto/elliptic"
	"math/big"
)

// CompressPoint compress a point to a byte string
func CompressPoint(curve elliptic.Curve, x, y *big.Int) []byte {
	if x.Sign() == 0 && y.Sign() == 0 {
		return []byte{0x00}
	}

	byteLen := (curve.Params().BitSize + 7) >> 3

	ret := make([]byte, 1+byteLen)
	if y.Bit(0) == 0 {
		ret[0] = 0x02
	} else {
		ret[0] = 0x03
	}

	xBytes := x.Bytes()
	copy(ret[1:], xBytes)
	return ret
}

// DecompressPoint decompress a byte string to a point
func DecompressPoint(curve elliptic.Curve, data []byte) (x, y *big.Int) {
	byteLen := (curve.Params().BitSize + 7) >> 3
	switch data[0] {
	case 0x00:
		if len(data) == 1 {
			return new(big.Int), new(big.Int)
		}
	case 0x02, 0x03:
		{
			if len(data) != 1+byteLen {
				return nil, nil
			}
			x = new(big.Int).SetBytes(data[1:])

			// x³ - 3x + b
			x3 := new(big.Int).Mul(x, x)
			x3.Mul(x3, x)

			threeX := new(big.Int).Lsh(x, 1)
			threeX.Add(threeX, x)

			x3.Sub(x3, threeX)
			x3.Add(x3, curve.Params().B)
			x3.Mod(x3, curve.Params().P)

			y = new(big.Int).ModSqrt(x3, curve.Params().P)
			if y == nil {
				x, y = nil, nil
			}
			if y.Bit(0) != uint(data[0]&0x01) {
				y.Sub(curve.Params().P, y)
			}
			return x, y
		}
	case 0x04:
		return elliptic.Unmarshal(curve, data)
	default:
		return nil, nil
	}
	return nil, nil
}
