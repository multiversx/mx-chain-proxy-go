package process

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-go/hashing/keccak"
	"github.com/ElrondNetwork/elrond-proxy-go/process/helpers"
)

var usernameHasher = &keccak.Keccak{}

const usernameHashLength = 32

const (
	lowUsernameLengthBoundary  = 3
	highUsernameLengthBoundary = 25
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
func (dp *DnsProcessor) GetDnsAddressForUsername(providedUsername string) (string, error) {
	username, err := computeUsername(providedUsername)
	if err != nil {
		return "", err
	}

	usernameBytes := usernameHasher.Compute(username)
	if len(usernameBytes) != usernameHashLength {
		return "", ErrHashedUsernameBelowLimit
	}

	address := dp.sortedEncodedAddresses[usernameBytes[usernameHashLength-1]]

	return address, nil
}

func computeUsername(providedUsername string) (string, error) {
	username := providedUsername

	if strings.Contains(providedUsername, ".") {
		splitStr := strings.Split(providedUsername, ".")
		if len(splitStr) != 2 {
			return "", ErrInvalidUsername
		}
		if fmt.Sprintf(".%s", splitStr[1]) != usernameSuffix {
			return "", ErrInvalidUsername
		}

		username = splitStr[0]
	}

	usernameLength := len(username)
	if !(usernameLength > lowUsernameLengthBoundary && usernameLength < highUsernameLengthBoundary) {
		return "", ErrInvalidUsernameLength
	}

	if !isUsernameAlphanumeric(username) {
		return "", fmt.Errorf("%w: only alphanumeric characters are allowed", ErrInvalidUsername)
	}

	username += usernameSuffix

	return username, nil
}

func isUsernameAlphanumeric(username string) bool {
	// use a basic for loop and check all the characters. 40x faster in benchmarks vs a regular expression
	for _, c := range username {
		if (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
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
