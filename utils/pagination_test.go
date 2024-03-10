package utils

import (
	"testing"
)

func TestSetTotalCount(t *testing.T) {
	tests := []struct {
		totalCount         int32
		pageSize           int32
		expectedTotalPages int32
	}{
		{totalCount: 100, pageSize: 10, expectedTotalPages: 10},
		{totalCount: 101, pageSize: 10, expectedTotalPages: 11},
		{totalCount: 0, pageSize: 10, expectedTotalPages: 0},
		{totalCount: 100, pageSize: 0, expectedTotalPages: 0},
		{totalCount: -10, pageSize: 10, expectedTotalPages: 0},
		{totalCount: 100, pageSize: -10, expectedTotalPages: 0},
		{totalCount: 0, pageSize: 0, expectedTotalPages: 0},
		{totalCount: 0, pageSize: 1, expectedTotalPages: 0},
		{totalCount: 1, pageSize: 1, expectedTotalPages: 1},
		{totalCount: 1, pageSize: 0, expectedTotalPages: 0},
	}

	for _, test := range tests {
		pm := PageMetadata{PageSize: test.pageSize}
		pm.SetTotalCount(test.totalCount)
		if pm.TotalPages != test.expectedTotalPages {
			t.Errorf("SetTotalCount(%d): expected total pages %d, got %d", test.totalCount, test.expectedTotalPages, pm.TotalPages)
		}
		if pm.TotalCount != test.totalCount && test.totalCount > 0 {
			t.Errorf("SetTotalCount(%d): expected total count %d, got %d", test.totalCount, test.totalCount, pm.TotalCount)
		}
	}
}

func TestPageMetadata_SetPage_BeforeSetTotalCount(t *testing.T) {
	pm := &PageMetadata{PageSize: 10}
	pm.SetPage(5)
	if pm.Page != 1 {
		t.Errorf("Expected Page to default to 1 when TotalPages is 0, got %d", pm.Page)
	}
}

func TestPageMetadata_Offset_ZeroValues(t *testing.T) {
	pm := &PageMetadata{PageSize: 0, Page: 0}
	offset := pm.Offset()
	if offset != 0 {
		t.Errorf("Expected offset to be 0 when PageSize and Page are 0, got %d", offset)
	}
}
