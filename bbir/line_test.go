package bbir

import "testing"

func Test_Line_NewLine(t *testing.T) {

	header := []string{
		"Summary",
		"Description",
		"StartDate",
		"DueDate",
		"EstimatedHours",
		"ActualHours",
		"IssueType",
		"Category",
		"Version",
		"Milestone",
		"Priority",
		"Assignee",
		"ParentIssue",
		"CustomCheckBox1",
		"CustomCheckBox2",
	}

	record := []string{
		"Summary",
		"Description",
		"2018-01-01",
		"2018-01-02",
		"1",
		"2",
		"task",
		"web",
		"sprint1",
		"sprint1",
		"Normal",
		"ken",
		"*",
		"1",
		"2",
	}
	line := NewLine(header, record)
	if record[0] != line.Summary {
		t.Errorf("Could not match a value. want: %s, result: %s", record[0], line.Summary)
	}
	if record[1] != line.Description {
		t.Errorf("Could not match a value. want: %s, result: %s", record[1], line.Description)
	}
	if record[2] != line.StartDate {
		t.Errorf("Could not match a value. want: %s, result: %s", record[2], line.StartDate)
	}
	if record[3] != line.DueDate {
		t.Errorf("Could not match a value. want: %s, result: %s", record[3], line.DueDate)
	}
	if record[4] != line.EstimatedHours {
		t.Errorf("Could not match a value. want: %s, result: %s", record[4], line.EstimatedHours)
	}
	if record[5] != line.ActualHours {
		t.Errorf("Could not match a value. want: %s, result: %s", record[5], line.ActualHours)
	}
	if record[6] != line.IssueType {
		t.Errorf("Could not match a value. want: %s, result: %s", record[6], line.IssueType)
	}
	if record[7] != line.Category {
		t.Errorf("Could not match a value. want: %s, result: %s", record[7], line.Category)
	}
	if record[8] != line.Version {
		t.Errorf("Could not match a value. want: %s, result: %s", record[8], line.Version)
	}
	if record[9] != line.Milestone {
		t.Errorf("Could not match a value. want: %s, result: %s", record[9], line.Milestone)
	}
	if record[10] != line.Priority {
		t.Errorf("Could not match a value. want: %s, result: %s", record[10], line.Priority)
	}
	if record[11] != line.Assignee {
		t.Errorf("Could not match a value. want: %s, result: %s", record[11], line.Assignee)
	}
	if record[12] != line.ParentIssue {
		t.Errorf("Could not match a value. want: %s, result: %s", record[12], line.ParentIssue)
	}
	if record[13] != line.CustomFields[header[13]] {
		t.Errorf("Could not match a value. want: %s, result: %s", record[13], line.CustomFields[header[13]])
	}
	if record[14] != line.CustomFields[header[14]] {
		t.Errorf("Could not match a value. want: %s, result: %s", record[14], line.CustomFields[header[14]])
	}
}
