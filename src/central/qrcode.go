package central

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skip2/go-qrcode"
	"github.com/volvinbur1/docs-chain/src/common"
	"image/color"
)

func CreateNFTQRCode(metadata common.PaperNftMetadata) ([]byte, error) {
	nftJson, err := json.Marshal(metadata)
	if err != nil {
		return nil, fmt.Errorf("paper nft metadata marshal failed. Error %s", err)
	}

	qrCode, err := qrcode.New(string(nftJson), qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("q")
	}
	qrCode.ForegroundColor = color.RGBA{R: 0xff, G: 0xd7, B: 0x00, A: 0xff}
	qrCode.BackgroundColor = color.RGBA{R: 0x00, G: 0x57, B: 0xb7, A: 0xff}

	var imageBuffer bytes.Buffer
	err = qrCode.Write(256, &imageBuffer)
	if err != nil {
		return nil, fmt.Errorf("qrCode image data writing. Error %s", err)
	}
	return imageBuffer.Bytes(), nil
}
