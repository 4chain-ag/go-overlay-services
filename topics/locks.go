package topics

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"log"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/script"
)

type Lock struct {
	Address string `json:"address"`
	Until   uint32 `json:"until"`
}

type LockTopicManager struct {
	topic.BaseTopicManager
	Network overlay.Network
}

func (tm *LockTopicManager) IdentifyAdmissableOutputs(subCtx overlay.SubmitContext) (*overlay.Admittance, error) {
	admit := &overlay.Admittance{
		OutputData: make(map[uint32]interface{}),
	}
	for vout, output := range subCtx.Tx.Outputs {
		if output.Satoshis != 1 {
			continue
		}
		scr := output.LockingScript

		if lock := tm.Parse(scr); lock != nil {
			admit.OutputsToAdmit = append(admit.OutputsToAdmit, uint32(vout))
			admit.OutputData[uint32(vout)] = lock
		}
	}
	return admit, nil
}

func (p *LockTopicManager) Parse(scr *script.Script) *Lock {
	lockPrefixIndex := bytes.Index(*scr, LockPrefix)
	if lockPrefixIndex > -1 && bytes.Contains((*scr)[lockPrefixIndex:], LockSuffix) {
		lock := &Lock{}
		pos := lockPrefixIndex + len(LockPrefix)
		if op, err := scr.ReadOp(&pos); err != nil {
			log.Println(err)
		} else if len(op.Data) != 20 {
			return nil
		} else if add, err := script.NewAddressFromPublicKeyHash(op.Data, p.Network == overlay.NetworkMainnet); err != nil {
			return nil
		} else {
			lock.Address = add.AddressString
		}
		if op, err := scr.ReadOp(&pos); err != nil {
			log.Println(err)
		} else {
			until := make([]byte, 4)
			copy(until, op.Data)
			lock.Until = binary.LittleEndian.Uint32(until)
		}
		return lock
	}
	return nil
}

var LockPrefix, _ = hex.DecodeString("2097dfd76851bf465e8f715593b217714858bbe9570ff3bd5e33840a34e20ff0262102ba79df5f8ae7604a9830f03c7933028186aede0675a16f025dc4f8be8eec0382201008ce7480da41702918d1ec8e6849ba32b4d65b1e40dc669c31a1e6306b266c0000")
var LockSuffix, _ = hex.DecodeString("610079040065cd1d9f690079547a75537a537a537a5179537a75527a527a7575615579014161517957795779210ac407f0e4bd44bfc207355a778b046225a7068fc59ee7eda43ad905aadbffc800206c266b30e6a1319c66dc401e5bd6b432ba49688eecd118297041da8074ce081059795679615679aa0079610079517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e01007e81517a75615779567956795679567961537956795479577995939521414136d08c5ed2bf3ba048afe6dcaebafeffffffffffffffffffffffffffffff00517951796151795179970079009f63007952799367007968517a75517a75517a7561527a75517a517951795296a0630079527994527a75517a6853798277527982775379012080517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f517f7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e7c7e01205279947f7754537993527993013051797e527e54797e58797e527e53797e52797e57797e0079517a75517a75517a75517a75517a75517a75517a75517a75517a75517a75517a75517a75517a756100795779ac517a75517a75517a75517a75517a75517a75517a75517a75517a7561517a75517a756169557961007961007982775179517954947f75517958947f77517a75517a756161007901007e81517a7561517a7561040065cd1d9f6955796100796100798277517951790128947f755179012c947f77517a75517a756161007901007e81517a7561517a756105ffffffff009f69557961007961007982775179517954947f75517958947f77517a75517a756161007901007e81517a7561517a75615279a2695679a95179876957795779ac7777777777777777")
