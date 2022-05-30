package analyzer

import (
	"github.com/golang/mock/gomock"
	"github.com/volvinbur1/docs-chain/src/common"
	"github.com/volvinbur1/docs-chain/src/storage/mock_storage"
	"testing"
)

func TestPaperPdfProcessor_PrepareFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newPaper := common.UploadedPaper{
		Id:       "test_id",
		FilePath: "test_data/test1.pdf",
	}
	db := mock_storage.NewMockDatabaseInterface(ctrl)
	processor := NewPaperPdfProcessor(newPaper, db, nil)
	err := processor.PrepareFile()
	if err != nil {
		t.Errorf("Error is expected not ot be nil")
	}
}
