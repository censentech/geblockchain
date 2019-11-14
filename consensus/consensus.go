/*
 * Copyright (C) 2018 The ontology Authors
 * This file is part of The ontology library.
 *
 * The ontology is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The ontology is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
 */

package consensus

import (
	"OntologyWithPOC/account"
	"OntologyWithPOC/common/config"
	"OntologyWithPOC/common/log"
	"OntologyWithPOC/consensus/dbft"
	"OntologyWithPOC/consensus/poc"
	"OntologyWithPOC/consensus/poc/config"
	"OntologyWithPOC/consensus/solo"
	"OntologyWithPOC/consensus/vbft"
	"github.com/ontio/ontology-eventbus/actor"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type ConsensusService interface {
	Start() error
	Halt() error
	GetPID() *actor.PID
}

const (
	CONSENSUS_DBFT = "dbft"
	CONSENSUS_SOLO = "solo"
	CONSENSUS_VBFT = "vbft"
	CONSENSUS_POC  = "poc"
)

var quitWg sync.WaitGroup

//调用os.MkdirAll递归创建文件夹
func createFile(filePath string) error {
	if !isExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func GetAllFileSize(pathname string) uint32 {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		panic(err)
	}
	var sumSize uint32
	for _, fi := range rd {
		if fi.IsDir() {
			log.Info("[%s]\n", pathname+"\\"+fi.Name())
		} else {
			sumSize += uint32(fi.Size())
		}
	}
	return sumSize
}

func GenShabalData(account *account.Account) {
	quitWg.Add(1)
	quitWg.Done()

	err := createFile(config.DefConfig.Genesis.POC.NonceDir)
	if err != nil {
		panic(err)
	}
	filespace := GetAllFileSize(config.DefConfig.Genesis.POC.NonceDir)
	pocspace := config.DefConfig.Genesis.POC.PocSpace
	dfspace := pocspace * 1024 * 1024
	if filespace < dfspace && (dfspace-filespace)/262144 != 0 {
		for i := uint32(0); i < (dfspace-filespace)/262144; i++ {
			nonceNr := rand.New(rand.NewSource(time.Now().UnixNano())).Uint64()
			pubkey := pocconfig.PubkeyID(account.PubKey())
			poc.Callshabal("genNonce256", []byte(strconv.FormatUint(nonceNr, 10)), []byte(pubkey),
				[]byte(strconv.Itoa(0)), []byte(""), []byte(config.DefConfig.Genesis.POC.NonceDir))
		}
	} else {
		log.Info("There is enough nonce file, the space is more than the default config!!!")
	}
}

func NewConsensusService(consensusType string, account *account.Account, txpool *actor.PID, ledger *actor.PID, p2p *actor.PID) (ConsensusService, error) {
	if consensusType == "" {
		consensusType = CONSENSUS_DBFT
	}
	var consensus ConsensusService
	var err error
	switch consensusType {
	case CONSENSUS_DBFT:
		consensus, err = dbft.NewDbftService(account, txpool, p2p)
	case CONSENSUS_SOLO:
		consensus, err = solo.NewSoloService(account, txpool)
	case CONSENSUS_VBFT:
		consensus, err = vbft.NewVbftServer(account, txpool, p2p)
	case CONSENSUS_POC:
		go GenShabalData(account)
		consensus, err = poc.NewPocServer(account, txpool, p2p)
	}
	log.Infof("ConsensusType:%s", consensusType)
	return consensus, err
}
