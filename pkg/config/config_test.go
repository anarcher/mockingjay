package config

import (
	"testing"
)

func TestReadConfig(t *testing.T) {
	yamlStr := `
proxies:
    - target: http://localhost:8081
      filter:
          form:
              action: GetStats
              namespace: aaa
`

	cfg, err := ReadConfigBytes([]byte(yamlStr))
	if err != nil {
		t.Error(err)
	}

	if len(cfg.Proxies) <= 0 {
		t.Error("len(cfg.Proxies) <= 0")
	}

	if cfg.Proxies[0].Filter.Form["action"] != "GetStats" {
		t.Error("cfg.filter.param.action != GetStats")
	}
}
