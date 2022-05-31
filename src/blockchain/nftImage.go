package blockchain

import (
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/ipfs"
)

func GetImageUrl(path string) (string, error) {

	cid, err := ipfs.AddFileToIpfs(path)
	if err != nil {
		return "", err
	}

	publicImgUrl := common.IpfsUrl + cid

	return publicImgUrl, nil
}
