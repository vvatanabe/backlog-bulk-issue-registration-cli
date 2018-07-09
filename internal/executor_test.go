package internal

import (
	"os"
	"testing"
)

func Test_BulkCommandExecutor(t *testing.T) {
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
	commands, err := converter.Convert(lines)
	if err != nil {
		t.Errorf("Could not convert:\n%s", err.Error())
	}

	executor := injector.Get(new(BulkCommandExecutor)).(*BulkCommandExecutor)
	if err := executor.Do(commands); err != nil {
		t.Errorf("Could not execute:\n%s", err.Error())
	}
}
