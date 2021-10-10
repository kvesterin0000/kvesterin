package db

import "testing"

func Test_hashPass(t *testing.T) {
	type args struct {
		pass string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default",
			args: args{pass: "password"},
			want: "X03MO1qnZdYdgyfeuILPmQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hashPass(tt.args.pass); got != tt.want {
				t.Errorf("hashPass() = %v, want %v", got, tt.want)
			}
		})
	}
}
