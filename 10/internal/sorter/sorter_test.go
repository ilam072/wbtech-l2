package sorter

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	data := []string{"banana", "Apple", "cherry"}
	expected := []string{"Apple", "banana", "cherry"}

	s := New()
	s.data = data
	s.Sort()

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("Sort() = %v, want %v", s.Data(), expected)
	}
}

func TestSortByColumn(t *testing.T) {
	data := []string{
		"apple\t2",
		"banana\t10",
		"cherry\t1",
	}
	expected := []string{
		"cherry\t1",
		"banana\t10",
		"apple\t2",
	}

	s := New()
	s.data = data

	s.SortByColumn(2, false)

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("SortByColumn() = %v, want %v", s.Data(), expected)
	}
}

func TestSortByColumn_Numeric(t *testing.T) {
	data := []string{
		"apple\t2",
		"banana\t10",
		"cherry\t1",
	}
	expected := []string{
		"cherry\t1",
		"apple\t2",
		"banana\t10",
	}

	s := New()
	s.data = data

	s.SortByColumn(2, true)

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("SortByColumn() = %v, want %v", s.Data(), expected)
	}
}

func TestSortNumeric(t *testing.T) {
	data := []string{"10", "2", "30", "1"}
	expected := []string{"1", "2", "10", "30"}

	s := New()
	s.data = data
	s.SortNumeric()

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("SortNumeric() = %v, want %v", s.Data(), expected)
	}
}

func TestReverse(t *testing.T) {
	data := []string{"a", "b", "c"}
	expected := []string{"c", "b", "a"}

	s := New()
	s.data = data
	s.Reverse()

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("Reverse() = %v, want %v", s.Data(), expected)
	}
}

func TestUnique(t *testing.T) {
	data := []string{"малина", "яблоко", "помидор", "помидор", "малина", "банан", "груша", "огурец"}

	s := New()
	s.data = data
	m := make(map[string]int)
	s.Unique()

	for _, v := range s.data {
		m[v]++
	}

	for _, v := range m {
		if v > 1 {
			t.Errorf("duplicate error")
		}
	}
}

func TestSortByMonth(t *testing.T) {
	data := []string{"Mar", "jan", "Feb", "Dec"}
	expected := []string{"jan", "Feb", "Mar", "Dec"}

	s := New()
	s.data = data
	s.SortByMonth(0)

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("SortByMonth() = %v, want %v", s.Data(), expected)
	}
}

func TestTrimTrailingSpaces(t *testing.T) {
	data := []string{"a   ", "b\t\t", "c"}
	expected := []string{"a", "b", "c"}

	s := New()
	s.data = data
	s.TrimTrailingSpaces()

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("TrimTrailingSpaces() = %v, want %v", s.Data(), expected)
	}
}

func TestSortHumanNumeric(t *testing.T) {
	data := []string{"10K", "500", "5M", "1G"}
	expected := []string{"500", "10K", "5M", "1G"}

	s := New()
	s.data = data
	s.SortHumanNumeric()

	if !reflect.DeepEqual(s.Data(), expected) {
		t.Errorf("SortHumanNumeric() = %v, want %v", s.Data(), expected)
	}
}
