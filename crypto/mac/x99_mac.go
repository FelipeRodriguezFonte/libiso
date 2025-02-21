/**
**/

package mac

import (
	"github.com/rkbalgi/libiso/crypto"
)

// GenerateMacX99 generates a X9.9 MAC using a single length key  data will be zero padded if required
func GenerateMacX99(inMacData []byte, keyData []byte) ([]byte, error) {

	macData := make([]byte, len(inMacData))
	copy(macData, inMacData)

	//add 0 padding
	if len(macData) < 8 || len(macData)%8 != 0 {
		pads := make([]byte, 8-(len(macData)%8))
		println("pads ", len(pads))
		macData = append(macData, pads...)
	}

	var err error
	result, err := crypto.EncryptDesCbc(macData, keyData)
	if err != nil {
		return nil, err
	}
	return result[len(result)-8:], nil

}
