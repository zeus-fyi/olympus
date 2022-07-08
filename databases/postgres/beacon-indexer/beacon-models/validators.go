package beacon_models

import (
	"context"
)

type Validator struct {
	Index  string
	Pubkey string
}

var insertValidators = "INSERT INTO validators (index, pubkey) VALUES"

func (v *Validator) GetFieldValues() []string {
	return []string{v.Index, v.Pubkey}
}

func (v *Validator) InsertValidators(ctx context.Context, vs ...Validator) error {

	//for _, val := range vs {
	//	val.Index
	//}
	return nil
}
