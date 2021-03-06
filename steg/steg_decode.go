package steg

import (
	"encoding/binary"
	"fmt"
	"github.com/DimitarPetrov/stegify/bits"
	"image"
	"io"
	"os"
)

//Decode performs steganography decoding of Reader with previously encoded data by the Encode function and writes to result Writer.
func Decode(carrier io.Reader, result io.Writer) error {

	RGBAImage, _, err := getImageAsRGBA(carrier)
	if err != nil {
		return fmt.Errorf("error parsing carrier image: %v", err)
	}

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	dataBytes := make([]byte, 0, 2048)
	resultBytes := make([]byte, 0, 2048)

	dataCount := extractDataCount(RGBAImage)

	var count int

	for x := 0; x < dx && dataCount > 0; x++ {
		for y := 0; y < dy && dataCount > 0; y++ {

			if count >= dataSizeHeaderReservedBytes {
				c := RGBAImage.RGBAAt(x, y)
				dataBytes = append(dataBytes, bits.GetLastTwoBits(c.R), bits.GetLastTwoBits(c.G), bits.GetLastTwoBits(c.B))
				dataCount -= 3
			} else {
				count += 4
			}
		}
	}

	if dataCount < 0 {
		dataBytes = dataBytes[:len(dataBytes)+dataCount] //remove bytes that are not part of data and mistakenly added
	}

	dataBytes = align(dataBytes) // len(dataBytes) must be aliquot of 4

	for i := 0; i < len(dataBytes); i += 4 {
		resultBytes = append(resultBytes, bits.ConstructByteOfQuartersAsSlice(dataBytes[i:i+4]))
	}

	result.Write(resultBytes)

	return nil
}

//DecodeByFileNames performs steganography decoding of data previously encoded by the Encode function.
//The data is decoded from file carrier and it is saved in separate new file
func DecodeByFileNames(carrierFileName string, newFileName string) error {
	carrier, err := os.Open(carrierFileName)
	defer carrier.Close()
	if err != nil {
		return fmt.Errorf("error opening carrier file: %v", err)
	}

	result, err := os.Create(newFileName)
	defer result.Close()
	if err != nil {
		return fmt.Errorf("error creating result file: %v", err)
	}

	err = Decode(carrier, result)
	if err != nil {
		os.Remove(newFileName)
	}

	return err

}

func align(dataBytes []byte) []byte {
	switch len(dataBytes) % 4 {
	case 1:
		dataBytes = append(dataBytes, byte(0), byte(0), byte(0))
	case 2:
		dataBytes = append(dataBytes, byte(0), byte(0))
	case 3:
		dataBytes = append(dataBytes, byte(0))
	}
	return dataBytes
}

func extractDataCount(RGBAImage *image.RGBA) int {
	dataCountBytes := make([]byte, 0, 16)

	dx := RGBAImage.Bounds().Dx()
	dy := RGBAImage.Bounds().Dy()

	count := 0

	for x := 0; x < dx && count < dataSizeHeaderReservedBytes; x++ {
		for y := 0; y < dy && count < dataSizeHeaderReservedBytes; y++ {

			c := RGBAImage.RGBAAt(x, y)
			dataCountBytes = append(dataCountBytes, bits.GetLastTwoBits(c.R), bits.GetLastTwoBits(c.G), bits.GetLastTwoBits(c.B))
			count += 4

		}
	}

	dataCountBytes = append(dataCountBytes, byte(0))

	var bs = []byte{bits.ConstructByteOfQuartersAsSlice(dataCountBytes[:4]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[4:8]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[8:12]),
		bits.ConstructByteOfQuartersAsSlice(dataCountBytes[12:])}

	return int(binary.LittleEndian.Uint32(bs))
}
