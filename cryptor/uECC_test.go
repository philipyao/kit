package cryptor

import (
    "testing"
    "encoding/hex"
)

func TestECDH(t *testing.T) {
    if true {
        return
    }
    ecdh := NewECDH()
    pub, priv, err := ecdh.GenerateKeys()
    if err != nil {
        t.Fatal(err)
    }
    secret := ecdh.MakeSecret(pub, priv)
    t.Logf("pub %v, priv %v, secret %v",
        hex.EncodeToString(pub), hex.EncodeToString(priv), hex.EncodeToString(secret))
}

func BenchmarkECDH(b *testing.B) {
    ecdh := NewECDH()
    for i := 0; i < b.N; i++ {
        pub, priv, err := ecdh.GenerateKeys()
        if err != nil {
            b.Fatal(err)
        }
        ecdh.MakeSecret(pub, priv)
    }
    //BenchmarkECDH-4         	    3000	    533496 ns/op
}
