package modules

import "testing"

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a       int
		b       int
		want    int
		wantErr bool
	}{
		{
			name:    "positive - simple division",
			a:       10,
			b:       2,
			want:    5,
			wantErr: false,
		},
		{
			name:    "positive - zero numerator",
			a:       0,
			b:       5,
			want:    0,
			wantErr: false,
		},
		{
			name:    "positive - both positive",
			a:       20,
			b:       4,
			want:    5,
			wantErr: false,
		},
		{
			name:    "negative - negative minus positive",
			a:       -10,
			b:       2,
			want:    -5,
			wantErr: false,
		},
		{
			name:    "negative - both negative",
			a:       -20,
			b:       -4,
			want:    5,
			wantErr: false,
		},
		{
			name:    "error - division by zero",
			a:       5,
			b:       0,
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Divide(%d, %d) error = %v, wantErr %v", tt.a, tt.b, err, tt.wantErr)
				return
			}
			if err == nil && got != tt.want {
				t.Errorf("Divide(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "both positive",
			a:    10,
			b:    3,
			want: 7,
		},
		{
			name: "positive minus zero",
			a:    5,
			b:    0,
			want: 5,
		},
		{
			name: "negative minus positive",
			a:    -5,
			b:    3,
			want: -8,
		},
		{
			name: "both negative",
			a:    -5,
			b:    -3,
			want: -2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
