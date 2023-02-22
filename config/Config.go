package config

import "github.com/pingcap/go-ycsb/pkg/generator"

type ProcessorConfig struct {
}

type ProposerConfig struct {
	Zipfian *generator.Zipfian
}

type ValidatorConfig struct {
}

type ConcurrentProcessorConfig struct {
	threadnumber uint64
}

type ChainConfig struct {
	ChainId uint64
}
type Config interface {
}

type Hash [32]byte

const ContractAddress = "0x0000000000000000000000000000000000000000"

const Blocksize = 1000

const SinglechainSize = 32

const MaxHeaderData = 32

const LeadingZeros = 1

const MaxBlockExtraData = 100 * 1024

const SmallbankAccessedkeys = 4

const Txsperblock = 100

//const Concurrentprocessor = 0
//const Sequentialprocessor = 1
//const Preprocessor = 2
//
//const Concurrent = true
//const Sequential = false
