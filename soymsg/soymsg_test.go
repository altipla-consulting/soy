package soymsg

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/altipla-consulting/soy/ast"
	"github.com/altipla-consulting/soy/parse"
)

// Test that NewMessage correctly splits a message string into message parts.
func TestNewMessage(t *testing.T) {
	type test struct {
		input  string
		output []Part
	}
	var txt = func(str string) Part { return RawTextPart{str} }
	var ph = func(name string) Part { return PlaceholderPart{name} }
	var tests = []test{
		{"", nil},
		{"hello world", []Part{txt("hello world")}},
		{"hello {WORLD}", []Part{txt("hello "), ph("WORLD")}},
		{"{HELLO_WORLD}", []Part{ph("HELLO_WORLD")}},
		{"{A_1}{A_2}", []Part{ph("A_1"), ph("A_2")}},
		{"{}", []Part{txt("{}")}},
		{"{ }", []Part{txt("{ }")}},
		{"{br}", []Part{txt("{br}")}},
		{"x{A}{B} {C}.", []Part{txt("x"), ph("A"), ph("B"), txt(" "), ph("C"), txt(".")}},
	}

	for _, test := range tests {
		var msg = NewMessage(0, test.input)
		if !reflect.DeepEqual(msg.Parts, test.output) {
			t.Errorf("(actual) %v != %v (expected)", msg.Parts, test.output)
		}
	}
}

// Tests that the set of messages (ids and placeholders) extracted from
// features.soy is the same as that generated by the official java
// implementation.
func TestFeatureExtractedMsgs(t *testing.T) {
	type test struct {
		msg   *ast.MsgNode
		id    uint64
		phstr string
	}

	// test data taken from closure-templates/examples/examples_extracted.xlf
	var tests = []test{
		// Simple messages
		{msg("noun", "The word 'Archive' used as a noun, i.e. an information store.", "Archive"),
			7224011416745566687, "Archive"},
		{msg("verb", "The word 'Archive' used as a verb, i.e. to store information.", "Archive"),
			4826315192146469447, "Archive"},
		{msg("", "", "A trip was taken."),
			3329840836245051515, "A trip was taken."},
		{msg("", "Ask user to pick best keyword", "Your favorite keyword"),
			2209690285855487595, "Your favorite keyword"},
		{msg("", "Link to Help", "Help"),
			7911416166208830577, "Help"},

		// Messages with dataref placeholders
		{msg("", "Example: Alice took a trip to wonderland.", "{$name} took a trip to {$destination}."),
			768490705511913603, "{NAME} took a trip to {DESTINATION}."},
		{msg("", "Example: 5 is nowhere near the value of pi.", "{$pi} is nowhere near the value of pi."),
			889614911019327165, "{PI} is nowhere near the value of pi."},
		{msg("", "Example: Alice took a trip.", "{$name} took a trip."),
			3179387603303514412, "{NAME} took a trip."},

		// Messages with html tags
		// {msg("", "Link to the unreleased 'Labs' feature.", `Click <a href="{$labsUrl}">here</a> to access Labs.`),
		// 	5539341884085868292, `Click {START_LINK}here{END_LINK} to access Labs.`},

		// Messages with calls
		{msg("", "Example: The set of prime numbers is {2, 3, 5, 7, 11, 13, ...}.", `
The set of {$setName} is {lb}
{call .buildCommaSeparatedList_}
  {param items: $setMembers /}
{/call}
, ...{rb}.`),
			135956960462609535, "The set of {SET_NAME} is {{XXX}, ...}."},

		// Plural

		// TODO: Clarify with closure-templates mailing list whether ids should be
		// calculated with braced PHs or not.  Presently we do not use braced phs for id.
		{msg("", "The number of eggs you need.", `
{plural $eggs}
  {case 1}You have one egg
  {default}You have {$eggs} eggs
{/plural}`),
			176798647517908084, "{EGGS_1,plural,=1{You have one egg}other{You have {EGGS_2} eggs}}"},
		// would be 8336954131281929964 without/ braced phs

		// TODO: Add test that needs placeholder index
		// TODO: Test equivalent nodes
	}

	for _, test := range tests {
		SetPlaceholdersAndID(test.msg)
		if test.id != test.msg.ID {
			t.Errorf("(actual) %v != %v (expected)", test.msg.ID, test.id)
		}

		var actual = PlaceholderString(test.msg)
		if test.phstr != actual {
			t.Errorf("(actual) %v != %v (expected)", actual, test.phstr)
		}
	}
}

func msg(meaning, desc string, body string) *ast.MsgNode {
	var msgtmpl = fmt.Sprintf("{msg meaning=%q desc=%q}%s{/msg}", meaning, desc, body)
	var sf, err = parse.SoyFile("", msgtmpl)
	if err != nil {
		panic(err)
	}
	return sf.Body[0].(*ast.MsgNode)
}

func txt(str string) *ast.RawTextNode {
	return &ast.RawTextNode{0, []byte(str)}
}
