package knowledge

import (
	"testing"
)

const (
	READ_UNTIL1 = " papera { pluto */ ciao"
	READ_UNTIL2 = "*/"
	SNIPPET     = `/* ciao a tutti */
	A = 1 * (B + C); /* good */
	Z=min(2,3,4,5);
alfa = 
	beta || gamma;
`
	BLOCK1 = ` { /*
		ciao */
alfa = 
	beta || gamma;
} ciao 
pippo
`
	BLOCK1_snippet = "/*ciao */alfa = beta || gamma;"
	BLOCK1_source  = "ciao pippo"

	BLOCK2         = "{ ciao { pippo } pluto }"
	BLOCK2_snippet = "ciao { pippo } pluto"
	BLOCK2_source  = ""
)

var EXPRS = []string{
	"A = 1 * (B + C)",
	"Z=min(2,3,4,5)",
	"alfa = beta || gamma",
}

func TestReadUntil(t *testing.T) {
	token, source, err := readUntil(READ_UNTIL1, "{")
	if err != nil {
		t.Error(err)
	}
	if token != "papera" {
		t.Errorf("1 - wrong token [%v]", token)
	}
	if source != "pluto */ ciao" {
		t.Errorf("1 - wrong source [%v]", source)
	}
	token, source, err = readUntil(READ_UNTIL1, "*/")
	if err != nil {
		t.Error(err)
	}
	if token != "papera { pluto" {
		t.Errorf("2 - wrong token [%v]", token)
	}
	if source != "ciao" {
		t.Errorf("2 - wrong source [%v]", source)
	}
	token, source, err = readUntil(READ_UNTIL2, "*/")
	if err != nil {
		t.Error(err)
	}
	if token != "" {
		t.Errorf("3 - wrong token %v", token)
	}
	if source != "" {
		t.Errorf("3 - wrong source %v", source)
	}
}

func TestGetEspressions(t *testing.T) {
	exprs, err := getExpressions(SNIPPET)
	if err != nil {
		t.Error(err)
	}
	if len(exprs) != 3 {
		t.Fatal("expected 3 expressions")
	}
	for i := range EXPRS {
		if exprs[i] != EXPRS[i] {
			t.Fatalf("expected [%v] not [%v]", EXPRS[i], exprs[i])
		}
	}
}

func TestGetBlock(t *testing.T) {
	snippet, source, err := getBlock(BLOCK1)
	if err != nil {
		t.Fatal(err)
	}
	if snippet != BLOCK1_snippet {
		t.Errorf("snippet - expected [%v] not [%v]", BLOCK1_snippet, snippet)
	}
	if source != BLOCK1_source {
		t.Errorf("source - expected [%v] not [%v]", BLOCK1_source, source)
	}
	snippet, source, err = getBlock(BLOCK2)
	if err != nil {
		t.Fatal(err)
	}
	if snippet != BLOCK2_snippet {
		t.Errorf("snippet - expected [%v] not [%v]", BLOCK2_snippet, snippet)
	}
	if source != BLOCK2_source {
		t.Errorf("source - expected [%v] not [%v]", BLOCK2_source, source)
	}
}

func TestParseEventFunction(t *testing.T) {
	name, weight, ok, err := parseEventFunction("event(test1,1.0)")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("that was a event function!")
	}
	if name != "test1" {
		t.Fatalf("unexpected [%v]", name)
	}
	if weight != 1.0 {
		t.Fatalf("unexpected [%v]", weight)
	}
}
