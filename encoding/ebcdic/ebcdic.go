package ebcdic

import (
	"bytes"
	"fmt"
	"strconv"
	_ "strings"
)

//an ebcdic to unicode mapper
//based on http://www-01.ibm.com/support/docview.wss?uid=swg21145758

//A unicode mapper is also available at http://www.unicode.org/Public/MAPPINGS/VENDORS/MICSFT/EBCDIC/CP037.TXT

/*
ASCII TO EBCDIC TRANSLATE TABLE

*                 0 1 2 3 4 5 6 7 8 9 A B C D E F
TOEBTABL DC    X'000102030405060708090A0B0C0D0E0F'  00-0F
         DC    X'101112131415161718191A1B1C1D1E1F'  10-1F
         DC    X'405A7F7B5B6C507D4D5D5C4E6B604B61'  20-2F
         DC    X'F0F1F2F3F4F5F6F7F8F97A5E4C7E6E6F'  30-3F
         DC    X'7CC1C2C3C4C5C6C7C8C9D1D2D3D4D5D6'  40-4F
         DC    X'D7D8D9E2E3E4E5E6E7E8E9ADE0BDB06D'  50-5F
         DC    X'79818283848586878889919293949596'  60-6F
         DC    X'979899A2A3A4A5A6A7A8A9C04FD0A17F'  70-7F
         DC    X'808182838485868788898A8B8C8D8E8F'  80-8F
         DC    X'909192939495969798999A9B9C9D9E9F'  90-9F
         DC    X'A0A1A2A3A4A5A6A7A8A9AAABACADAEAF'  A0-AF
         DC    X'B0B1B2B3B4B5B6B7B8B9BABBBCBDBEBF'  B0-BF
         DC    X'C0C1C2C3C4C5C6C7C8C9CACBCCCDCECF'  C0-CF
         DC    X'D0D1D2D3D4D5D6D7D8D9DADBDCDDDEDF'  D0-DF
         DC    X'E0E1E2E3E4E5E6E7E8E9EAEBECEDEEEF'  E0-EF
         DC    X'F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF'  F0-FF

EBCDIC TO ASCII TRANSLATE TABLE

*                 0 1 2 3 4 5 6 7 8 9 A B C D E F
TOASTABL DC    X'000102030405060708090A0B0C0D0E0F'  00-0F
         DC    X'101112131415161718191A1B1C1D1E1F'  10-1F
         DC    X'202122232425262728292A2B2C2D2E2F'  20-2F
         DC    X'303132333435363738393A3B3C3D3E3F'  30-3F
         DC    X'202122232425262728292A2E3C282B6A'  40-4F
         DC    X'2651525354555657585921242A293B7E'  50-5F
         DC    X'2D2F62636465666768696A2C255F3E3F'  60-6F
         DC    X'707172737475767778603A2340273D22'  70-7F
         DC    X'806162636465666768698A8B8C8D8E8F'  80-8F
         DC    X'906A6B6C6D6E6F7071729A9B9C9D9E9F'  90-9F
         DC    X'A0A1737475767778797AAAABACADAEAF'  A0-AF
         DC    X'5EB1B2B3B4B5B6B7B8B9BABBBCBDBEBF'  B0-BF
         DC    X'7B414243444546474849CACBCCCDCECF'  C0-CF
         DC    X'7D4A4B4C4D4E4F505152DADBDCDDDEDF'  D0-DF
         DC    X'5CE1535455565758595AEAEBECEDEEEF'  E0-EF
         DC    X'30313233343536373839FAFBFCFDFEFF'  F0-FF

**/

var ebcdic_to_ascii = 
    "000102030405060708090A0B0C0D0E0F" +
	"101112131415161718191A1B1C1D1E1F" +
	"202122232425262728292A2B2C2D2E2F" +
	"303132333435363738393A3B3C3D3E3F" +
	"202122232425262728292A2E3C282B6A" +
	"2651525354555657585921242A293B7E" +
	"2D2F62636465666768696A2C255F3E3F" +
	"707172737475767778603A2340273D22" +
	"806162636465666768698A8B8C8D8E8F" +
	"906A6B6C6D6E6F7071729A9B9C9D9E9F" +
	"A0A1737475767778797AAAABACADAEAF" +
	"5EB1B2B3B4B5B6B7B8B9BABBBCBDBEBF" +
	"7B414243444546474849CACBCCCDCECF" +
	"7D4A4B4C4D4E4F505152DADBDCDDDEDF" +
	"5CE1535455565758595AEAEBECEDEEEF" +
	"30313233343536373839FAFBFCFDFEFF"

var ascii_to_ebcdic = "000102030405060708090A0B0C0D0E0F" +
	"101112131415161718191A1B1C1D1E1F" +
	"405A7F7B5B6C507D4D5D5C4E6B604B61" +
	"F0F1F2F3F4F5F6F7F8F97A5E4C7E6E6F" +
	"7CC1C2C3C4C5C6C7C8C9D1D2D3D4D5D6" +
	"D7D8D9E2E3E4E5E6E7E8E9ADE0BDB06D" +
	"79818283848586878889919293949596" +
	"979899A2A3A4A5A6A7A8A9C04FD0A17F" +
	"808182838485868788898A8B8C8D8E8F" +
	"909192939495969798999A9B9C9D9E9F" +
	"A0A1A2A3A4A5A6A7A8A9AAABACADAEAF" +
	"B0B1B2B3B4B5B6B7B8B9BABBBCBDBEBF" +
	"C0C1C2C3C4C5C6C7C8C9CACBCCCDCECF" +
	"D0D1D2D3D4D5D6D7D8D9DADBDCDDDEDF" +
	"E0E1E2E3E4E5E6E7E8E9EAEBECEDEEEF" +
	"F0F1F2F3F4F5F6F7F8F9FAFBFCFDFEFF"

//convert from ebcdic bytes to a ascii string
func EncodeToString(data []byte) string {

	buf := bytes.NewBufferString("")

	for _, b := range data {
		var x uint32 = uint32(b)
		fmt.Println("x",x);
		tmp := ebcdic_to_ascii[(x * 2) : (x*2)+2]
		i, _ := strconv.ParseUint(tmp, 16, 8)
		fmt.Println("i",i);
		buf.WriteString(string(i))

	}

	return buf.String()
}

//convert from a ascii encoded string to ebcdic bytes
func Decode(str string) []byte {

	data := make([]byte, len(str))

	for i := 0; i < len(str); i++ {
		b := uint32(str[i])
		b_val, _ := strconv.ParseUint(ascii_to_ebcdic[b*2:b*2+2], 16, 8)
		data[i] = byte(b_val)
	}

	return data
}
