package pow

import (
	"fmt"
	"testing"
)

func Test_POW(t *testing.T) {
	str, err := Generate("IMnbxy37CPIC_MeImjFH1YualVscQcX8c9pCioW50oGsjiHqkGuM60gMpw6FcIW6")
	fmt.Println(str, err)
}

/*
1:20:230228:9eGkLyBiQDkhJLZkw1IhPDzqNuGgbUBEuvzE5odFJCRD3fo9awpv8bBDSqI2iDS_::AMdiOfOT:7FFFFFFFFFFFFFFF
*/
