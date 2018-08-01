package cryptor

/*
#cgo CFLAGS: -I .

#include <stdlib.h>
#include <xxtea.h>
*/
import "C"

import (
    "unsafe"
    "errors"
)

// void * xxtea_encrypt(const void * data, size_t len, const void * key, size_t * out_len);
// void * xxtea_decrypt(const void * data, size_t len, const void * key, size_t * out_len);

func XXTeaEncrypt(data []byte, key []byte) ([]byte, error) {
    dataLen := C.size_t(len(data))
    cdataIn := unsafe.Pointer(&data[0])
    ckey := unsafe.Pointer(&key[0])

    var coutLen C.size_t
    coutLen = 0

    coutData := C.xxtea_encrypt(cdataIn, dataLen, ckey, &coutLen)
    realLen := int(coutLen)
    if coutData == nil || realLen == 0 {
        return nil, errors.New("enc failed")
    }
    dataPrt := unsafe.Pointer(coutData)
    defer C.free(dataPrt)

    outData := C.GoBytes(dataPrt, C.int(realLen))
    return outData, nil
}

func XXTeaDecrypt(data []byte, key []byte) ([]byte, error) {
    dataLen := C.size_t(len(data))
    cdataIn := unsafe.Pointer(&data[0])
    ckey := unsafe.Pointer(&key[0])

    var coutLen C.size_t
    coutLen = 0

    coutData := C.xxtea_decrypt(cdataIn, dataLen, ckey, &coutLen)
    realLen := int(coutLen)
    if coutData == nil || realLen == 0 {
        return nil, errors.New("dec failed")
    }

    dataPrt := unsafe.Pointer(coutData)
    defer C.free(dataPrt)

    outData := C.GoBytes(dataPrt, C.int(realLen))
    return outData, nil
}

