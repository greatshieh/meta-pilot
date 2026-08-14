package response

import "server/pkg/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
