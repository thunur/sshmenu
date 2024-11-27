package menu

import (
	"github.com/thunur/sshmenu/internal/config"

	"github.com/manifoldco/promptui"
)

type MenuType int

const (
	MenuTypeConnect MenuType = iota // 连接服务器
	MenuTypeDelete                  // 删除服务器
)

func ShowMenu(configs []config.SSHConfig, menuType MenuType) (*config.SSHConfig, error) {
	var templates *promptui.SelectTemplates
	var label string

	switch menuType {
	case MenuTypeConnect:
		templates = &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➤ {{ .Name | cyan }} ({{ .User }}@{{ .Host }})",
			Inactive: "  {{ .Name | white }} ({{ .User }}@{{ .Host }})",
			Selected: "✔ {{ .Name | green }} ({{ .User }}@{{ .Host }})",
		}
		label = "选择要连接的服务器"
	case MenuTypeDelete:
		templates = &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➤ {{ .Name | red }} ({{ .User }}@{{ .Host }})",
			Inactive: "  {{ .Name | white }} ({{ .User }}@{{ .Host }})",
			Selected: "✗ {{ .Name | red }} ({{ .User }}@{{ .Host }})",
		}
		label = "选择要删除的服务器"
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     configs,
		Templates: templates,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}

	return &configs[i], nil
}
