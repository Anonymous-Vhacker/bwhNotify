package util

import "testing"

func TestGenerateDingTalkSign(t *testing.T) {
	type args struct {
		timestamp int64
		secret    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: struct {
				timestamp int64
				secret    string
			}{
				timestamp: 1720438284166,
				secret:    "SEC1da7326e83ea5dc27bfa71eed304d4f93803fdae3487fd0e54633fbe9f38bbc5",
			},
			want: "Kc6ieSszXt1GPCgiiRXHyezV2Yz8knM6GM8Xj9tJm8Y=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateDingTalkSign(tt.args.timestamp, tt.args.secret); got != tt.want {
				t.Errorf("GenerateDingTalkSign() = %v, want %v", got, tt.want)
			}
		})
	}
}
