package trigger

import (
	"fmt"
	"testing"
	"time"
)

func TestTimeTrigger_Start(t *testing.T) {
	type fields struct {
		targetTime time.Time
	}
	target := time.Now().Add(3 * time.Second)
	tests := []struct {
		name    string
		fields  fields
		want    <-chan time.Time
		want1   chan error
		wantErr bool
	}{
		{
			name: "time",
			fields: fields{
				targetTime: target,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &TimeTrigger{
				targetTime: tt.fields.targetTime,
			}
			got, _, _ := b.Listen()
			tm := <-got
			fmt.Println(tm.Sub(target).Milliseconds())
		})
	}
}
