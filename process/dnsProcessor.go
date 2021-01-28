package process

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/hashing/keccak"
	"github.com/ElrondNetwork/elrond-proxy-go/process/helpers"
)

var usernameHasher = &keccak.Keccak{}

const (
	lowUsernameLengthBoundary  = 3
	highUsernameLengthBoundary = 20
	usernameSuffix             = ".elrond"
)

// DnsProcessor handles dns operations
type DnsProcessor struct {
	sortedEncodedAddresses []string
	pubKeyConverter        core.PubkeyConverter
}

// NewDnsProcessor initializes all the inner components and returns the newly created instance of DnsProcessor
func NewDnsProcessor(pubKeyConverter core.PubkeyConverter) (*DnsProcessor, error) {
	if check.IfNil(pubKeyConverter) {
		return nil, ErrNilPubKeyConverter
	}

	dp := &DnsProcessor{
		pubKeyConverter: pubKeyConverter,
	}

	err := dp.computeDnsAddresses()
	if err != nil {
		return nil, err
	}

	return dp, nil
}

// GetDnsAddresses returns all the dns addresses
func (dp *DnsProcessor) GetDnsAddresses() ([]string, error) {
	return dp.sortedEncodedAddresses, nil
}

// GetDnsAddressForUsername returns the corresponding dns address for the provided username
func (dp *DnsProcessor) GetDnsAddressForUsername(username string) (string, error) {
	err := checkUsername(username)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(username, usernameSuffix) {
		username += usernameSuffix
		if len(username) > highUsernameLengthBoundary {
			return "", ErrInvalidUsernameLength
		}
	}

	usernameBytes := usernameHasher.Compute(username)
	if len(usernameBytes) < 32 {
		return "", ErrHashedUsernameBelowLimit
	}

	address := dp.sortedEncodedAddresses[usernameBytes[31]]

	return address, nil

}

func checkUsername(username string) error {
	usernameLength := len(username)
	if !(usernameLength > lowUsernameLengthBoundary && usernameLength < highUsernameLengthBoundary) {
		return ErrInvalidUsernameLength
	}

	for _, c := range username {
		if (c < 'a' || c > 'z') && c != '.' {
			return errors.New(fmt.Sprintf("invalid character: %c . only lowercase letters are allowed", c))
		}
	}

	return nil
}

func (dp *DnsProcessor) computeDnsAddresses() error {
	var initialDNSAddress = bytes.Repeat([]byte{1}, 32)
	addressLen := len(initialDNSAddress)

	addresses := make([][]byte, 0)
	for i := 0; i < 256; i++ {
		shardInBytes := []byte{0, uint8(i)}
		newDNSPk := string(initialDNSAddress[:(addressLen-core.ShardIdentiferLen)]) + string(shardInBytes)
		addr, err := helpers.CreateScAddress([]byte(newDNSPk), 0)
		if err != nil {
			return err
		}

		addresses = append(addresses, addr)
	}

	sort.SliceStable(addresses, func(i, j int) bool {
		iEl := addresses[i]
		jEl := addresses[j]
		return iEl[len(iEl)-1] < jEl[len(jEl)-1]
	})

	encodedAddresses := make([]string, 0)
	for _, addr := range addresses {
		encodedAddress := dp.pubKeyConverter.Encode(addr)
		encodedAddresses = append(encodedAddresses, encodedAddress)
	}

	dp.sortedEncodedAddresses = encodedAddresses

	return nil
}
