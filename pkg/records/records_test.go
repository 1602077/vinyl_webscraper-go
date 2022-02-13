package records

import (
	"testing"
)

var WKM = NewRecord("Tom Misch", "What Kinda Music", "", float64(30))
var LF = NewRecord("Jorja Smith", "Lost & Found", "", float64(22.75))
var NWBD = NewRecord("Loyle Carner", "Not Waving, But Drowning", "", float64(25))

func TestMergeRecordHistories(t *testing.T) {
	t.Run("it merges a RecordInstance with Record History", func(t *testing.T) {
		var TestRecords = Records{WKM, LF}
		var rh = RecordHistory{RecordInstance{Date: "yesterday", Records: TestRecords}}
		var ri = RecordInstance{Date: "today", Records: TestRecords}

		rh.MergeRecordHistories(ri)

		got := len(rh)
		want := 2

		if got != want {
			t.Errorf("expected length %v, got %v", want, got)
		}
	})

	t.Run("it replaces duplicate date with RecordInstance being merged in", func(t *testing.T) {
		var rh = RecordHistory{RecordInstance{Date: "today", Records: Records{WKM, LF}}}
		var ri = RecordInstance{Date: "today", Records: Records{NWBD}}

		rh.MergeRecordHistories(ri)

		if len(rh) != 1 {
			t.Errorf("expected 1 date of data, got %v", len(rh))
		}

		if len(ri.Records) != 1 {
			t.Errorf("expected 1 record in MergedRecordHistory, got %v", len(ri.Records))
		}
	})
}
