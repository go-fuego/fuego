package fuego

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Constructor returns new MetadataParsers instance with empty registeredParsers slice
func TestNewMetadataParsersReturnsEmptyInstance(t *testing.T) {
	parsers := NewMetadataParsers()

	if parsers == nil {
		t.Error("Expected non-nil MetadataParsers instance")
	}

	if len(parsers.registeredParsers) != 0 {
		t.Errorf("Expected empty registeredParsers, got %d items", len(parsers.registeredParsers))
	}

	if len(parsers.registeredNames) != 0 {
		t.Errorf("Expected empty registeredNames map, got %d items", len(parsers.registeredNames))
	}
}

// Constructor initializes registeredNames as empty map with bool values
func TestNewMetadataParsersInitializesRegisteredNamesAsEmptyMap(t *testing.T) {
	parsers := NewMetadataParsers()

	if parsers == nil {
		t.Error("Expected non-nil MetadataParsers instance")
	}

	if len(parsers.registeredNames) != 0 {
		t.Errorf("Expected empty registeredNames map, got %d items", len(parsers.registeredNames))
	}
}

// Zero value initialization without using constructor
func TestMetadataParsersZeroValueInitialization(t *testing.T) {
	var parsers MetadataParsers

	if parsers.registeredParsers != nil {
		t.Error("Expected nil registeredParsers slice")
	}

	if parsers.registeredNames != nil {
		t.Error("Expected nil registeredNames map")
	}
}

// Return all registered metadata parsers from the internal slice
func TestGetRegisteredParsersReturnsAllParsers(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: func(params MetadataParserParams) { return }},
			{Name: "parser2", Parser: func(params MetadataParserParams) { return }},
		},
	}

	parsers := mp.GetRegisteredParsers()

	if len(parsers) != 2 {
		t.Errorf("Expected 2 parsers, got %d", len(parsers))
	}

	if parsers[0].Name != "parser1" || parsers[1].Name != "parser2" {
		t.Errorf("Unexpected parser names")
	}
}

// Return empty slice when no parsers are registered
func TestGetRegisteredParsersReturnsEmptySlice(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{},
	}

	parsers := mp.GetRegisteredParsers()

	if len(parsers) != 0 {
		t.Errorf("Expected empty slice, got slice with length %d", len(parsers))
	}
}

// Verify concurrent access is properly synchronized using mutex lock
func TestGetRegisteredParsersConcurrentAccess(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: func(params MetadataParserParams) { return }},
			{Name: "parser2", Parser: func(params MetadataParserParams) { return }},
		},
	}

	var wg sync.WaitGroup
	numGoroutines := 10
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			parsers := mp.GetRegisteredParsers()
			if len(parsers) != 2 {
				t.Errorf("Expected 2 parsers, got %d", len(parsers))
			}
		}()
	}

	wg.Wait()
}

// Check that returned slice contains expected parser entries
func TestGetRegisteredParsersContainsExpectedEntries(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "expectedParser1", Parser: func(params MetadataParserParams) { return }},
			{Name: "expectedParser2", Parser: func(params MetadataParserParams) { return }},
		},
	}

	parsers := mp.GetRegisteredParsers()

	if len(parsers) != 2 {
		t.Errorf("Expected 2 parsers, got %d", len(parsers))
	}

	if parsers[0].Name != "expectedParser1" || parsers[1].Name != "expectedParser2" {
		t.Errorf("Unexpected parser names")
	}
}

// Verify returned slice matches original registered parsers
func TestGetRegisteredParsersMatchesOriginal(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: func(params MetadataParserParams) { return }},
			{Name: "parser2", Parser: func(params MetadataParserParams) { return }},
		},
	}

	parsers := mp.GetRegisteredParsers()

	if !reflect.DeepEqual(parsers, mp.registeredParsers) {
		t.Errorf("Returned parsers do not match the original registered parsers")
	}
}

// Reset clears all registered parsers from the slice
func TestResetClearsRegisteredParsers(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: nil},
			{Name: "parser2", Parser: nil},
		},
		registeredNames: map[string]bool{
			"parser1": true,
			"parser2": true,
		},
	}

	mp.Reset()

	if len(mp.registeredParsers) != 0 {
		t.Errorf("Expected registeredParsers to be empty, got %d items", len(mp.registeredParsers))
	}

	if len(mp.registeredNames) != 0 {
		t.Errorf("Expected registeredNames to be empty, got %d items", len(mp.registeredNames))
	}
}

// Reset called on nil MetadataParsers struct
func TestResetOnNilMetadataParsers(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when calling Reset on nil MetadataParsers")
		}
	}()

	var mp *MetadataParsers
	mp.Reset()
}

// Reset reinitializes the registeredNames map as empty
func TestResetClearsRegisteredNames(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: nil},
			{Name: "parser2", Parser: nil},
		},
		registeredNames: map[string]bool{
			"parser1": true,
			"parser2": true,
		},
	}

	mp.Reset()

	if len(mp.registeredNames) != 0 {
		t.Errorf("Expected registeredNames to be empty, got %d items", len(mp.registeredNames))
	}
}

// Reset called when parsers list is already empty
func TestResetWhenParsersListIsEmpty(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{},
		registeredNames:   map[string]bool{},
	}

	mp.Reset()

	if len(mp.registeredParsers) != 0 {
		t.Errorf("Expected registeredParsers to be empty, got %d items", len(mp.registeredParsers))
	}

	if len(mp.registeredNames) != 0 {
		t.Errorf("Expected registeredNames to be empty, got %d items", len(mp.registeredNames))
	}
}

// State remains consistent between parsers list and names map
func TestResetMaintainsConsistencyBetweenParsersAndNames(t *testing.T) {
	mp := &MetadataParsers{
		registeredParsers: []MetadataParserEntry{
			{Name: "parser1", Parser: nil},
			{Name: "parser2", Parser: nil},
		},
		registeredNames: map[string]bool{
			"parser1": true,
			"parser2": true,
		},
	}

	mp.Reset()

	if len(mp.registeredParsers) != 0 {
		t.Errorf("Expected registeredParsers to be empty, got %d items", len(mp.registeredParsers))
	}

	if len(mp.registeredNames) != 0 {
		t.Errorf("Expected registeredNames to be empty, got %d items", len(mp.registeredNames))
	}
}

// Register new parser at start position successfully
func TestRegisterMetadataParserAtStart(t *testing.T) {
	parserName := "testParser"

	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{})

	err := mp.RegisterMetadataParser(parserName, func(params MetadataParserParams) { return }, "prepend", "")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	registeredParsers := mp.GetRegisteredParsers()
	if len(registeredParsers) == 0 {
		t.Error("Expected parser to be registered")
	}

	if registeredParsers[0].Name != parserName {
		t.Errorf("Expected first parser to be %s, got %s", parserName, registeredParsers[0].Name)
	}
}

// Invalid position parameter returns appropriate error
func TestRegisterMetadataParserInvalidPosition(t *testing.T) {
	parserName := "testParser"

	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{})

	err := mp.RegisterMetadataParser(parserName, func(params MetadataParserParams) {
		return
	}, "invalid", "")

	if err == nil {
		t.Error("Expected error for invalid position, got nil")
	}

	expectedErr := "Invalid position. Use 'prepend', 'append', 'before', or 'after'"
	if err.Error() != expectedErr {
		t.Errorf("Expected error message '%s', got '%s'", expectedErr, err.Error())
	}
}

// Register parser at end position with unique name
func TestRegisterMetadataParserAtEndPosition(t *testing.T) {
	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{})

	mockParser := func(params MetadataParserParams) { return }
	err := mp.RegisterMetadataParser("unique-parser", mockParser, "append", "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mp.registeredParsers))
	assert.Equal(t, "unique-parser", mp.registeredParsers[0].Name)
}

// Register parser before existing parser with valid relativeTo
func TestRegisterMetadataParserBeforeExistingParser(t *testing.T) {
	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{
		{Name: "existing-parser", Parser: func(params MetadataParserParams) { return }},
	})

	mockParser := func(params MetadataParserParams) { return }
	err := mp.RegisterMetadataParser("new-parser", mockParser, "before", "existing-parser")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(mp.registeredParsers))
	assert.Equal(t, "new-parser", mp.registeredParsers[0].Name)
	assert.Equal(t, "existing-parser", mp.registeredParsers[1].Name)
}

// Register parser after existing parser with valid relativeTo
func TestRegisterMetadataParserAfterRelativePosition(t *testing.T) {
	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{
		{Name: "existing-parser", Parser: func(params MetadataParserParams) { return }},
	})

	mockParser := func(params MetadataParserParams) { return }
	err := mp.RegisterMetadataParser("new-parser", mockParser, "after", "existing-parser")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(mp.registeredParsers))
	assert.Equal(t, "existing-parser", mp.registeredParsers[0].Name)
	assert.Equal(t, "new-parser", mp.registeredParsers[1].Name)
}

// Register duplicate parser name returns nil without modifying list
func TestRegisterDuplicateMetadataParserName(t *testing.T) {
	mp := NewMetadataParsers()
	mp.InitializeMetadataParsers([]MetadataParserEntry{
		{Name: "existing-parser", Parser: func(params MetadataParserParams) { return }},
	})

	mockParser := func(params MetadataParserParams) { return }
	err := mp.RegisterMetadataParser("existing-parser", mockParser, "prepend", "")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(mp.registeredParsers))
	assert.Equal(t, "existing-parser", mp.registeredParsers[0].Name)
}
