package differential

import (
	"testing"

	"elle-go/fixtures"
)

func TestToEDNWriteSkew(t *testing.T) {
	want := `{:process 1 :type :invoke :value [[:r "y" nil] [:append "x" 1]]}
{:process 1 :type :ok :value [[:r "y" nil] [:append "x" 1]]}
{:process 2 :type :invoke :value [[:r "x" nil] [:append "y" 1]]}
{:process 2 :type :ok :value [[:r "x" nil] [:append "y" 1]]}
`
	if got := ToEDN(fixtures.WriteSkew()); got != want {
		t.Errorf("ToEDN mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}

func TestToEDNSerializable(t *testing.T) {
	want := `{:process 1 :type :invoke :value [[:append "x" 1]]}
{:process 1 :type :ok :value [[:append "x" 1]]}
{:process 2 :type :invoke :value [[:r "x" nil] [:append "x" 2]]}
{:process 2 :type :ok :value [[:r "x" [1]] [:append "x" 2]]}
{:process 3 :type :invoke :value [[:r "x" nil] [:r "y" nil]]}
{:process 3 :type :ok :value [[:r "x" [1 2]] [:r "y" nil]]}
`
	if got := ToEDN(fixtures.Serializable()); got != want {
		t.Errorf("ToEDN mismatch:\n got:\n%s\nwant:\n%s", got, want)
	}
}
