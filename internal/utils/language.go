package utils

func GetFileExtension(language string) string {
	extensions := map[string]string{
		"go":     "go",
		"python": "py",
		"java":   "java",
		"cpp":    "cpp",
		"c":      "c",
		"js":     "js",
		"ts":     "ts",
		"rust":   "rs",
		"ruby":   "rb",
	}
	if ext, ok := extensions[language]; ok {
		return ext
	}
	return language
}

func GetLanguageFromExtension(extension string) string {
	languages := map[string]string{
		"go":   "go",
		"py":   "python",
		"java": "java",
		"cpp":  "cpp",
		"c":    "c",
		"js":   "javascript",
		"ts":   "typescript",
		"rs":   "rust",
		"rb":   "ruby",
	}
	if lang, ok := languages[extension]; ok {
		return lang
	}
	return extension
}

func GetCompilerCommand(language string) []string {
	commands := map[string][]string{
		"go":     {"go", "build"},
		"python": {"python3"},
		"java":   {"javac"},
		"cpp":    {"g++", "-std=c++17"},
		"c":      {"gcc"},
		"js":     {"node"},
		"ts":     {"tsc"},
		"rust":   {"rustc"},
		"ruby":   {"ruby"},
	}
	if cmd, ok := commands[language]; ok {
		return cmd
	}
	return []string{language}
}

func GetRunCommand(language string) []string {
	commands := map[string][]string{
		"go":     {"./main"},
		"python": {"python3", "main.py"},
		"java":   {"java", "Main"},
		"cpp":    {"./main"},
		"c":      {"./main"},
		"js":     {"node", "main.js"},
		"ts":     {"node", "main.js"},
		"rust":   {"./main"},
		"ruby":   {"ruby", "main.rb"},
	}
	if cmd, ok := commands[language]; ok {
		return cmd
	}
	return []string{"./" + language}
}
