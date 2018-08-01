package cryptor

import (
    "testing"
    "fmt"
    "encoding/hex"
    "io"
    "crypto/rand"
)

func TestAES(t *testing.T) {
    if true {
        return
    }
    //256位密钥
    key := []byte("6368616e676520746869732070617373")

    //测试 CBC 模式
    plaintext := []byte("exampleplaintextlslslsllslss")
    cipherText, err := AesCBCEncrypt(plaintext, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("plain: %v, cipher: %v\n", string(plaintext), hex.EncodeToString(cipherText))
    text, err := AesCBCDecrypt(cipherText, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("cipher: %v, plain: %v\n", hex.EncodeToString(cipherText), string(text))
    if string(text) != string(plaintext) {
        t.Fatalf("text mismatch: %v %v", text, plaintext)
    }

    //测试 CFB 模式
    cipherText, err = AesCFBEncrypt(plaintext, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("plain: %v, cipher: %v\n", string(plaintext), hex.EncodeToString(cipherText))
    text, err = AesCFBDecrypt(cipherText, key)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Printf("cipher: %v, plain: %v\n", hex.EncodeToString(cipherText), string(text))
    if string(text) != string(plaintext) {
        t.Fatalf("text mismatch: %v %v", text, plaintext)
    }
}

func BenchmarkAesCFB(b *testing.B) {
    key := []byte("6368616e67652074")

    plaintext := make([]byte, 10)
    if _, err := io.ReadFull(rand.Reader, plaintext); err != nil {
        panic(err)
    }
    var err error
    for i := 0; i < b.N; i++ {
        _, err = AesCFBEncrypt(plaintext, key)
        if err != nil {
            b.Fatal(err)
        }
    }

    //mac 4核
    //BenchmarkAesCFB-4   	 1000000	      1896 ns/op     //short message
    //BenchmarkAesCFB-4   	  200000	      7029 ns/op     //2k message
    //BenchmarkAesCFB-4   	    5000	    236397 ns/op     //200k message

}

func BenchmarkAesCBC(b *testing.B) {
    key := []byte("6368616e67652074")

    plaintext := make([]byte, 10)
    if _, err := io.ReadFull(rand.Reader, plaintext); err != nil {
        panic(err)
    }
    var err error
    for i := 0; i < b.N; i++ {
        _, err = AesCBCEncrypt(plaintext, key)
        if err != nil {
            b.Fatal(err)
        }
    }

    //mac 4核
    //BenchmarkAesCBC-4   	 1000000	      2057 ns/op       //short message
    //BenchmarkAesCBC-4   	  200000	      7411 ns/op       //2k message
    //BenchmarkAesCBC-4   	    5000	    227901 ns/op       //200k message
}
