package reflect_util

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertMapToStruct(t *testing.T) {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	type args struct {
		m map[string]any
		s Person
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "happy case",
			args: args{
				m: map[string]any{
					"name": "Dat",
					"age":  26,
				},
				s: Person{
					Name: "Dat",
					Age:  26,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConvertMapToStruct(tt.args.m, &tt.args.s); (err != nil) != tt.wantErr {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.args.s.Name, tt.args.m["name"])
			assert.Equal(t, tt.args.s.Age, tt.args.m["age"])
		})
	}
}

func TestCopyStruct(t *testing.T) {
	type source struct {
		A string
		B string
		C string
	}

	type destination struct {
		A string
		B string
		D string
	}
	type args struct {
		source      source
		destination destination
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    *destination
	}{
		{
			name: "happy case",
			args: args{
				source: source{
					A: "a",
					B: "b",
					C: "c",
				},
			},
			wantErr: false,
			want: &destination{
				A: "a",
				B: "b",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyStruct(&tt.args.source, &tt.args.destination)
			assert.Equal(t, tt.wantErr, err != nil)
			assert.True(t, reflect.DeepEqual(tt.want, &tt.args.destination))
		})
	}
}
