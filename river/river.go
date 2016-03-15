package river

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
)

// Non-optimized way to create a key.
func SeriesKey(measurement string, tags map[string]string) []byte {
	var b bytes.Buffer
	b.WriteString(measurement)

	keys := make([]string, 0, len(tags))
	for k, _ := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		b.WriteByte(',')

		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(tags[k])
	}

	return b.Bytes()
}

// WriteLine writes the line represented by seriesKey, fields, and time, to w.
// Returns any error returned during write.
func WriteLine(w io.Writer, seriesKey []byte, fields []Field, time int64) error {
	var buf bytes.Buffer // TODO: use sync.pool?

	buf.Write(seriesKey)
	buf.WriteByte(' ')

	for i, f := range fields {
		if i != 0 {
			buf.WriteByte(',')
		}

		f.writeToBuf(&buf)
	}
	buf.WriteByte(' ')

	// Timestamp in nanoseconds, formatted in base 10, should fit in exactly 19 bytes for the foreseeable future
	tsBuf := make([]byte, 0, 19)
	tsBuf = strconv.AppendInt(tsBuf, time, 10)
	buf.Write(tsBuf)

	if _, err := io.Copy(w, &buf); err != nil {
		return fmt.Errorf("Error writing line: %s", err.Error())
	}

	return nil
}