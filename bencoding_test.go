package bazbittorrent_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	bazbittorrent "github.com/CalvoM/baz_bit-torrent"
)

func TestSimple(t *testing.T) {
	actual := 2 + 2
	expected := 4
	if actual != expected {
		t.Errorf("got %q, wanted %q", actual, expected)
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{"4:spam", "spam"}, {"5:oceans", "ocean"}, {"6:eleven", "eleven"}, {"1:erft", "e"}}
	for _, test := range tests {
		bencoding := bazbittorrent.NewDecoder(strings.NewReader(test.input))
		actual, _ := bencoding.DecodeString()
		if actual != test.expected {
			t.Errorf("got %q, wanted %q", actual, test.expected)
		}
	}
}

func TestDecodeInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{{"i90e", 90}, {"i5e", 5}, {"i1000e", 1000}, {"i9e", 9}}
	for _, test := range tests {
		bencoding := bazbittorrent.NewDecoder(strings.NewReader(test.input))
		actual, _ := bencoding.DecodeInt()
		if actual != test.expected {
			t.Errorf("got %q, wanted %q", actual, test.expected)
		}
	}
}

func TestDecodeList(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{{"li90ei1ei12345eei8e", []any{90, 1, 12345}}, {"li90e4:spame", []any{90, "spam"}}}
	for _, test := range tests {
		bencoding := bazbittorrent.NewDecoder(strings.NewReader(test.input))
		actual, _ := bencoding.DecodeList()
		if !slices.Equal(actual, test.expected) {
			t.Errorf("got %q, wanted %q", actual, test.expected)
		}
	}
}

func TestDecodeDict(t *testing.T) {
	tests := []struct {
		input    string
		expected map[any]any
	}{{"d1:a1:be", map[any]any{"a": "b"}}, {"d5:namesl4:sean5:shawni5ee", map[any]any{"names": []any{"sean", "shawn", 5}}}, {"d4:infod4:test5:valueee", map[any]any{"info": map[any]any{"test": "value"}}}}
	for _, test := range tests {
		bencoding := bazbittorrent.NewDecoder(strings.NewReader(test.input))
		actual, _ := bencoding.DecodeDict()
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("got %q, wanted %q", actual, test.expected)
		}
	}
}
