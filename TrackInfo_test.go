package radiowatch

import "testing"

var (
	tests = []struct {
		input    string
		expected string
	}{
		{"already_normalized", "already_normalized"},
		{"NoSpecialChars", "nospecialchars"},
		{"Some White Space", "some_white_space"},
		{"Moar!?Oh, yeah!<.<", "moar_oh_yeah"},
	}
)

func TestTrackInfo_NormalizedStationName(t *testing.T) {

	for _, tl := range tests {
		ti := TrackInfo{Station:tl.input}
		actual := ti.NormalizedStationName()
		if actual != tl.expected {
			t.Errorf("NormalizedStationName(%v): Expected %v, got %v", tl.input, tl.expected, actual)
		}
	}
}

func benchNormalize(input string, b *testing.B) {
	ti := TrackInfo{Station:input}
	for i := 0; i < b.N; i++ {
		ti.NormalizedStationName()
	}
}

func BenchmarkTI_NormalizedAlready(b *testing.B) {benchNormalize(tests[0].input, b)}
func BenchmarkTI_NormalizedUpper(b *testing.B) {benchNormalize(tests[1].input, b)}
func BenchmarkTI_NormalizedWhitespace(b *testing.B) {benchNormalize(tests[2].input, b)}
func BenchmarkTI_NormalizedSpecialchars(b *testing.B) {benchNormalize(tests[3].input, b)}
