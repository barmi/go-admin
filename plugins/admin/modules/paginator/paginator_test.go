package paginator

import (
	"testing"

	_ "github.com/barmi/go-admin-themes/sword"
	"github.com/barmi/go-admin/modules/config"
	"github.com/barmi/go-admin/plugins/admin/modules/parameter"
)

func TestGet(t *testing.T) {
	config.Initialize(&config.Config{Theme: "sword"})
	param := parameter.BaseParam()
	param.Page = "7"
	Get(nil, Config{
		Size:         105,
		Param:        param,
		PageSizeList: []string{"10", "20", "50", "100"},
	})
}
