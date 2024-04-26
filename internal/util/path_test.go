package util_test

import (
	"booking/internal/util"
	"strings"
	"testing"
)

func TestGetAbsoluteProjectPath(t *testing.T) {
	t.Run("Function should return valid path", func(t *testing.T) {
		util := util.New()
		got := util.GetAbsoluteProjectPath()
		rootFolderName := "booking"
		if !(strings.HasSuffix(got, rootFolderName+"/") || strings.HasSuffix(got, rootFolderName+"\\")) {
			t.Errorf(`the path should end with %v/ or %v\ but got %v.`, rootFolderName, rootFolderName, got)
		}
	})
}
