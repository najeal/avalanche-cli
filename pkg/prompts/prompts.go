// Copyright (C) 2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package prompts

import (
	"errors"
	"math/big"
	"os"
	"strconv"

	"github.com/ava-labs/avalanche-cli/pkg/models"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/formatting/address"
	"github.com/ethereum/go-ethereum/common"
	"github.com/manifoldco/promptui"
)

const (
	Yes = "Yes"
	No  = "No"
)

func validatePositiveBigInt(input string) error {
	n := new(big.Int)
	n, ok := n.SetString(input, 10)
	if !ok {
		return errors.New("invalid number")
	}
	if n.Cmp(big.NewInt(0)) == -1 {
		return errors.New("invalid number")
	}
	return nil
}

func validateAddress(input string) error {
	if !common.IsHexAddress(input) {
		return errors.New("invalid address")
	}
	return nil
}

func validateExistingFilepath(input string) error {
	if fileInfo, err := os.Stat(input); err == nil && !fileInfo.IsDir() {
		return nil
	}
	return errors.New("file doesn't exist")
}

func validateBiggerThanZero(input string) error {
	val, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		return err
	}
	if val == 0 {
		return errors.New("the value must be bigger than zero")
	}
	return nil
}

func CaptureUint64(promptStr string) (uint64, error) {
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: validateBiggerThanZero,
	}

	amountStr, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	val, err := strconv.ParseUint(amountStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

func CapturePositiveBigInt(promptStr string) (*big.Int, error) {
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: validatePositiveBigInt,
	}

	amountStr, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	amountInt := new(big.Int)
	amountInt, ok := amountInt.SetString(amountStr, 10)
	if !ok {
		return nil, errors.New("SetString: error")
	}
	return amountInt, nil
}

func validatePChainAddress(input string) (string, error) {
	chainID, hrp, _, err := address.Parse(input)
	if err != nil {
		return "", err
	}

	if chainID != "P" {
		return "", errors.New("this is not a PChain address")
	}
	return hrp, nil
}

func validatePChainFujiAddress(input string) error {
	hrp, err := validatePChainAddress(input)
	if err != nil {
		return err
	}
	if hrp != constants.FujiHRP {
		return errors.New("this is not a fuji address")
	}
	return nil
}

func validatePChainMainAddress(input string) error {
	hrp, err := validatePChainAddress(input)
	if err != nil {
		return err
	}
	if hrp != constants.MainnetHRP {
		return errors.New("this is not a mainnet address")
	}
	return nil
}

func validatePChainLocalAddress(input string) error {
	hrp, err := validatePChainAddress(input)
	if err != nil {
		return err
	}
	// ANR uses the `custom` HRP for local networks,
	// but the `local` HRP also exists...
	if hrp != constants.LocalHRP && hrp != constants.FallbackHRP {
		return errors.New("this is not a local nor custom address")
	}
	return nil
}

func getPChainValidationFunc(network models.Network) func(string) error {
	switch network {
	case models.Fuji:
		return validatePChainFujiAddress
	case models.Mainnet:
		return validatePChainMainAddress
	case models.Local:
		return validatePChainLocalAddress
	default:
		return func(string) error {
			return errors.New("unsupported network")
		}
	}
}

func CapturePChainAddress(promptStr string, network models.Network) (string, error) {
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: getPChainValidationFunc(network),
	}

	return prompt.Run()
}

func CaptureAddress(promptStr string) (common.Address, error) {
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: validateAddress,
	}

	addressStr, err := prompt.Run()
	if err != nil {
		return common.Address{}, err
	}

	addressHex := common.HexToAddress(addressStr)
	return addressHex, nil
}

func CaptureExistingFilepath(promptStr string) (string, error) {
	prompt := promptui.Prompt{
		Label:    promptStr,
		Validate: validateExistingFilepath,
	}

	pathStr, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return pathStr, nil
}

func yesNoBase(promptStr string, orderedOptions []string) (bool, error) {
	prompt := promptui.Select{
		Label: promptStr,
		Items: orderedOptions,
	}

	_, decision, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return decision == Yes, nil
}

func CaptureYesNo(promptStr string) (bool, error) {
	return yesNoBase(promptStr, []string{Yes, No})
}

func CaptureNoYes(promptStr string) (bool, error) {
	return yesNoBase(promptStr, []string{No, Yes})
}

func CaptureList(promptStr string, options []string) (string, error) {
	prompt := promptui.Select{
		Label: promptStr,
		Items: options,
	}

	_, listDecision, err := prompt.Run()
	if err != nil {
		return "", err
	}
	return listDecision, nil
}

func CaptureString(promptStr string) (string, error) {
	prompt := promptui.Prompt{
		Label: promptStr,
		Validate: func(input string) error {
			if input == "" {
				return errors.New("string cannot be empty")
			}
			return nil
		},
	}

	str, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return str, nil
}

func CaptureIndex(promptStr string, options []common.Address) (int, error) {
	prompt := promptui.Select{
		Label: promptStr,
		Items: options,
	}

	listIndex, _, err := prompt.Run()
	if err != nil {
		return 0, err
	}
	return listIndex, nil
}
