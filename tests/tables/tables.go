package tables

import "github.com/barmi/go-admin/plugins/admin/modules/table"

var Generators = map[string]table.Generator{
	"posts":    GetPostsTable,
	"authors":  GetAuthorsTable,
	"external": GetExternalTable,
}
