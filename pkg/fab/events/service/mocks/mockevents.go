/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	"github.com/golang/protobuf/proto"
	cb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/common"
	pb "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/protos/peer"
)

// NewBlock returns a new mock block initialized with the given channel
func NewBlock(channelID string, transactions ...*TxInfo) *cb.Block {
	var data [][]byte
	txValidationFlags := make([]uint8, len(transactions))
	for i, txInfo := range transactions {
		envBytes, err := proto.Marshal(newEnvelope(channelID, txInfo))
		if err != nil {
			panic(err)
		}
		data = append(data, envBytes)
		txValidationFlags[i] = uint8(txInfo.TxValidationCode)
	}

	blockMetaData := make([][]byte, 4)
	blockMetaData[cb.BlockMetadataIndex_TRANSACTIONS_FILTER] = txValidationFlags

	return &cb.Block{
		Header:   &cb.BlockHeader{},
		Metadata: &cb.BlockMetadata{Metadata: blockMetaData},
		Data:     &cb.BlockData{Data: data},
	}
}

// TxInfo contains the data necessary to
// construct a mock transaction
type TxInfo struct {
	TxID             string
	TxValidationCode pb.TxValidationCode
	HeaderType       cb.HeaderType
	ChaincodeID      string
	EventName        string
}

// NewTransaction creates a new transaction
func NewTransaction(txID string, txValidationCode pb.TxValidationCode, headerType cb.HeaderType) *TxInfo {
	return &TxInfo{
		TxID:             txID,
		TxValidationCode: txValidationCode,
		HeaderType:       headerType,
	}
}

// NewTransactionWithCCEvent creates a new transaction with the given chaincode event
func NewTransactionWithCCEvent(txID string, txValidationCode pb.TxValidationCode, ccID string, eventName string) *TxInfo {
	return &TxInfo{
		TxID:             txID,
		TxValidationCode: txValidationCode,
		ChaincodeID:      ccID,
		EventName:        eventName,
		HeaderType:       cb.HeaderType_ENDORSER_TRANSACTION,
	}
}

// NewFilteredBlock returns a new mock filtered block initialized with the given channel
// and filtered transactions
func NewFilteredBlock(channelID string, filteredTx ...*pb.FilteredTransaction) *pb.FilteredBlock {
	return &pb.FilteredBlock{
		ChannelId:            channelID,
		FilteredTransactions: filteredTx,
	}
}

// NewFilteredTx returns a new mock filtered transaction
func NewFilteredTx(txID string, txValidationCode pb.TxValidationCode) *pb.FilteredTransaction {
	return &pb.FilteredTransaction{
		Txid:             txID,
		TxValidationCode: txValidationCode,
	}
}

// NewFilteredTxWithCCEvent returns a new mock filtered transaction
// with the given chaincode event
func NewFilteredTxWithCCEvent(txID, ccID, event string) *pb.FilteredTransaction {
	return &pb.FilteredTransaction{
		Txid: txID,
		Data: &pb.FilteredTransaction_TransactionActions{
			TransactionActions: &pb.FilteredTransactionActions{
				ChaincodeActions: []*pb.FilteredChaincodeAction{
					&pb.FilteredChaincodeAction{
						ChaincodeEvent: &pb.ChaincodeEvent{
							ChaincodeId: ccID,
							EventName:   event,
							TxId:        txID,
						},
					},
				},
			},
		},
	}
}

func newEnvelope(channelID string, txInfo *TxInfo) *cb.Envelope {
	tx := &pb.Transaction{
		Actions: []*pb.TransactionAction{newTxAction(txInfo.TxID, txInfo.ChaincodeID, txInfo.EventName)},
	}
	txBytes, err := proto.Marshal(tx)
	if err != nil {
		panic(err)
	}

	channelHeader := &cb.ChannelHeader{
		ChannelId: channelID,
		TxId:      txInfo.TxID,
		Type:      int32(txInfo.HeaderType),
	}
	channelHeaderBytes, _ := proto.Marshal(channelHeader)

	payload := &cb.Payload{
		Header: &cb.Header{
			ChannelHeader: channelHeaderBytes,
		},
		Data: txBytes,
	}
	payloadBytes, _ := proto.Marshal(payload)

	return &cb.Envelope{
		Payload: payloadBytes,
	}
}

func newTxAction(txID string, ccID string, eventName string) *pb.TransactionAction {
	ccEvent := &pb.ChaincodeEvent{
		TxId:        txID,
		ChaincodeId: ccID,
		EventName:   eventName,
	}
	eventBytes, err := proto.Marshal(ccEvent)
	if err != nil {
		panic(err)
	}

	chaincodeAction := &pb.ChaincodeAction{
		ChaincodeId: &pb.ChaincodeID{
			Name: ccID,
		},
		Events: eventBytes,
	}
	extBytes, err := proto.Marshal(chaincodeAction)
	if err != nil {
		panic(err)
	}

	prp := &pb.ProposalResponsePayload{
		Extension: extBytes,
	}

	prpBytes, err := proto.Marshal(prp)
	if err != nil {
		panic(err)
	}

	cap := &pb.ChaincodeActionPayload{
		Action: &pb.ChaincodeEndorsedAction{
			ProposalResponsePayload: prpBytes,
		},
	}
	payloadBytes, err := proto.Marshal(cap)
	if err != nil {
		panic(err)
	}

	return &pb.TransactionAction{
		Payload: payloadBytes,
		Header:  nil,
	}
}
