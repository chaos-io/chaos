package logs

type Config struct {
	InitFields   map[string]interface{} `json:"initFields"`
	File         FileConfig             `json:"file"`
	Level        string                 `json:"level" default:"debug"`
	Encode       string                 `json:"encode" default:"console"`
	LevelPattern string                 `json:"levelPattern" default:""`
	Output       string                 `json:"output" default:"console"`
	LevelPort    int                    `json:"levelPort" default:"0"`
}

type FileConfig struct {
	Path       string `json:"path" default:"./logs/app.log"`
	Encode     string `json:"encode"  default:"json"`
	MaxSize    int    `json:"maxSize" default:"100"`
	MaxBackups int    `json:"maxBackups" default:"10"`
	MaxAge     int    `json:"maxAge" default:"30"`
}
