package blockchain

import (
	"encoding/base64"
	"fmt"
	"github.com/volvinbur1/docs-chain/internal/blockchainresponse"
	"github.com/volvinbur1/docs-chain/src/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ReadImage(imagePath string) (string, error) {
	bytes, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	//mimeType := http.DetectContentType(bytes)
	//
	//switch mimeType {
	//case "image/jpeg":
	//	base64Encoding += "data:image/jpeg;base64,"
	//case "image/png":
	//	base64Encoding += "data:image/png;base64,"
	//}

	base64Encoding := toBase64(bytes)

	return base64Encoding, nil
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func (blockChain *BlockChain) GetImageUrl(base64Enc string) (string, error) {

	requestUrl := fmt.Sprintf("%s?%s=600&%s=%s",
		common.ImageBaseUrl, common.Expiration, common.Key, blockChain.ImgUrlApiKey)

	params := url.Values{}
	params.Add(common.Image, base64Enc)

	res, err := http.PostForm(requestUrl, params)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	return blockchainresponse.ParseImageResponse(res)
}
