package bbir

import (
	"os"
	"testing"
)

const fixturesPath = "../testdata/"

func Test_CommandConverter(t *testing.T) {
	injector := NewInjectorForCommandBuilderTest(t)
	converter := injector.Get(new(CommandConverter)).(*CommandConverter)
	filePath := fixturesPath + "example.csv"
	file, err := os.Open(filePath)
	if err != nil {
		t.Errorf("Could not open fixture file: %v", filePath)
	}
	_, lines, err := NewCSVReader(file).ReadAll()
	if err != nil {
		t.Errorf("Could not read csv data: %v", err.Error())
	}
	_, err = converter.Convert(lines)
	if err != nil {
		t.Errorf("Could not convert:\n%s", err.Error())
	}
}
