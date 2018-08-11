package cryptor

import (
    "crypto/elliptic"
    "crypto/rand"
)

type GoECDH struct {
    elliptic.Curve
}

func NewGoECDH(curve elliptic.Curve) *GoECDH {
    if curve == nil {
        curve = elliptic.P224()
    }
    return &GoECDH{
        Curve: curve,
    }
}

func (ge *GoECDH) GenerateKeys() (pub []byte, priv []byte, err error) {
    p, x, y, e := elliptic.GenerateKey(ge.Curve, rand.Reader)
    if e != nil {
        err = e
        return
    }
    pub = elliptic.Marshal(ge.Curve, x, y)
    priv = p
    return
}

func (ge *GoECDH) MakeSecret(pub []byte, priv []byte) []byte {
    x, y := elliptic.Unmarshal(ge.Curve, pub)
    secret, _ := ge.Curve.ScalarMult(x, y, priv)
    return secret.Bytes()
}


