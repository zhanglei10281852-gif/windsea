package export_test

import (
	"github.com/zhanglei10281852-gif/windsea/internal/export"
	"testing"
)

func TestBundleRoundTrip(t *testing.T) {
	raw, err := export.BuildBundle([]export.Artifact{{Name: "summary.txt", Content: []byte("ready")}, {Name: "meta.json", Content: []byte("{}")}})
	if err != nil {
		t.Fatal(err)
	}
	files, err := export.ReadBundle(raw)
	if err != nil || string(files["summary.txt"]) != "ready" {
		t.Fatalf("files=%v err=%v", files, err)
	}
	if _, err := export.BuildBundle([]export.Artifact{{Name: "a", Content: nil}, {Name: "a", Content: nil}}); err == nil {
		t.Fatal("duplicate artifact accepted")
	}
}
