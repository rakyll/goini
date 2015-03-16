package ini

import (
	"io/ioutil"
	"testing"
	"strings"
)

const (
	exampleStr = `key1 = true

[section1]
key1 = value2
key2 = 5
key3 = 1.3

[section2]
key1 = 5

`
)

var (
	dict Dict
	err  error
)

func init() {
	dict, err = Load("example.ini")
}

func TestLoad(t *testing.T) {
	if err != nil {
		t.Fatal("Example: load error:", err)
	}
	t.Logf("dict:%v:\n", dict)

	_, err2 := Load("badExample.ini")
	if err2 == nil {
		t.Errorf("BadExample: loaded without error:")
	}
	t.Logf("badExample:err:%v:\n", err2)

}

func TestLoadString(t *testing.T) {

	exampleDict, err := LoadString(exampleStr)
	if err != nil {
		t.Errorf("LoadString: failed to load exampleStr:%v:", err)
	}

	nsk1, _ := exampleDict.GetString("section1", "key1")

	if nsk1 != "value2" {
		t.Errorf("Dict not loaded from string as expected.")
	}


	nk1, _ := exampleDict.GetBool("", "key1")

	if nk1 != true {
		t.Errorf("Dict not loaded from string as expected.")
	}

}
func TestNoSemiColon(t *testing.T) {
	b, found := dict.GetDouble("wine", "zup")
	if !found  {
		t.Error("Example: failed to find key for line with no semi-colon.")
	}
	if b != 12.5 {
		t.Error("Example: failed to find value 12.5 for line with no semi-colon.")
	}
}

func TestOnlyKey(t *testing.T) {
	_, found := dict.GetString("wine", "nuch")
	if !found  {
		t.Error("Example: failed to find key for line with no value.")
	}
}

func TestWrite(t *testing.T) {
	d, err := Load("empty.ini")
	if err != nil {
		t.Error("Example: load error:", err)
	}
	d.SetString("", "key", "value")
	tempFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Error("Write: Couldn't create temp file.", err)
	}
	err = Write(tempFile.Name(), &d)
	if err != nil {
		t.Error("Write: Couldn't write to temp config file.", err)
	}
	contents, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Error("Write: Couldn't read from the temp config file.", err)
	}
	if string(contents) != "key = value\n\n" {
		t.Error("Write: Contents of the config file doesn't match the expected.")
	}
}

func TestGetBool(t *testing.T) {
	b, found := dict.GetBool("pizza", "ham")
	if !found || !b {
		t.Error("Example: parse error for key ham of section pizza.")
	}
	b, found = dict.GetBool("pizza", "mushrooms")
	if !found || !b {
		t.Error("Example: parse error for key mushrooms of section pizza.")
	}
	b, found = dict.GetBool("pizza", "capres")
	if !found || b {
		t.Error("Example: parse error for key capres of section pizza.")
	}
	b, found = dict.GetBool("pizza", "cheese")
	if !found || b {
		t.Error("Example: parse error for key cheese of section pizza.")
	}
}

func TestGetStringIntAndDouble(t *testing.T) {
	str, found := dict.GetString("wine", "grape")
	if !found || str != "Cabernet Sauvignon" {
		t.Error("Example: parse error for key grape of section wine.")
	}
	i, found := dict.GetInt("wine", "year")
	if !found || i != 1989 {
		t.Error("Example: parse error for key year of section wine.")
	}
	str, found = dict.GetString("wine", "country")
	if !found || str != "Spain" {
		t.Error("Example: parse error for key grape of section wine.")
	}
	d, found := dict.GetDouble("wine", "alcohol")
	if !found || d != 12.5 {
		t.Error("Example: parse error for key grape of section wine.")
	}
}

func TestSetBoolAndStringAndIntAndDouble(t *testing.T) {
	dict.SetBool("pizza", "ham", false)
	b, found := dict.GetBool("pizza", "ham")
	if !found || b {
		t.Error("Example: bool set error for key ham of section pizza.")
	}
	dict.SetString("pizza", "ham", "no")
	n, found := dict.GetString("pizza", "ham")
	if !found || n != "no" {
		t.Error("Example: string set error for key ham of section pizza.")
	}
	dict.SetInt("wine", "year", 1978)
	i, found := dict.GetInt("wine", "year")
	if !found || i != 1978 {
		t.Error("Example: int set error for key year of section wine.")
	}
	dict.SetDouble("wine", "not-exists", 5.6)
	d, found := dict.GetDouble("wine", "not-exists")
	if !found || d != 5.6 {
		t.Error("Example: float set error for not existing key for wine.")
	}
}

func TestDelete(t *testing.T) {
	d, err := Load("empty.ini")
	if err != nil {
		t.Error("Example: load error:", err)
	}
	d.SetString("pizza", "ham", "yes")
	d.Delete("pizza", "ham")
	_, found := d.GetString("pizza", "ham")
	if found {
		t.Error("Example: delete error for key ham of section pizza.")
	}
	if len(d.GetSections()) > 1 {
		t.Error("Only a single section should exist after deletion.")
	}
}

func TestGetNotExist(t *testing.T) {
	_, found := dict.GetString("not", "exist")
	if found {
		t.Error("There is no key exist of section not.")
	}
}

func TestGetSections(t *testing.T) {
	sections := dict.GetSections()
	if len(sections) != 3 {
		t.Error("The number of sections is wrong:", len(sections))
	}
	for _, section := range sections {
		if section != "" && section != "pizza" && section != "wine" {
			t.Errorf("Section '%s' should not be exist.", section)
		}
	}
}

func TestString(t *testing.T) {
	d, err := Load("empty.ini")
	if err != nil {
		t.Error("Example: load error:", err)
	}
	d.SetBool("", "key1", true)
	d.SetString("section1", "key1", "value2")
	d.SetInt("section1", "key2", 5)
	d.SetDouble("section1", "key3", 1.3)
	d.SetDouble("section2", "key1", 5.0)

	s2k1, _ := d.GetDouble("section2", "key1")
	if s2k1 != 5.0 {
		t.Errorf("GetDouble: failed to load section2 key1:d:%v:", d)
	}

	stringified := d.String()

	if strings.Contains(stringified, "section2") != true {
		t.Fatalf("LoadString: failed to load stringified:section2 missing:%v:", stringified)
	}
	if strings.Contains(stringified, "key1 = true") != true {
		t.Fatalf("LoadString: failed to load stringified:top scope key1 missing:%v:", stringified)
	}
	if strings.Contains(stringified, "key1 = 5") != true {
		t.Fatalf("LoadString: failed to load stringified:section2 key1 is 5 not 5.0:%v:", stringified)
	}
	if strings.Contains(stringified, "key1 = value2") != true {
		t.Fatalf("LoadString: failed to load stringifiedsection1 key1 missing:%v:", stringified)
	}
	newDict, err := LoadString(stringified)
	if err != nil {
		t.Fatalf("LoadString: failed to load stringified:%v:", err)
	}

	a, _ := newDict.GetBool("", "key1")

	if a != true {
		t.Errorf("Dict cannot be stringified as expected:new:%v:old:true:newDict:%v:", a, newDict)
	}

	nsk1, _ := newDict.GetString("section1", "key1")

	if nsk1 != "value2" {
		t.Errorf("Dict cannot be stringified as expected.")
	}

}
