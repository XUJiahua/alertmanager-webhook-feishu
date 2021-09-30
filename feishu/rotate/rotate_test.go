package rotate

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func TestEveryN(t *testing.T) {
	type args struct {
		i int
		n int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				i: 1,
				n: 2,
			},
			want: 1,
		},
		{
			args: args{
				i: 2,
				n: 2,
			},
			want: 1,
		},
		{
			args: args{
				i: 3,
				n: 2,
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bucketIndexEveryN(tt.args.i, tt.args.n); got != tt.want {
				t.Errorf("bucketIndexEveryN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMentionRotator_Rotate(t *testing.T) {
	type fields struct {
		BaseDate  time.Time
		CycleDays int
		OpenIDs   []string
	}
	type args struct {
		t time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			fields: fields{
				BaseDate:  time.Now(),
				CycleDays: 14,
				OpenIDs:   []string{"a", "b"},
			},
			args: args{
				t: time.Now(),
			},
			want: []string{"a"},
		},
		{
			fields: fields{
				BaseDate:  time.Now(),
				CycleDays: 14,
				OpenIDs:   []string{"a", "b"},
			},
			args: args{
				t: time.Now().AddDate(0, 0, 13),
			},
			want: []string{"a"},
		},
		{
			fields: fields{
				BaseDate:  time.Now(),
				CycleDays: 14,
				OpenIDs:   []string{"a", "b"},
			},
			args: args{
				t: time.Now().AddDate(0, 0, 14),
			},
			want: []string{"b"},
		},
		{
			fields: fields{
				BaseDate:  time.Now(),
				CycleDays: 14,
				OpenIDs:   []string{"a", "b"},
			},
			args: args{
				t: time.Now().AddDate(0, 0, -1),
			},
			want: []string{"b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MentionRotator{
				baseDate:  tt.fields.BaseDate,
				cycleDays: tt.fields.CycleDays,
				openIDs:   tt.fields.OpenIDs,
			}
			if got := r.Rotate(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rotate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_adjustDays(t *testing.T) {
	type args struct {
		relativeDays int
		cycleDays    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				relativeDays: 0,
				cycleDays:    2,
			},
			want: 1,
		},
		{
			args: args{
				relativeDays: 0,
				cycleDays:    1,
			},
			want: 1,
		},
		{
			args: args{
				relativeDays: 0,
				cycleDays:    1,
			},
			want: 1,
		},
		{
			args: args{
				relativeDays: -1,
				cycleDays:    2,
			},
			want: 3,
		},
		{
			args: args{
				relativeDays: -3,
				cycleDays:    2,
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := adjustDays(tt.args.relativeDays, tt.args.cycleDays); got != tt.want {
				t.Errorf("adjustDays() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseDays(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			args: args{
				s: "1w1d",
			},
			want: 8,
		},
		{
			args: args{
				s: "1w",
			},
			want: 7,
		},
		{
			args: args{
				s: "1d",
			},
			want: 1,
		},
		{
			args: args{
				s: "1d1w",
			},
			wantErr: true,
		},
		{
			args: args{
				s: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDuration(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseDuration() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	rotator, err := New("2021-09-27:2w", nil)
	require.NoError(t, err)
	spew.Dump(rotator)
	spew.Dump(time.Now())
	spew.Dump(time.Now().Format(time.RFC3339))
	spew.Dump(time.Now().Format("Z07:00"))
}
