package initx

import (
	"github.com/flowshield/deca/internal/config"
	"github.com/flowshield/deca/pkg/errors"
	"github.com/flowshield/deca/pkg/storage"
	"github.com/flowshield/deca/pkg/storage/ipfs"
	web3_storage "github.com/flowshield/deca/pkg/storage/web3-storage"
)

// InitStorage 初始化存储引擎
func InitStorage() (storage.ExecCloser, func(), error) {
	cfg := config.C
	var db storage.ExecCloser
	var err error
	switch cfg.Storage.Type {
	case "ipfs":
		db, err = ipfs.New()
		if err != nil {
			return nil, nil, err
		}
	case "web3.storage":
		db, err = web3_storage.New()
		if err != nil {
			return nil, nil, err
		}
	default:
		return nil, nil, errors.New("unknown storage")
	}
	return db, func() {
		if db != nil {
			db.Close()
		}
	}, nil
}
