package ipfs

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"os"
)

func AddFileToIpfs(path string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	ipfsFile, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file for IPFS. Error %s", err)
	}

	cid, err := sh.Add(ipfsFile)
	if err != nil {
		return "", fmt.Errorf("addinf file to IPFS. Error %s", err)
	}

	defer ipfsFile.Close()

	return cid, nil

}

func GetFileFromIpfs(cid string, outDir string) error {
	sh := shell.NewShell("localhost:5001")

	err := sh.Get(cid, outDir)
	if err != nil {
		return err
	}
	return nil
}
