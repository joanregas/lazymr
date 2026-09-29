package gitlab

type File struct {
	OldPath       string
	NewPath       string
	AMode         string
	BMode         string
	Diff          string
	NewFile       bool
	RenamedFile   bool
	DeletedFile   bool
	GeneratedFile bool
}
