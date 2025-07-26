package problems

type Problem struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Difficulty         string   `json:"difficulty"`
	Tags               []string `json:"tags"`
	LanguagesSupported []string `json:"languages_supported"`
}

type SelectedProblem struct {
	Problem
	Language string `json:"language"`
}
