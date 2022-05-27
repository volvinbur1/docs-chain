package blockchain

import (
	"github.com/volvinbur1/docs-chain/internal/blockchainresponse"
	"github.com/volvinbur1/docs-chain/src/common"
	"net/http"
	"net/url"
	"strings"
)

func (blockChain *BlockChain) AddToSolana(nftStruct NftStruct) (common.NftMintResponse, error) {
	reqUrl := common.NftBaseUrl + common.MintEndpoint

	params := url.Values{}
	params.Add(common.Mnemonic, blockChain.Mnemonic)
	params.Add(common.DerivationPath, "")
	params.Add(common.NfrName, nftStruct.NftName)
	params.Add(common.NftSymbol, nftStruct.NftSymbol)
	params.Add(common.Description, nftStruct.NftDescription)
	params.Add(common.Network, common.DevNetwork)
	params.Add(common.NftUrl, nftStruct.NftImageUrl)
	params.Add(common.NftUploadMethod, common.Link)

	req, err := http.NewRequest(http.MethodPost, reqUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return common.NftResponse{}, err
	}

	req.Header.Add(common.ApiKey, blockChain.BlockApiKey)
	req.Header.Add(common.ApiSecret, blockChain.BlockApiSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return common.NftResponse{}, err
	}

	defer res.Body.Close()

	return blockchainresponse.ParseMintInfo(res)
}
