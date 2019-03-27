package synth

import (
	"reflect"
	"testing"
)

func Test_cartesianProduct(t *testing.T) {
	type args struct {
		lists [][]interface{}
	}
	tests := []struct {
		name string
		args args
		want [][]interface{}
	}{
		{
			name: "3 * 2",
			args: args{
				lists: [][]interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5},
				},
			},
			want: [][]interface{}{
				[]interface{}{1, 4},
				[]interface{}{2, 4},
				[]interface{}{3, 4},
				[]interface{}{1, 5},
				[]interface{}{2, 5},
				[]interface{}{3, 5},
			},
		},
		{
			name: "3 * 2 * 2",
			args: args{
				lists: [][]interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5},
					[]interface{}{7, 8},
				},
			},
			want: [][]interface{}{
				[]interface{}{1, 4, 7},
				[]interface{}{2, 4, 7},
				[]interface{}{3, 4, 7},
				[]interface{}{1, 5, 7},
				[]interface{}{2, 5, 7},
				[]interface{}{3, 5, 7},
				[]interface{}{1, 4, 8},
				[]interface{}{2, 4, 8},
				[]interface{}{3, 4, 8},
				[]interface{}{1, 5, 8},
				[]interface{}{2, 5, 8},
				[]interface{}{3, 5, 8},
			},
		},
		{
			name: "3",
			args: args{
				lists: [][]interface{}{
					[]interface{}{1, 2, 3},
				},
			},
			want: [][]interface{}{
				[]interface{}{1},
				[]interface{}{2},
				[]interface{}{3},
			},
		},
		{
			name: "(empty)",
			args: args{
				lists: [][]interface{}{},
			},
			want: [][]interface{}{
				[]interface{}{},
			},
		},
		{
			name: "3 * 2 * 0",
			args: args{
				lists: [][]interface{}{
					[]interface{}{1, 2, 3},
					[]interface{}{4, 5},
					[]interface{}{},
				},
			},
			want: [][]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cartesianProduct(tt.args.lists); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cartesianProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}
