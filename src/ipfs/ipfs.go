package ipfs

import (
	shell "github.com/ipfs/go-ipfs-api"
	"os"
)

func AddFileToIpfs(path string) (string, error) {
	sh := shell.NewShell("localhost:5001")

	ipfsFile, err := os.Open(path)
	if err != nil {
		return "", err
	}

	cid, err := sh.Add(ipfsFile)
	if err != nil {
		return "", err
	}

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
