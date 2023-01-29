package main

/// collider.go
/// Solidity 4byte function selector collision finder
/// author: lcfr.eth
///
/// usage: go run main.go -target 0x099aba56 -args "(address)" -prefix "lcfr_" -randLength 8
/// go mod init main
/// go get github.com/ethereum/go-ethereum

import (
	"time"
	"os"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
  "math/rand"
  "unsafe"
  "flag"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                 
    letterIdxMask = 1<<letterIdxBits - 1
    letterIdxMax  = 63 / letterIdxBits   
)
// found on stackoverflow https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func RandStringBytesMaskImprSrcUnsafe(n int) string {
    var src = rand.NewSource(time.Now().UnixNano())
    b := make([]byte, n)
    for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
        if remain == 0 {
            cache, remain = src.Int63(), letterIdxMax
        }
        if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
            b[i] = letterBytes[idx]
            i--
        }
        cache >>= letterIdxBits
        remain--
    }
    return *(*string)(unsafe.Pointer(&b))
}

func randomReverseCapitalization(s string) string {
	rand.Seed(time.Now().UTC().UnixNano())
	runes := []rune(s)
	randomIndex := rand.Intn(len(runes))

  if runes[randomIndex] >= 'A' && runes[randomIndex] <= 'Z' {
		runes[randomIndex] = runes[randomIndex] + ('a' - 'A')
	} else if runes[randomIndex] >= 'a' && runes[randomIndex] <= 'z' {
		runes[randomIndex] = runes[randomIndex] - ('a' - 'A')
	}
	return string(runes)
}

func getKeccak4(signature string, target string) {
  data := []byte(signature)
  hash := crypto.Keccak256Hash(data)
  fourByte := hash.Hex()[0:10]
  //fmt.Printf("got: %v\n", fourByte)
  if( fourByte == target ) {
    fmt.Printf("Collision found! Target %s = %s, signature: %s\n", target, fourByte, signature)
    elapsed := time.Since(start)
    fmt.Printf("Found in: %s\n", elapsed)
    os.Exit(2)
  }
}

var start = time.Now()

func main() {

  targetFlag := flag.String("target", "0x099aba56", "target to match") // isTalentToken(address)
  argsFlag  := flag.String("args", "(address)", "args to match")
  prefixFlag := flag.String("prefix", "", "prefix to match")
  padFlag := flag.Int("pad", 8, "bytes to pad")
  flag.Parse()
  fmt.Printf("starting collider ... target: %s, args: %s\n", *targetFlag, *argsFlag)

	for { 
    var hack string
    if ( *prefixFlag != "" ) {
      //hack = randomReverseCapitalization(*prefixFlag) + RandStringBytesMaskImprSrcUnsafe(*randLengthFlag)
      hack = *prefixFlag + RandStringBytesMaskImprSrcUnsafe(*padFlag)

    } else {
      hack = RandStringBytesMaskImprSrcUnsafe(*padFlag)
    }

    final := hack + *argsFlag
    //fmt.Printf("using: %s\n", final)
	  //go getKeccak4(final, target)
	  //getKeccak4(final, *targetFlag)

    // use a go routine to get the keccak hash but manage the max number of go routines
    // so we don't run out of memory
    go func() {
      getKeccak4(final, *targetFlag)
    } ()
	}	
  
}
