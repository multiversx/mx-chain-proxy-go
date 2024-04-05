package facade

import "github.com/pkg/errors"

// ErrNilActionsProcessor signals that a nil actions processor has been provided
var ErrNilActionsProcessor = errors.New("nil actions processor provided")

// ErrNilAccountProcessor signals that a nil account processor has been provided
var ErrNilAccountProcessor = errors.New("nil account processor provided")

// ErrNilTransactionProcessor signals that a nil transaction processor has been provided
var ErrNilTransactionProcessor = errors.New("nil transaction processor provided")

// ErrNilSCQueryService signals that a nil smart contracts query service has been provided
var ErrNilSCQueryService = errors.New("nil smart contracts query service provided")

// ErrNilNodeGroupProcessor signals that a nil node group processor has been provided
var ErrNilNodeGroupProcessor = errors.New("nil node group processor provided")

// ErrNilValidatorStatisticsProcessor signals that a nil validator statistics processor has been provided
var ErrNilValidatorStatisticsProcessor = errors.New("nil validator statistics processor provided")

// ErrNilFaucetProcessor signals that a nil faucet processor has been provided
var ErrNilFaucetProcessor = errors.New("nil faucet processor provided")

// ErrNilNodeStatusProcessor signals that a nil node status processor has been provided
var ErrNilNodeStatusProcessor = errors.New("nil node status processor provided")

// ErrNilBlockProcessor signals that a nil block processor has been provided
var ErrNilBlockProcessor = errors.New("nil block processor provided")

// ErrNilBlocksProcessor signals that a nil blocks processor has been provided
var ErrNilBlocksProcessor = errors.New("nil blocks processor provided")

// ErrNilProofProcessor signals that a nil proof processor has been provided
var ErrNilProofProcessor = errors.New("nil proof processor provided")

// ErrNilESDTSuppliesProcessor signals that a nil esdt supplies processor has been provided
var ErrNilESDTSuppliesProcessor = errors.New("nil esdt supplies processor")

// ErrNilStatusProcessor signals that a nil status processor has been provided
var ErrNilStatusProcessor = errors.New("nil status processor")

// ErrNilAboutInfoProcessor signals that a nil about info processor has been provided
var ErrNilAboutInfoProcessor = errors.New("nil about info processor")
