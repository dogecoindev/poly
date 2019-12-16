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

package relayer_manager

import (
	"fmt"

	"github.com/ontio/multi-chain/common"
	"github.com/ontio/multi-chain/core/genesis"
	"github.com/ontio/multi-chain/core/types"
	"github.com/ontio/multi-chain/native"
	"github.com/ontio/multi-chain/native/service/utils"
)

const (
	//function name
	REGISTER_RELAYER = "registerRelayer"
	REMOVE_RELAYER   = "RemoveRelayer"

	//key prefix
	RELAYER = "relayer"
)

//Register methods of node_manager contract
func RegisterRelayerManagerContract(native *native.NativeService) {
	native.Register(REGISTER_RELAYER, RegisterRelayer)
	native.Register(REMOVE_RELAYER, RemoveRelayer)
}

func RegisterRelayer(native *native.NativeService) ([]byte, error) {
	params := new(RelayerParam)
	if err := params.Deserialization(common.NewZeroCopySource(native.GetInput())); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RegisterRelayer, contract params deserialize error: %v", err)
	}

	// get operator from database
	operatorAddress, err := types.AddressFromBookkeepers(genesis.GenesisBookkeepers)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	//check witness
	err = utils.ValidateOwner(native, operatorAddress)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RegisterRelayer, checkWitness error: %v", err)
	}

	//get relayer
	relayerRaw, err := GetRelayerRaw(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RegisterRelayer, get relayer error: %v", err)
	}
	if relayerRaw != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RegisterRelayer, relayer is already registered")
	}

	err = putRelayer(native, params)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RegisterRelayer, putRelayer error: %v", err)
	}

	return utils.BYTE_TRUE, nil
}

func RemoveRelayer(native *native.NativeService) ([]byte, error) {
	params := new(RelayerParam)
	if err := params.Deserialization(common.NewZeroCopySource(native.GetInput())); err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RemoveRelayer, contract params deserialize error: %v", err)
	}
	contract := utils.RelayerManagerContractAddress

	// get operator from database
	operatorAddress, err := types.AddressFromBookkeepers(genesis.GenesisBookkeepers)
	if err != nil {
		return utils.BYTE_FALSE, err
	}

	//check witness
	err = utils.ValidateOwner(native, operatorAddress)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RemoveRelayer, checkWitness error: %v", err)
	}

	//get relayer
	relayerRaw, err := GetRelayerRaw(native, params.Address)
	if err != nil {
		return utils.BYTE_FALSE, fmt.Errorf("RemoveRelayer, get relayer error: %v", err)
	}
	if relayerRaw == nil {
		return utils.BYTE_FALSE, fmt.Errorf("RemoveRelayer, relayer is not registered")
	}

	native.GetCacheDB().Delete(utils.ConcatKey(contract, []byte(RELAYER), params.Address))

	return utils.BYTE_TRUE, nil
}