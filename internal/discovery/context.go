package discovery

type Evidence struct {
	Path    string
	Kind    string
	Snippet string
}

type RepoContext struct {
	Root          string
	Evidence      []Evidence
	DetectedTools []string
}
