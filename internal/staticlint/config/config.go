package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config хранит значения параметров линтера и позволяет получить
// их через геттеры.
type Config struct {
	parameters *parameters
}

// Builder реализует методы для загрузки значений параметров.
type Builder struct {
	parameters *parameters
	err        error
}

type parameters struct {
	AnalyzersNames []string `json:"analyzers"`
}

const configFileName = "staticlint.json"

// NewBuilder возвращает указатель на новый экземпляр Builder.
func NewBuilder() *Builder {
	return &Builder{
		parameters: &parameters{
			AnalyzersNames: []string{
				"bodyclose",
				"errcheck",
				"osexit",
				"asmdecl",
				"assign",
				"atomic",
				"atomicalign",
				"bools",
				"buildssa",
				"buildtag",
				"cgocall",
				"composite",
				"copylock",
				"ctrlflow",
				"deepequalerrors",
				"directive",
				"errorsas",
				"fieldalignment",
				"findcall",
				"framepointer",
				"httpresponse",
				"ifaceassert",
				"inspect",
				"loopclosure",
				"lostcancel",
				"nilfunc",
				"nilness",
				"pkgfact",
				"printf",
				"reflectvaluecompare",
				"shadow",
				"shift",
				"sigchanyzer",
				"sortslice",
				"stdmethods",
				"stringintconv",
				"structtag",
				"testinggoroutine",
				"tests",
				"timeformat",
				"unmarshal",
				"unreachable",
				"unsafeptr",
				"unusedresult",
				"unusedwrite",
				"usesgenerics",
				"SA1000",
				"SA1001",
				"SA1002",
				"SA1003",
				"SA1004",
				"SA1005",
				"SA1006",
				"SA1007",
				"SA1008",
				"SA1010",
				"SA1011",
				"SA1012",
				"SA1013",
				"SA1014",
				"SA1015",
				"SA1016",
				"SA1017",
				"SA1018",
				"SA1019",
				"SA1020",
				"SA1021",
				"SA1023",
				"SA1024",
				"SA1025",
				"SA1026",
				"SA1027",
				"SA1028",
				"SA1029",
				"SA1030",
				"SA2000",
				"SA2001",
				"SA2002",
				"SA2003",
				"SA3000",
				"SA3001",
				"SA4000",
				"SA4001",
				"SA4003",
				"SA4004",
				"SA4005",
				"SA4006",
				"SA4008",
				"SA4009",
				"SA4010",
				"SA4011",
				"SA4012",
				"SA4013",
				"SA4014",
				"SA4015",
				"SA4016",
				"SA4017",
				"SA4018",
				"SA4019",
				"SA4020",
				"SA4021",
				"SA4022",
				"SA4023",
				"SA4024",
				"SA4025",
				"SA4026",
				"SA4027",
				"SA4028",
				"SA4029",
				"SA4030",
				"SA4031",
				"SA5000",
				"SA5001",
				"SA5002",
				"SA5003",
				"SA5004",
				"SA5005",
				"SA5007",
				"SA5008",
				"SA5009",
				"SA5010",
				"SA5011",
				"SA5012",
				"SA6000",
				"SA6001",
				"SA6002",
				"SA6003",
				"SA6005",
				"SA9001",
				"SA9002",
				"SA9003",
				"SA9004",
				"SA9005",
				"SA9006",
				"SA9007",
				"SA9008",
				"ST1001",
				"ST1003",
				"ST1005",
				"ST1006",
				"ST1008",
				"ST1011",
				"ST1012",
				"ST1013",
				"ST1015",
				"ST1016",
				"ST1017",
				"ST1018",
				"ST1019",
				"ST1020",
				"ST1021",
				"ST1022",
				"ST1023",
				"S1000",
				"S1001",
				"S1002",
				"S1003",
				"S1004",
				"S1005",
				"S1006",
				"S1007",
				"S1008",
				"S1009",
				"S1010",
				"S1011",
				"S1012",
				"S1016",
				"S1017",
				"S1018",
				"S1019",
				"S1020",
				"S1021",
				"S1023",
				"S1024",
				"S1025",
				"S1028",
				"S1029",
				"S1030",
				"S1031",
				"S1032",
				"S1033",
				"S1034",
				"S1035",
				"S1036",
				"S1037",
				"S1038",
				"S1039",
				"S1040",
				"U1000",
			},
		},
	}
}

// SetDefaultAnalyzersNames устанавливает имена анализаторов по умолчанию.
func (b *Builder) SetDefaultAnalyzersNames(names ...string) *Builder {
	b.parameters.AnalyzersNames = names

	return b
}

// LoadFile загружает значения из файла конфигурации. Если файл не существует,
// используются значения по умолчанию.
func (b *Builder) LoadFile() *Builder {
	pwd, err := getPWD()
	if err != nil {
		return b
	}

	data, err := os.ReadFile(filepath.Join(pwd, configFileName))
	if err != nil {
		return b
	}

	b.err = json.Unmarshal(data, &b.parameters)

	return b
}

// Build возвращает Config для чтения загруженных значений параметров.
func (b *Builder) Build() (*Config, error) {
	return &Config{b.parameters}, b.err
}

// AnalyzersNames возвращает имена анализаторов.
func (c *Config) AnalyzersNames() []string {
	return c.parameters.AnalyzersNames
}

func getPWD() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	return filepath.Dir(ex), nil
}
