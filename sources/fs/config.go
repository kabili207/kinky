package fs

type SourceConfig struct {
	Extensions        []string
	Folder            string
	Recursive         bool
	EnableNSFWSuffix  bool
	EnableContentText bool
	EnableNSFWFolder  bool
}
