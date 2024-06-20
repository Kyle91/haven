// @Author Eric
// @Date 2024/6/16 14:40:00
// @Desc
package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Kyle91/haven/crypto"
)

func main() {
	keyStr := "NTMxODU3NjMyNTM3NjU0MTA4NDEwMDE2MDA2NjEyNzM="
	key, _ := base64.StdEncoding.DecodeString(keyStr)
	fmt.Println(string(key))
	c := []byte{125, 169, 113, 86, 176, 3, 136, 19, 158, 193, 70, 35, 193, 96, 118, 140}
	decrypt, err := crypto.Aes256Decrypt(key, c)
	fmt.Println(err)
	fmt.Println(decrypt)

}
