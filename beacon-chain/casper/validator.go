package casper

import (
	"fmt"

	"github.com/ovcharovvladimir/essentiaHybrid/common"
	"github.com/ovcharovvladimir/Prysm/beacon-chain/params"
	"github.com/ovcharovvladimir/Prysm/beacon-chain/utils"
	pb "github.com/ovcharovvladimir/Prysm/proto/beacon/p2p/v1"
)

// RotateValidatorSet is called every dynasty transition. The primary functions are:
// 1.) Go through queued voter indices and induct them to be active by setting start
// dynasty to current cycle.
// 2.) Remove bad active voter whose balance is below threshold to the exit set by
// setting end dynasty to current cycle.
func RotateValidatorSet(validators []*pb.ValidatorRecord, dynasty uint64) []*pb.ValidatorRecord {
	upperbound := len(ActiveValidatorIndices(validators, dynasty))/30 + 1

	// Loop through active voter set, remove voter whose balance is below 50%.
	for _, index := range ActiveValidatorIndices(validators, dynasty) {
		if validators[index].Balance < params.DefaultBalance/2 {
			validators[index].EndDynasty = dynasty
		}
	}
	// Get the total number of voter we can induct.
	inductNum := upperbound
	if len(QueuedValidatorIndices(validators, dynasty)) < inductNum {
		inductNum = len(QueuedValidatorIndices(validators, dynasty))
	}

	// Induct queued voter to active voter set until the switch dynasty is greater than current number.
	for _, index := range QueuedValidatorIndices(validators, dynasty) {
		validators[index].StartDynasty = dynasty
		inductNum--
		if inductNum == 0 {
			break
		}
	}
	return validators
}

// ActiveValidatorIndices filters out active validators based on start and end dynasty
// and returns their indices in a list.
func ActiveValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []uint32 {
	var indices []uint32
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty <= dynasty && dynasty < validators[i].EndDynasty {
			indices = append(indices, uint32(i))
		}
	}
	return indices
}

// ExitedValidatorIndices filters out exited validators based on start and end dynasty
// and returns their indices in a list.
func ExitedValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []uint32 {
	var indices []uint32
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty < dynasty && validators[i].EndDynasty <= dynasty {
			indices = append(indices, uint32(i))
		}
	}
	return indices
}

// QueuedValidatorIndices filters out queued validators based on start and end dynasty
// and returns their indices in a list.
func QueuedValidatorIndices(validators []*pb.ValidatorRecord, dynasty uint64) []uint32 {
	var indices []uint32
	for i := 0; i < len(validators); i++ {
		if validators[i].StartDynasty > dynasty {
			indices = append(indices, uint32(i))
		}
	}
	return indices
}

// SampleAttestersAndProposers returns lists of random sampled attesters and proposer indices.
func SampleAttestersAndProposers(seed common.Hash, validators []*pb.ValidatorRecord, dynasty uint64) ([]uint32, uint32, error) {
	attesterCount := params.MinCommiteeSize
	if len(validators) < params.MinCommiteeSize {
		attesterCount = len(validators)
	}
	indices, err := utils.ShuffleIndices(seed, ActiveValidatorIndices(validators, dynasty))
	if err != nil {
		return nil, 0, err
	}
	return indices[:int(attesterCount)], indices[len(indices)-1], nil
}

// GetAttestersTotalDeposit from the pending attestations.
func GetAttestersTotalDeposit(attestations []*pb.AttestationRecord) uint64 {
	var numOfBits int
	for _, attestation := range attestations {
		for _, byte := range attestation.AttesterBitfield {
			numOfBits += int(utils.BitSetCount(byte))
		}
	}
	// Assume there's no slashing condition, the following logic will change later phase.
	return uint64(numOfBits) * params.DefaultBalance
}

// GetIndicesForHeight returns the attester set of a given height.
func GetIndicesForHeight(shardCommittees []*pb.ShardAndCommitteeArray, lcs uint64, height uint64) (*pb.ShardAndCommitteeArray, error) {
	if !(lcs <= height && height < lcs+params.CycleLength*2) {
		return nil, fmt.Errorf("can not return attester set of given height, input height %v has to be in between %v and %v", height, lcs, lcs+params.CycleLength*2)
	}
	return shardCommittees[height-lcs], nil
}
