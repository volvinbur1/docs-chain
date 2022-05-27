package blockchainresponse

import (
	"errors"
	"github.com/volvinbur1/docs-chain/internal/bodyresponse"
	"net/http"
)

func ParseImageResponse(response *http.Response) (string, error) {
	bodyI, err := bodyresponse.ParseResponseBody(response)
	if err != nil {
		return "", err
	}

	image, isOkay := bodyI["data"].(map[string]interface{})
	if isOkay != true {
		return "", errors.New("error when parsing `data` to map[string]interface{}")
	}

	imageUrl, isOkay := image["url"].(string)
	if isOkay != true {
		return "", errors.New("error when parsing `url` to string")
	}

	return imageUrl, nil
}
