package topics

import (
	"bytes"
	"crypto/sha256"
	"unicode/utf8"

	"github.com/bsv-blockchain/go-sdk/overlay"
	"github.com/bsv-blockchain/go-sdk/overlay/topic"
	"github.com/bsv-blockchain/go-sdk/script"
)

type File struct {
	Hash    []byte `json:"hash"`
	Size    uint32 `json:"size"`
	Type    string `json:"type"`
	Content []byte `json:"-"`
	Name    string `json:"name,omitempty"`
}

type Inscription struct {
	File   *File             `json:"file,omitempty"`
	Parent *overlay.Outpoint `json:"parent,omitempty"`
	Bitcom map[string][]byte `json:"bitcom,omitempty"`
}

type InscriptionTopicManager struct {
	topic.BaseTopicManager
	Network overlay.Network
}

func (tm *InscriptionTopicManager) IdentifyAdmissableOutputs(subCtx overlay.SubmitContext) (*overlay.Admittance, error) {
	admit := &overlay.Admittance{
		OutputData: make(map[uint32]interface{}),
	}
	for vout, output := range subCtx.Tx.Outputs {
		if output.Satoshis != 1 {
			continue
		}
		scr := output.LockingScript

		if insc := tm.Parse(scr); insc != nil {
			admit.OutputsToAdmit = append(admit.OutputsToAdmit, uint32(vout))
			admit.OutputData[uint32(vout)] = insc
		}
	}
	return admit, nil
}

func (tm *InscriptionTopicManager) Parse(scr *script.Script) *Inscription {
	for pos := 0; pos < len(*scr); {
		startI := pos
		if op, err := scr.ReadOp(&pos); err != nil {
			break
		} else if pos > 2 && op.Op == script.OpDATA3 && bytes.Equal(op.Data, []byte("ord")) && (*scr)[startI-2] == 0 && (*scr)[startI-1] == script.OpIF {
			insc := &Inscription{
				File: &File{},
			}

		ordLoop:
			for {
				var field int
				var err error
				var op, op2 *script.ScriptChunk
				if op, err = scr.ReadOp(&pos); err != nil || op.Op > script.Op16 {
					return insc
				} else if op2, err = scr.ReadOp(&pos); err != nil || op2.Op > script.Op16 {
					return insc
				} else if op.Op > script.OpPUSHDATA4 && op.Op <= script.Op16 {
					field = int(op.Op) - 80
				} else if len(op.Data) == 1 {
					field = int(op.Data[0])
				} else if len(op.Data) > 1 {
					if add, err := script.NewAddressFromString(string(op.Data)); err == nil {
						if insc.Bitcom == nil {
							insc.Bitcom = make(map[string][]byte)
						}
						insc.Bitcom[add.AddressString] = op2.Data
					}
					continue
				}
				switch field {
				case 0:
					insc.File.Content = op2.Data
					break ordLoop
				case 1:
					if len(op2.Data) < 256 && utf8.Valid(op2.Data) {
						insc.File.Type = string(op2.Data)
					}
				case 3:
					insc.Parent = overlay.NewOutpointFromBytes(op2.Data)
				}

			}
			op, err := scr.ReadOp(&pos)
			if err != nil || op.Op != script.OpENDIF {
				return insc
			}

			insc.File.Size = uint32(len(insc.File.Content))
			hash := sha256.Sum256(insc.File.Content)
			insc.File.Hash = hash[:]

			return insc
		}
	}
	return nil
}
