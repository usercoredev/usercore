package database

import (
	"github.com/usercoredev/usercore/pagination"
	"testing"
)

func TestGetSortableFields(t *testing.T) {
	type sortableStruct struct {
		UnsortableField string
		SortableField   string `sortable:"true"`
	}

	fields := getSortableFields(sortableStruct{})
	if len(fields) != 1 {
		t.Errorf("Expected 1 sortable field, got %d", len(fields))
	}
	if !fields["SortableField"] {
		t.Errorf("Expected sortable_field to be true, got false")
	}

	if fields["UnsortableField"] {
		t.Errorf("Expected unsortable_field to be false, got true")
	}
}

func TestConvertToOrder(t *testing.T) {
	user := User{}
	pm := &pagination.PageMetadata{OrderBy: "1; DROP TABLE users", Order: ""}
	orderClause := user.ConvertToOrder(*pm)
	if orderClause != "created_at desc" {
		t.Errorf("Expected default order clause for unsafe OrderBy, got %s", orderClause)
	}
}
