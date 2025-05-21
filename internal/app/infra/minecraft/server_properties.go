package minecraft

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"path"

	"github.com/candbright/mc-dashboard/internal/pkg/utils"
	"github.com/candbright/mc-dashboard/pkg/dw"
	"github.com/magiconair/properties"
)

var (
	serverPropertiesFilter = map[string][]string{
		"gamemode": {"survival", "creative", "adventure"},
	}
)

type ServerProperties struct {
	path    string
	content *template.Template
	*dw.DataWriter[map[string]string]
}

type ServerPropertiesConfig struct {
	RootDir string
}

func NewServerProperties(config ServerPropertiesConfig) (*ServerProperties, error) {
	sp := &ServerProperties{
		path: path.Join(config.RootDir, "server.properties"),
	}
	d, err := dw.New[map[string]string](dw.Config{
		Path: sp.path,
		Marshal: func(v any) ([]byte, error) {
			var result bytes.Buffer
			properties := v.(map[string]string)
			propertiesCamelCase := make(map[string]string)
			for k, v := range properties {
				propertiesCamelCase[utils.ToCamelCase(k)] = v
			}
			err := sp.content.Execute(&result, propertiesCamelCase)
			if err != nil {
				return nil, err
			}
			return result.Bytes(), nil
		},
		Unmarshal: func(data []byte, v any) error {
			var err error
			p := properties.MustLoadString(string(data))
			if mapPtr, ok := v.(*map[string]string); ok {
				*mapPtr = p.Map()
				return nil
			}
			err = p.Decode(v)
			if err != nil {
				return err
			}
			return nil
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create data writer: %w", err)
	}
	templateStr, err := utils.ConvertToTemplate(string(d.DataBytes), utils.TemplateConfig{
		KeepComments:    true,
		KeepEmptyLines:  true,
		KeepIndentation: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert to template: %w", err)
	}
	content, err := template.New("server_properties").Parse(templateStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	sp.DataWriter = d
	sp.content = content
	return sp, nil
}

func (sp *ServerProperties) GetAll() map[string]string {
	return sp.Data
}

func (sp *ServerProperties) Get(key string) string {
	return sp.Data[key]
}

func (sp *ServerProperties) SetAll(data map[string]string) error {
	for k, v := range data {
		err := sp.Set(k, v, false)
		if err != nil {
			return err
		}
	}
	return sp.Write()
}

func (sp *ServerProperties) Set(k, v string, write bool) error {
	_, kExist := sp.Data[k]
	if !kExist {
		return fmt.Errorf("unsupported key [%s]", k)
	}
	filter, fExist := serverPropertiesFilter[k]
	if fExist && !utils.Contains(filter, v) {
		return fmt.Errorf("unsupported value [%s]", v)
	}
	sp.Data[k] = v
	if write {
		return sp.Write()
	} else {
		return nil
	}
}

func (sp *ServerProperties) GetServerName() string {
	return sp.Get("server-name")
}
func (sp *ServerProperties) SetServerName(name string) error {
	return sp.Set("server-name", name, true)
}
