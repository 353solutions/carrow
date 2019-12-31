package csv

import (
	"os"
	"testing"
)

func TestRead(t *testing.T) {
	file, err := os.Open("cart.csv")
	if err != nil {
		t.Fatal(err)
	}

	table, err := Read(file)
	if err != nil {
		t.Fatal(err)
	}

	if table.NumRows() != 4 {
		t.Fatal(table.NumRows())
	}

}
