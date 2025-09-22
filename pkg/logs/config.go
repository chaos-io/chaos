package logs

type Config struct {
	InitFields   map[string]interface{} `json:"initFields"`
	Level        string                 `json:"level" default:"debug"`    // debug,info,warn,error,fatal
	Encode       string                 `json:"encode" default:"console"` // console,json
	LevelPattern string                 `json:"levelPattern" default:""`
	LevelPort    int                    `json:"levelPort" default:"0"`
	Output       string                 `json:"output" default:"console"` // console,file
	File         FileConfig             `json:"file"`
}

type FileConfig struct {
	Path       string `json:"path" default:"./logs/app.log"`
	Encode     string `json:"encode" default:"json"`
	MaxSize    int    `json:"maxSize" default:"100"`
	MaxBackups int    `json:"maxBackups" default:"10"`
	MaxAge     int    `json:"maxAge" default:"30"`
	Compress   bool   `json:"compress"`
}
