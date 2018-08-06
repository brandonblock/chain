package chain

import (
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

func TestBlock_Serialize(t *testing.T) {
	type fields struct {
		Timestamp     int64
		Data          []byte
		PrevBlockHash []byte
		Hash          []byte
		Nonce         int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "no error please",
			fields: fields{
				Timestamp: time.Now().Unix(),
				Data:      []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Block{
				Timestamp:     tt.fields.Timestamp,
				Data:          tt.fields.Data,
				PrevBlockHash: tt.fields.PrevBlockHash,
				Hash:          tt.fields.Hash,
				Nonce:         tt.fields.Nonce,
			}
			got, err := b.Serialize()
			assert.NoError(t, err, "this should be error-free")
			assert.NotNil(t, got, "if there's no error, this shouldn't be nil")
		})
	}
}

func TestDeserializeBlock(t *testing.T) {

	b := &Block{
		Timestamp: time.Now().Unix(),
		Data:      []byte{},
	}

	encoded, _ := b.Serialize()

	type args struct {
		d []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				d: encoded,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DeserializeBlock(tt.args.d)
			assert.NoError(t, err, "this should be error-free")
			assert.IsType(t, &Block{}, got, "return should be of type *Block")
		})
	}
}
