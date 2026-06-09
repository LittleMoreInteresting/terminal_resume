package data

import (
	"testing"
)

func TestLoadResume(t *testing.T) {
	resume, err := LoadResume()
	if err != nil {
		t.Fatalf("LoadResume failed: %v", err)
	}

	if resume.Name == "" {
		t.Error("expected Name to be non-empty")
	}
	if resume.Title == "" {
		t.Error("expected Title to be non-empty")
	}
	if len(resume.Experience) == 0 {
		t.Error("expected Experience to be non-empty")
	}
	if len(resume.Projects) == 0 {
		t.Error("expected Projects to be non-empty")
	}

	t.Logf("Loaded resume for: %s (%s)", resume.Name, resume.Title)
	t.Logf("Experience count: %d", len(resume.Experience))
	t.Logf("Projects count: %d", len(resume.Projects))
}
