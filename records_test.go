package main

import (
	"testing"
)

var WKM = Record{
	Artist:      "Tom Misch",
	Album:       "What Kinda Music",
	AmazonPrice: "£30",
}

var LF = Record{
	Artist:      "Jorja Smith",
	Album:       "Lost & Found",
	AmazonPrice: "£22.75",
}

var NWBD = Record{
	Artist:      "Loyle Carner",
	Album:       "Not Waving, But Drowning",
	AmazonPrice: "£25",
}

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
