package cryptor

import (
    "testing"
    "fmt"
    "encoding/hex"
    "io"
    "crypto/rand"
)

func TestXXTea(t *testing.T) {
    if false {
        return
    }
    //128位密钥
    key := []byte("6368616e67652074")

    //测试 CBC 模式
    plaintext := []byte("exampleplaintextlslslsllslss")
    cipherText, err := XXTeaEncrypt(plaintext, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("plain: %v, cipher: %v\n", string(plaintext), hex.EncodeToString(cipherText))

    text, err := XXTeaDecrypt(cipherText, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("cipher: %v, plain: %v\n", hex.EncodeToString(cipherText), string(text))
    if string(text) != string(plaintext) {
        t.Fatalf("text mismatch: %v %v", text, plaintext)
    }
}

func BenchmarkXXTea(b *testing.B) {
    key := []byte("6368616e67652074")

    plaintext := make([]byte, 10)
    if _, err := io.ReadFull(rand.Reader, plaintext); err != nil {
        panic(err)
    }
    var err error
    for i := 0; i < b.N; i++ {
        _, err = XXTeaEncrypt(plaintext, key)
        if err != nil {
            b.Fatal(err)
        }
    }
}