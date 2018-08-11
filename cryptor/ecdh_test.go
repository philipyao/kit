package cryptor

import (
    "testing"
    "crypto/elliptic"
    "bytes"
)

func TestGoECDH(t *testing.T) {
    if false {
        return
    }
    curves := []elliptic.Curve{elliptic.P224(), elliptic.P256(), elliptic.P384(), elliptic.P521()}
    testGoECDH(curves, t)
}

func BenchmarkGoECDHP224(b *testing.B) {
    curves := []elliptic.Curve{elliptic.P224()}
    for i := 0; i < b.N; i++ {
        testGoECDH(curves, b)
    }
}

func BenchmarkGoECDHP256(b *testing.B) {
    curves := []elliptic.Curve{elliptic.P256()}
    for i := 0; i < b.N; i++ {
        testGoECDH(curves, b)
    }
}

func BenchmarkGoECDHP384(b *testing.B) {
    curves := []elliptic.Curve{elliptic.P384()}
    for i := 0; i < b.N; i++ {
        testGoECDH(curves, b)
    }
}

func BenchmarkGoECDHP521(b *testing.B) {
    curves := []elliptic.Curve{elliptic.P521()}
    for i := 0; i < b.N; i++ {
        testGoECDH(curves, b)
    }
}
func testGoECDH(curves []elliptic.Curve, tb testing.TB) {
    for _, curve := range curves {
        ECDH := NewGoECDH(curve)
        pub, priv, err := ECDH.GenerateKeys()
        if err != nil {
            tb.Fatal(err)
        }
        pub2, priv2, err := ECDH.GenerateKeys()
        if err != nil {
            tb.Fatal(err)
        }
        secret := ECDH.MakeSecret(pub2, priv)
        secret2 := ECDH.MakeSecret(pub, priv2)
        if !bytes.Equal(secret, secret2) {
            tb.Fatalf("secret mismatch: %v, %v", secret, secret2)
        }
        //tb.Logf("publen %v, privlen %v, secretlen %v", len(pub), len(priv), len(secret))
    }
}

//Mac
//Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
//8G mem
//machdep.cpu.core_count: 2
//machdep.cpu.thread_count: 4

//BenchmarkGoECDHP224-4   	     500	   3229132 ns/op
//BenchmarkGoECDHP256-4   	   10000	    169742 ns/op
//BenchmarkGoECDHP384-4   	     100	  16802830 ns/op
//BenchmarkGoECDHP521-4   	      50	  31517620 ns/op
