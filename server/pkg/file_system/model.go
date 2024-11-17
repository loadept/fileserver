package filesystem

type DirectoryStructure struct {
	Dirs  []Content `json:"dirs"`
	Files []Content `json:"files"`
}

type Content struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}
