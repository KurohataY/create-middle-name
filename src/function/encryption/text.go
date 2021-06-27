package encryption

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"strings"
)

// RFC 1950 で定義される圧縮方式。
func Compress(text string) string {

	buf := new(bytes.Buffer)
	zw := zlib.NewWriter(buf)
	defer zw.Close()

	if _, err := io.Copy(zw, bytes.NewBufferString(text)); err != nil {
		return buf.String()
	}

	//display compressed data
	b := buf.Bytes()

	return string(b)
}

func Extract(text string) string {
	zr := strings.NewReader(text)
	r, err := zlib.NewReader(zr)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return "No result"
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.String()
}

// func test(text string) string {
// 	content := text //raw data
// 	//compress raw data
// 	zr, err := compress(bytes.NewBufferString(content))
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return "No result"
// 	}
// 	//display compressed data
// 	b := zr.Bytes()

// 	fmt.Printf(string(b))
// 	//extract from compressed data
// 	r, err := extract(zr)
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, err)
// 		return "No result"
// 	}
// 	buf := new(bytes.Buffer)
// 	buf.ReadFrom(r)
// 	return buf.String()
// }
