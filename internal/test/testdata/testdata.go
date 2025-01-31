package testdata

import (
	_ "embed"
)

//go:embed data/insert_test_data.sql
var InsertTestData string

//go:embed data/clear_test_data.sql
var ClearTestData string
