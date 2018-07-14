package bbir

import (
	"encoding/csv"
	"io"
)

func NewCSVReader(r io.Reader) *CSVReader {
	reader := csv.NewReader(r)
	reader.Comma = ','
	reader.LazyQuotes = true
	return &CSVReader{reader}
}

type CSVReader struct {
	*csv.Reader
}

func (r *CSVReader) ReadAll() (size int, lines []*Line, err error) {
	for i := 0; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, []*Line{}, err
		}
		if i == 0 {
			continue
		}
		lines = append(lines, NewLine(record))
	}
	return len(lines), lines, nil
}
