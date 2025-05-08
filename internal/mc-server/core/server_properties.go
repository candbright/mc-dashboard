package core

import (
	"bytes"
	_ "embed"
	"html/template"
	"path"

	"github.com/candbright/go-server/internal/mc-server/utils"
	"github.com/candbright/go-server/pkg/dw"
	"github.com/magiconair/properties"
	"github.com/pkg/errors"
)

var (
	serverPropertiesFilter = map[string][]string{
		"gamemode": {"survival", "creative", "adventure"},
	}
)

type ServerProperties struct {
	Version string
	rootDir string
	tmplStr string
	*dw.DataWriter[map[string]string]
}

type ServerPropertiesConfig struct {
	Version string
	RootDir string
}

func NewServerProperties(config ServerPropertiesConfig) *ServerProperties {
	sp := &ServerProperties{
		Version: config.Version,
		rootDir: config.RootDir,
	}
	sp.Init()
	return sp
}

func (sp *ServerProperties) FileName() string {
	return "server.properties"
}

func (sp *ServerProperties) FilePath() string {
	return path.Join(sp.rootDir, sp.FileName())
}

func (sp *ServerProperties) Init() {
	d, err := dw.New[map[string]string](dw.Config{
		Path: sp.FilePath(),
		Marshal: func(v any) ([]byte, error) {
			content, err := template.New("server_properties").Parse(sp.tmplStr)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			var result bytes.Buffer
			err = content.Execute(&result, v)
			if err != nil {
				return nil, errors.WithStack(err)
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
				return errors.WithStack(err)
			}
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	templateStr, err := utils.ConvertToTemplate(string(d.DataBytes), utils.TemplateConfig{
		KeepComments:    true,
		KeepEmptyLines:  true,
		KeepIndentation: true,
	})
	if err != nil {
		panic(err)
	}
	sp.DataWriter = d
	sp.tmplStr = templateStr
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
		return errors.Errorf("unsupported key [%s]", k)
	}
	filter, fExist := serverPropertiesFilter[k]
	if fExist && !utils.Contains(filter, v) {
		return errors.Errorf("unsupported value [%s]", v)
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
