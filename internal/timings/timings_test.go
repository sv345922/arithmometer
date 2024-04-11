package timings

import (
	"arithmometer/internal/entities"
	"testing"
)

var times = entities.Timings{
	1, 2, 3, 4,
}

func TestTimings_GetDuration(t *testing.T) {
	var timings = new(Timings)
	timings.Timings = times

	if tm := timings.GetDuration("+"); tm != T_const*1 {
		t.Errorf("GetDuratuin(+)=%v, wont %v", tm, T_const*1)
	}
	if tm := timings.GetDuration("-"); tm != T_const*2 {
		t.Errorf("GetDuratuin(+)=%v, wont %v", tm, T_const*2)
	}
	if tm := timings.GetDuration("*"); tm != T_const*3 {
		t.Errorf("GetDuratuin(+)=%v, wont %v", tm, T_const*3)
	}
	if tm := timings.GetDuration("/"); tm != T_const*4 {
		t.Errorf("GetDuratuin(+)=%v, wont %v", tm, T_const*4)
	}
	if tm := timings.GetDuration("invalid symbol"); tm != T_const*0 {
		t.Errorf("GetDuratuin(+)=%v, wont %v", tm, T_const*0)
	}
}
func TestTimings_String(t *testing.T) {
	var timings = new(Timings)
	timings.Timings = times

	if timings.String() != "+: 1s, -: 2s, *: 3s, /: 4s" {
		t.Errorf("invalid timings.String()")
	}
}
