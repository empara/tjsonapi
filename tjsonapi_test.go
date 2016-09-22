package tjsonapi

import (
	"os"
	"reflect"
	"testing"
)

type TestStruct struct {
	ID                int    `jsonapi:"identifier,test"`
	FirstAttr         int    `jsonapi:"attribute,first"`
	SecondAttr        string `jsonapi:"attribute,second"`
	OneRelationship   int    `jsonapi:"relationship,one,data,other"`
	ManyRelationships []int  `jsonapi:"relationship,many,data,other"`
}

var basicTestStruct = TestStruct{
	ID:                42,
	FirstAttr:         84,
	SecondAttr:        "a string",
	OneRelationship:   4242,
	ManyRelationships: []int{21, 42},
}

var basicExpectedRoot *Root

func TestMain(m *testing.M) {
	r := NewResource()
	r.ID = "42"
	r.Type = "test"
	r.Attributes.AddAttribute("first", 84)
	r.Attributes.AddAttribute("second", "a string")
	r.Relationships["one"] = NewRelationship()
	r.Relationships["one"].Data = NewResourceLinkageToOne()
	r.Relationships["one"].Data.SetResourceIdentifier(&ResourceIdentifier{
		ID:   "4242",
		Type: "other",
		Meta: NewMeta(),
	})
	r.Relationships["many"] = NewRelationship()
	r.Relationships["many"].Data = NewResourceLinkageToMany()
	r.Relationships["many"].Data.AddResourceIdentifier(&ResourceIdentifier{
		ID:   "21",
		Type: "other",
		Meta: NewMeta(),
	})
	r.Relationships["many"].Data.AddResourceIdentifier(&ResourceIdentifier{
		ID:   "42",
		Type: "other",
		Meta: NewMeta(),
	})

	basicExpectedRoot = NewRoot()
	basicExpectedRoot.Data = NewResourcesOne()
	basicExpectedRoot.Data.SetResource(r)

	os.Exit(m.Run())
}

func resourceEqual(first *Resource, second *Resource) bool {
	if first.ID != second.ID ||
		first.Type != second.Type ||
		!reflect.DeepEqual(first.Attributes, second.Attributes) ||
		!reflect.DeepEqual(first.Relationships, second.Relationships) {
		return false
	}
	return true
}

func TestEncode(t *testing.T) {
	basicRoot, err := Marshal(basicTestStruct)
	if err != nil {
		t.Error("Error while marshaling root")
	}
	if !resourceEqual(basicRoot.Data.Data[0], basicExpectedRoot.Data.Data[0]) {
		t.Error("Encoded root does not match expected")
	}
}

func TestDecode(t *testing.T) {
	basicRoot, _ := Marshal(basicTestStruct)

	var reencodedStruct TestStruct
	err := Unmarshal(basicRoot, &reencodedStruct)
	if err != nil {
		t.Error("Error while unmarshaling re-encoded root")
	}
	if !reflect.DeepEqual(basicTestStruct, reencodedStruct) {
		t.Error("Re-encoded root does not match expected")
	}

	redecodedRoot, err := Marshal(reencodedStruct)
	if err != nil {
		t.Error("Error while marshaling re-encoded root")
	}
	if !resourceEqual(basicRoot.Data.Data[0], redecodedRoot.Data.Data[0]) {
		t.Error("Re-encoded root does not match encoded root")
	}
}
