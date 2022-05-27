package blockchainresponse

import (
	"errors"
	"github.com/volvinbur1/docs-chain/internal/bodyresponse"
	"github.com/volvinbur1/docs-chain/src/common"
	"net/http"
)

func ParseMintInfo(response *http.Response) (common.NftMintResponse, error) {
	nftInfo, err := bodyresponse.ParseResponseBody(response)
	if err != nil {
		return common.NftMintResponse{}, err
	}
	var isOkay bool
	nftResponse := common.NftMintResponse{}

	nftResponse.NftMint, isOkay = nftInfo["mint"].(string)
	if isOkay != true {
		return common.NftMintResponse{}, errors.New("error when casting `mint` to string")
	}

	nftResponse.TransactionSignature, isOkay = nftInfo["transaction_signature"].(string)
	if isOkay != true {
		return common.NftMintResponse{}, errors.New("error when casting `transaction_signature` to string")
	}

	nftResponse.MintRecoveryPhrase, isOkay = nftInfo["mint_secret_recovery_phrase"].(string)
	if isOkay != true {
		return common.NftMintResponse{}, errors.New("error when casting `mint_secret_recovery_phrase` to string")
	}

	return nftResponse, nil
}
