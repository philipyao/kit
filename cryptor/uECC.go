package cryptor

/*
#include "uECC.h"
 */
import "C"
import (
    "unsafe"
    "errors"
)

type ECDH struct {
    curve C.uECC_Curve
}

func NewECDH() *ECDH {
    return &ECDH{
        curve: C.uECC_secp256k1(),
    }
}

func (ecdh *ECDH) GenerateKeys() (pub []byte, priv []byte, err error) {
    pub = make([]byte, 64)
    priv = make([]byte, 32)
    ret := C.uECC_make_key((*C.uint8_t)(unsafe.Pointer(&pub[0])),
        (*C.uint8_t)(unsafe.Pointer(&priv[0])),
        ecdh.curve)
    if int(ret) == 0 {
        err = errors.New("err C.uECC_make_key")
    }
    return
}

func (ecdh *ECDH) MakeSecret(pub []byte, priv []byte) []byte {
    secret := make([]byte, 32)
    ret := C.uECC_shared_secret((*C.uint8_t)(unsafe.Pointer(&pub[0])),
        (*C.uint8_t)(unsafe.Pointer(&priv[0])),
        (*C.uint8_t)(unsafe.Pointer(&secret[0])),
        ecdh.curve)
    if int(ret) == 0 {
        return nil
    }
    return secret
}