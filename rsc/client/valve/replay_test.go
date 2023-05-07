//go:build !integration
// +build !integration

package valve

import (
	"testing"
)

func Test_getFileName(t *testing.T) {
	type args struct {
		replayURL string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get file name from replay url should success",
			args: args{
				replayURL: "http://replay153.valve.net/570/7132230434_1635105612.dem.bz2",
			},
			want: "7132230434_1635105612.dem.bz2",
		},
		{
			name: "get file name from replay url should success",
			args: args{
				replayURL: "http://replay153.valve.net/570/7132230435_1635105612.dem.bz2",
			},
			want: "7132230435_1635105612.dem.bz2",
		},
		{
			name: "get file name from replay url should success",
			args: args{
				replayURL: "http://replay153.valve.net/570/7132230455_1635105612.dem.bz2",
			},
			want: "7132230455_1635105612.dem.bz2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getFileName(tt.args.replayURL); got != tt.want {
				t.Errorf("getFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}
