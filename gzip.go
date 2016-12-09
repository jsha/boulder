package main

import (
	"bytes"
	"log"
	"os"
)
import "compress/gzip"
import "encoding/hex"
import "encoding/base64"

func main() {
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	input, err := base64.StdEncoding.DecodeString(`CidncnBjL3Rlc3RfcHJvdG8vaW50ZXJjZXB0b3JzX3Rlc3QucHJvdG8iFAoEVGltZRIMCgR0aW1lGAEgASgDMiIKB0NoaWxsZXISFwoFQ2hpbGwSBS5UaW1lGgUuVGltZSIAQgxaCnRlc3RfcHJvdG8=`)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(input)
	w.Close()
	os.Stdout.Write([]byte(hex.EncodeToString(buf.Bytes())))
	os.Stdout.Write([]byte("\n"))
}
