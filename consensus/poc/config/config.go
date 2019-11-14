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

package pocconfig

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"OntologyWithPOC/common"
	"OntologyWithPOC/common/serialization"
)

var (
	Version uint32 = 1
)

type PeerConfig struct {
	Index uint32 `json:"index"`
	ID    string `json:"id"`
}

type ChainConfig struct {
	Version              uint32        `json:"version"` // software version
	View                 uint32        `json:"view"`    // config-updated version
	N                    uint32        `json:"n"`       // network size
	C                    uint32        `json:"c"`       // consensus quorum
	BlockMsgDelay        time.Duration `json:"block_msg_delay"`
	HashMsgDelay         time.Duration `json:"hash_msg_delay"`
	PeerHandshakeTimeout time.Duration `json:"peer_handshake_timeout"`
	Peers                []*PeerConfig `json:"peers"`
	PosTable             []uint32      `json:"pos_table"`
	MaxBlockChangeView   uint32        `json:"MaxBlockChangeView"`
}

//
// poc consensus payload, stored on each block header
//
type PocBlockInfo struct {
	Proposer           uint32       `json:"leader"`
	LastConfigBlockNum uint32       `json:"last_config_block_num"`
	NewChainConfig     *ChainConfig `json:"new_chain_config"`
}

const (
	MAX_PROPOSER_COUNT  = 32
	MAX_ENDORSER_COUNT  = 240
	MAX_COMMITTER_COUNT = 240
)

func VerifyChainConfig(cfg *ChainConfig) error {

	// TODO

	return nil
}

//Serialize the ChainConfig
func (cc *ChainConfig) Serialize(w io.Writer) error {
	data, err := json.Marshal(cc)
	if err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func (pc *PeerConfig) Serialize(w io.Writer) error {
	if err := serialization.WriteUint32(w, pc.Index); err != nil {
		return fmt.Errorf("ChainConfig peer index length serialization failed %s", err)
	}
	if err := serialization.WriteString(w, pc.ID); err != nil {
		return fmt.Errorf("ChainConfig peer ID length serialization failed %s", err)
	}
	return nil
}

func (pc *PeerConfig) Deserialize(r io.Reader) error {
	index, err := serialization.ReadUint32(r)
	if err != nil {
		return fmt.Errorf("serialization PeerConfig index err:%s", err)
	}
	pc.Index = index

	nodeid, err := serialization.ReadString(r)
	if err != nil {
		return fmt.Errorf("serialization PeerConfig nodeid err:%s", err)
	}
	pc.ID = nodeid
	return nil
}

func (cc *ChainConfig) Hash() common.Uint256 {
	buf := new(bytes.Buffer)
	cc.Serialize(buf)
	hash := sha256.Sum256(buf.Bytes())
	return hash
}