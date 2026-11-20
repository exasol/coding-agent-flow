package openfasttrace

import (
	"reflect"
	"strings"
	"testing"

	"github.com/exasol/coding-agent-flow/internal/oftvalidation/domain"
)

func parseID(t *testing.T, text string) domain.OFTSpecificationItemID {
	t.Helper()
	id, err := domain.ParseOFTSpecificationItemID(text)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// [utest -> req~oft-requirements-detection~1]
func TestReadASpecReportReturnsTheOFTSpecificationItemsWithTheirCoverage(t *testing.T) {
	report := []byte(`<specdocument>
  <specobjects doctype="req">
    <specobject><id>r</id><version>1</version>
      <coverage><coveringSpecObjects></coveringSpecObjects></coverage>
      <covering><coveredType><id>f</id><version>1</version><doctype>feat</doctype></coveredType></covering>
    </specobject>
  </specobjects>
  <specobjects doctype="feat">
    <specobject><id>f</id><version>1</version>
      <coverage><needscoverage><needsobj>req</needsobj></needscoverage></coverage>
    </specobject>
  </specobjects>
</specdocument>`)

	trace, err := readASpecReport(report, domain.OFTSuccess)

	if err != nil {
		t.Fatal(err)
	}
	want := domain.NewOFTTrace([]domain.OFTSpecificationItem{
		domain.NewOFTSpecificationItem(parseID(t, "req~r~1"), domain.NewOFTCoverage(nil, []domain.OFTSpecificationItemID{parseID(t, "feat~f~1")})),
		domain.NewOFTSpecificationItem(parseID(t, "feat~f~1"), domain.NewOFTCoverage([]domain.OFTArtifactType{"req"}, []domain.OFTSpecificationItemID{})),
	}, domain.OFTSuccess)
	if !reflect.DeepEqual(trace, want) {
		t.Fatalf("got %+v, want %+v", trace, want)
	}
}

func TestReadASpecReportRefusesAMalformedCoverageLink(t *testing.T) {
	report := `<specdocument><specobjects doctype="req"><specobject><id>r</id><version>1</version>
<covering><coveredType><id>f</id><version>1</version><doctype></doctype></coveredType></covering>
</specobject></specobjects></specdocument>`

	if _, err := readASpecReport([]byte(report), domain.OFTSuccess); err == nil || !strings.Contains(err.Error(), "coverage link of req~r~1") {
		t.Fatalf("expected an error naming the coverage link, got %v", err)
	}
}

func TestReadASpecReportRefusesAnotherRootElement(t *testing.T) {
	for _, output := range []string{
		`<html><body>OpenFastTrace</body></html>`,
		`<specobjects doctype="feat"><specobject><id>f</id><version>1</version></specobject></specobjects>`,
	} {
		if _, err := readASpecReport([]byte(output), domain.OFTSuccess); err == nil || !strings.Contains(err.Error(), "specdocument") {
			t.Errorf("expected an error naming specdocument for %q, got %v", output, err)
		}
	}
}

func TestReadASpecReportReportsAnUnreadableReport(t *testing.T) {
	if _, err := readASpecReport([]byte("Exception in thread main"), domain.OFTFailure); err == nil {
		t.Fatal("expected an error for output that is no ASpec report")
	}
}
