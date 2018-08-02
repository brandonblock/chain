package chain

import (
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestBlock_SetHash(t *testing.T) {
	emptyHash := "[95 236 235 102 255 200 111 56 217 82 120 108 109 105 108 121 194 219 194 57 221 78 145 180 103 41 215 58 39 251 87 233]"
	type fields struct {
		Timestamp     int64
		Data          []byte
		PrevBlockHash []byte
		Hash          []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Empty Block",
			fields: fields{
				Timestamp:     0,
				Data:          []byte{},
				PrevBlockHash: []byte{},
				Hash:          []byte{},
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
			}
			b.SetHash()
			assert.Equal()
		})
	}
}
