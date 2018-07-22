package bbir

import (
	"encoding/csv"
	"io"
	"strings"
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

	header, err := r.Read()
	if err == io.EOF {
		return len(lines), lines, nil
	} else if err != nil {
		return 0, []*Line{}, err
	}

	for i, name := range header {
		name = strings.TrimLeft(name, " ")
		name = strings.TrimRight(name, " ")
		header[i] = name
	}

	for i := 1; ; i++ {
		record, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, []*Line{}, err
		}
		lines = append(lines, NewLine(header, record))
	}
	return len(lines), lines, nil
}
