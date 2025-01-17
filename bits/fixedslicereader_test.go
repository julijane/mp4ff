package bits

import (
	"fmt"
	"testing"
)

func TestFixedSliceReaderInts(t *testing.T) {
	nrBytes := 64
	data := make([]byte, nrBytes)
	for i := 0; i < nrBytes; i++ {
		data[i] = byte(i)
	}
	sr := NewFixedSliceReader(data)
	if sr.Length() != nrBytes {
		t.Errorf("sliceReader length = %d instead of %d", sr.Length(), nrBytes)
	}
	for i := 0; i < nrBytes+2; i++ {
		val := sr.ReadUint8()
		switch {
		case i < nrBytes:
			if i != int(val) {
				t.Errorf("val at %d not %d but %d", i, i, val)
			}
		default: // Read passed end
			if sr.AccError() == nil {
				t.Errorf("should have had an accumulated error")
			}
			if 0 != val {
				t.Errorf("should have zero value")
			}
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals := nrBytes / 2
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadUint16()
		switch {
		case i < nrVals:
			expectedVal := uint16(expVal(i, 2))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 2
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadInt16()
		switch {
		case i < nrVals:
			expectedVal := int16(expVal(i, 2))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 3
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadUint24()
		switch {
		case i < nrVals:
			expectedVal := uint32(expVal(i, 3))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 4
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadUint32()
		switch {
		case i < nrVals:
			expectedVal := uint32(expVal(i, 4))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 4
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadInt32()
		switch {
		case i < nrVals:
			expectedVal := int32(expVal(i, 4))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 8
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadUint64()
		switch {
		case i < nrVals:
			expectedVal := uint64(expVal(i, 8))
			if val != expectedVal {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}

	sr = NewFixedSliceReader(data)
	nrVals = nrBytes / 8
	for i := 0; i < nrVals+2; i++ {
		val := sr.ReadInt64()
		switch {
		case i < nrVals:
			expectedVal := int64(expVal(i, 8))
			if val != int64(expVal(i, 8)) {
				t.Errorf("val at %d not %d but %d", i, expectedVal, val)
			}
		default: // Read passed end
			verifyAccErrorInt(t, sr, int(val))
		}
	}
}

func verifyAccErrorInt(t *testing.T, sr *FixedSliceReader, val int) {
	if sr.AccError() == nil {
		t.Errorf("should have had an accumulated error")
	}
	if 0 != val {
		t.Errorf("should have zero value")
	}
}

func verifyAccErrorString(t *testing.T, sr *FixedSliceReader, val string) {
	if sr.AccError() == nil {
		t.Errorf("should have had an accumulated error")
	}
	if "" != val {
		t.Errorf("should be empty string")
	}
}

func verifyAccErrorBytes(t *testing.T, sr *FixedSliceReader, val []byte) {
	if sr.AccError() == nil {
		t.Errorf("should have had an accumulated error")
	}
	if len(val) != 0 {
		t.Errorf("should be empty slice")
	}
}

func expVal(i, nrBytes int) (val uint64) {
	for j := 0; j < nrBytes; j++ {
		byteVal := uint64(i*nrBytes + j)
		val += byteVal << (8 * (nrBytes - 1 - j))
	}
	return val
}

func TestFixedSliceReaderStrings(t *testing.T) {
	data := []byte("0123456789\x00abcdef")
	nrBytes := len(data)
	sr := NewFixedSliceReader(data)
	if sr.Length() != nrBytes {
		t.Errorf("sliceReader length = %d instead of %d", sr.Length(), nrBytes)
	}
	val := sr.ReadFixedLengthString(4)
	if val != "0123" {
		t.Errorf(`read string is %q instead of "0123"`, val)
	}
	val = sr.ReadZeroTerminatedString(10)
	if val != "456789" {
		t.Errorf(`read string is %q instead of "456789"`, val)
	}
	val = sr.ReadZeroTerminatedString(2)
	verifyAccErrorString(t, sr, val)

	val = sr.ReadFixedLengthString(2)
	verifyAccErrorString(t, sr, val)

	val = sr.ReadZeroTerminatedString(2)
	verifyAccErrorString(t, sr, val)

	sr = NewFixedSliceReader(data)
	val = sr.ReadFixedLengthString(nrBytes + 2)
	verifyAccErrorString(t, sr, val)
}

func TestFixedSliceReadBytes(t *testing.T) {
	data := []byte("0123456789abcdef")
	nrBytes := len(data)
	sr := NewFixedSliceReader(data)
	if sr.Length() != nrBytes {
		t.Errorf("sliceReader length = %d instead of %d", sr.Length(), nrBytes)
	}
	val := sr.ReadBytes(4)
	if string(val) != "0123" {
		t.Errorf(`read bytes are %q instead of "0123"`, string(val))
	}
	if sr.NrRemainingBytes() != nrBytes-4 {
		t.Errorf("nr remaining = %d instead of %d", sr.NrRemainingBytes(), nrBytes-4)
	}
	sr.SkipBytes(2)
	if sr.NrRemainingBytes() != nrBytes-6 {
		t.Errorf("nr remaining = %d instead of %d", sr.NrRemainingBytes(), nrBytes-6)
	}
	sr.SetPos(8)
	if sr.NrRemainingBytes() != nrBytes-8 {
		t.Errorf("nr remaining = %d instead of %d", sr.NrRemainingBytes(), nrBytes-8)
	}
	lookAhead := make([]byte, 8)
	err := sr.LookAhead(0, lookAhead)
	if err != nil {
		t.Error(err)
	}
	if string(lookAhead) != "89abcdef" {
		t.Errorf(`lookahead %q instead of "89abcdef"`, lookAhead)
	}
	lookAhead = make([]byte, 9)
	err = sr.LookAhead(0, lookAhead)
	if err == nil {
		t.Errorf("should be out of bounds")
	}
	remaining := sr.RemainingBytes()
	if string(remaining) != "89abcdef" {
		t.Errorf(`remaining %q instead of "89abcdef"`, remaining)
	}

	// Beyond end
	val = sr.ReadBytes(2)
	verifyAccErrorBytes(t, sr, val)

	val = sr.ReadBytes(2)
	verifyAccErrorBytes(t, sr, val)

	val = sr.RemainingBytes()
	verifyAccErrorBytes(t, sr, val)

	valNr := sr.NrRemainingBytes()
	verifyAccErrorInt(t, sr, valNr)

	sr.SkipBytes(2)
	verifyAccErrorInt(t, sr, 0)

	sr = NewFixedSliceReader(data)
	sr.SkipBytes(nrBytes + 2)
	verifyAccErrorInt(t, sr, 0)

	sr.SetPos(nrBytes + 2)
	wantedErrMsg := fmt.Sprintf("attempt to set pos %d beyond slice len %d", nrBytes+2, nrBytes)
	if sr.AccError().Error() != wantedErrMsg {
		t.Errorf("got error msg %q instead of %q", sr.AccError().Error(), wantedErrMsg)
	}
}
