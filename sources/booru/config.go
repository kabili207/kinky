package booru

type SourceConfig struct {
	Dialect     string
	Service     string
	Tags        string
	Content     string
	ExtractTags []string `yaml:"extract_tags"`
}
