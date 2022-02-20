package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func readBytes(r io.Reader, n int) []byte {
	buf := make([]byte, n)
	_, err := r.Read(buf)
	if err != nil{
		return nil
	}
	return buf
}

func readBytesAsInt(r io.Reader, n int) int {
	if n == 4{
	return int(binary.BigEndian.Uint32(readBytes(r, n)))
	} else {
		return int(readBytes(r,n)[0])
	}
}

func main(){
	fileNama := os.Args[1]
	buf, err := ioutil.ReadFile(fileNama)
	if err != nil {
		return 
	}
	r := bytes.NewReader(buf)

	if !bytes.Equal(readBytes(r, 8), []byte("\x89PNG\r\n\x1a\n")) {
		fmt.Println("This file is not PNG")
		return
	}

	loop := true
	var colorType int

	for loop{
		Length := readBytesAsInt(r,4)
		Type := string(readBytes(r, 4))
		data := readBytes(r, Length)
		_ = readBytes(r, 4)
		fmt.Println("Chunk:",Type)

		// TO DO: 必須チャンク以外の追加
		switch Type {
		case "IHDR":
			ihdrNR := bytes.NewReader(data)
			width := readBytesAsInt(ihdrNR, 4)
			height := readBytesAsInt(ihdrNR, 4)
			depth := readBytesAsInt(ihdrNR, 1)
			colorType := readBytesAsInt(ihdrNR, 1)
			compression := readBytesAsInt(ihdrNR, 1)
			filter := readBytesAsInt(ihdrNR, 1)
			interlace := readBytesAsInt(ihdrNR, 1)
			fmt.Println("Width:",width,"Height:",height,"depth:",depth,"ColorType:",colorType,"Compression:",compression,"FilterType:",filter,"Interlace:",interlace)
		
		case "PLTE":
			plteNR := bytes.NewReader(data)
			paletteData := readBytes(plteNR, Length)
			fmt.Println("paletteData:", paletteData)

		case "tRNS":
			trnsNR := bytes.NewReader(data)
			if colorType == 3 {
				var PaletteAlpha []byte
				for i := 0; i < Length; i++ {
					PaletteAlpha[i] = readBytes(trnsNR, 1)[0]
					fmt.Println("PaletteNo.",i," Alpha:",PaletteAlpha[i])
				}
			} else if colorType == 0 {
				var GlayAlpha []byte
				for i := 0; i < Length/2; i++ {
					GlayAlpha[i] = readBytes(trnsNR, 2)[0]
					fmt.Println("GlayLevel.",i," Alpha:",GlayAlpha[i])
				}
			} else if colorType == 2 {
				var TransAlphaR []byte
				var TransAlphaG []byte
				var TransAlphaB []byte
				for i := 0; i < Length/6; i++ {
					TransAlphaR[i] = readBytes(trnsNR, 2)[0]
					TransAlphaG[i] = readBytes(trnsNR, 2)[0]
					TransAlphaB[i] = readBytes(trnsNR, 2)[0]
					fmt.Println("No.",i," Alpha R:",TransAlphaR[i],", G:",TransAlphaG[i],", B:",TransAlphaB[i])
				}
			}

		case "gAMA":
			gamaNR := bytes.NewReader(data)
			gamma := readBytesAsInt(gamaNR, Length)
			fmt.Println("gammaValue:",gamma)

		case "cHRM":
			chrmNR := bytes.NewReader(data)
			whitePointX := readBytesAsInt(chrmNR, 4)
			whitePointY := readBytesAsInt(chrmNR, 4)
			redX := readBytesAsInt(chrmNR, 4)
			redY := readBytesAsInt(chrmNR, 4)
			greenX := readBytesAsInt(chrmNR, 4)
			greenY := readBytesAsInt(chrmNR, 4)
			blueX := readBytesAsInt(chrmNR, 4)	
			blueY := readBytesAsInt(chrmNR, 4)
			fmt.Println("White Point X:", whitePointX, "White Point Y:", whitePointY, "Red X:", redX, "Red Y:", redY, "Green X:", greenX, "Green Y:", greenY, "Blue X:", blueX, "Blue Y:", blueY)

		case "tEXt":
			textNR := bytes.NewReader(data)
			keyWords := readBytes(textNR, Length)
			fmt.Println("KeyWords:",string(keyWords))

		// TO DO: Data部の展開
		case "IDAT":
			idatNR := bytes.NewReader(data)
			imgData := readBytes(idatNR, Length)
			fmt.Println("imageData:",imgData)
			
		case "IEND":
			loop = false
		}
	}
	fmt.Println("Complete")
}