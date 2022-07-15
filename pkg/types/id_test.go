package types

import (
	"reflect"
	"testing"
)

//nolint:scopelint
func TestID_Marshal(t *testing.T) {
	type fields struct {
		ID int64
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				ID: 70579728276262912,
			},
			want:    []byte(`"av8FMdRdcb"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		id := NewID(tt.fields.ID)
		t.Run(tt.name, func(t *testing.T) {
			got, err := id.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//nolint:scopelint
func TestID_Unmarshal(t *testing.T) {
	type fields struct {
		id *ID
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantInt64 int64
		wantErr   bool
	}{
		{
			name: "test1",
			fields: fields{
				id: &ID{},
			},
			args: args{
				data: []byte(`"av8FMdRdcb"`),
			},
			wantInt64: 70579728276262912,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.id.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields.id.Int64() != tt.wantInt64 {
				t.Errorf("Unmarshal() got = %v, want %v", tt.fields.id.Int64(), tt.wantInt64)
			}
		})
	}
}
