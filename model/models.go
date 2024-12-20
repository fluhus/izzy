// Concrete models.

package model

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// BUG(amit): Make JSON parsing lazy (upon model request) and test all parsing
// in a unit-test?

//go:generate go run ./gen/gen.go

var (
	HiSeqModel, NovaSeqModel, MiSeqModel, BasicModel, PerfectModel *Model
)

var (
	//go:embed HiSeq.json
	hiseqData []byte
	//go:embed NovaSeq.json
	novaseqData []byte
	//go:embed MiSeq.json
	miseqData []byte
	//go:embed Basic.json
	basicData []byte
	//go:embed Perfect.json
	perfectData []byte
)

func init() {
	HiSeqModel = &Model{}
	if err := json.Unmarshal(hiseqData, HiSeqModel); err != nil {
		panic(fmt.Sprintf("failed to create HiSeq model: %v", err))
	}
	NovaSeqModel = &Model{}
	if err := json.Unmarshal(novaseqData, NovaSeqModel); err != nil {
		panic(fmt.Sprintf("failed to create NovaSeq model: %v", err))
	}
	MiSeqModel = &Model{}
	if err := json.Unmarshal(miseqData, MiSeqModel); err != nil {
		panic(fmt.Sprintf("failed to create MiSeq model: %v", err))
	}
	BasicModel = &Model{}
	if err := json.Unmarshal(basicData, BasicModel); err != nil {
		panic(fmt.Sprintf("failed to create basic model: %v", err))
	}
	PerfectModel = &Model{}
	if err := json.Unmarshal(perfectData, PerfectModel); err != nil {
		panic(fmt.Sprintf("failed to create perfect model: %v", err))
	}
}
