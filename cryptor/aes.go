package cryptor

import (
    "bytes"
    "crypto/cipher"
    "crypto/aes"
    "crypto/rand"
    "io"
)

//cipher feedback mode
func AesCFBEncrypt(text, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the cipherText.
    cipherText := make([]byte, aes.BlockSize + len(text))
    iv := cipherText[:aes.BlockSize]
    //生成随机 iv 值
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    stream := cipher.NewCFBEncrypter(block, iv)
    //原文和密文长度要一致
    stream.XORKeyStream(cipherText[aes.BlockSize:], text)
    //iv向量和加密后的密文一起返回
    return cipherText, nil
}

func AesCFBDecrypt(cipherText, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(cipherText) < aes.BlockSize {
        panic("ciphterText too short")
    }
    //密文前一部分为初始化向量 iv，无需保密
    iv := cipherText[:aes.BlockSize]
    data := cipherText[aes.BlockSize:]

    stream := cipher.NewCFBDecrypter(block, iv)

    // XORKeyStream can work in-place if the two arguments are the same.
    dst := make([]byte, len(data))
    stream.XORKeyStream(dst, data)
    return dst, nil
}

// cipher block chaining mode
//CBC 需要补齐不足的块
func AesCBCEncrypt(text, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    text = PKCS5Padding(text, aes.BlockSize)
    // text = ZeroPadding(text, aes.BlockSize)

    // The IV needs to be unique, but not secure. Therefore it's common to
    // include it at the beginning of the cipherText.
    cipherText := make([]byte, aes.BlockSize + len(text))
    iv := cipherText[:aes.BlockSize]
    //生成随机 iv 值
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        panic(err)
    }

    mode := cipher.NewCBCEncrypter(block, iv)
    //原文和密文长度相同
    mode.CryptBlocks(cipherText[aes.BlockSize:], text)
    //iv 和加密后的密文一起返回
    return cipherText, nil
}

//CBC 需要补齐不足的块
func AesCBCDecrypt(cipherText, key []byte) ([]byte, error) {
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    if len(cipherText) < aes.BlockSize {
        panic("ciphterText too short")
    }
    //密文前一部分为初始化向量 iv，无需保密
    iv := cipherText[:aes.BlockSize]
    data := cipherText[aes.BlockSize:]

    // CBC mode always works in whole blocks.
    if len(data) % aes.BlockSize != 0 {
        panic("cipherText is not a multiple of the block size")
    }

    mode := cipher.NewCBCDecrypter(block, iv)
    // CryptBlocks can work in-place if the two arguments are the same.
    dst := make([]byte, len(data))
    mode.CryptBlocks(dst, data)
    dst = PKCS5UnPadding(dst)
    // dst = ZeroUnPadding(dst)
    return dst, nil
}

func ZeroPadding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padtext := bytes.Repeat([]byte{0}, padding)
    return append(cipherText, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
    length := len(origData)
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}

func PKCS5Padding(cipherText []byte, blockSize int) []byte {
    padding := blockSize - len(cipherText)%blockSize
    padtext := bytes.Repeat([]byte{byte(padding)}, padding)
    return append(cipherText, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
    length := len(origData)
    // 去掉最后一个字节 unpadding 次
    unpadding := int(origData[length-1])
    return origData[:(length - unpadding)]
}
