package hera_v1_codegen

type CodeDirectoryMetadata struct {
	Map map[string]CodeFilesMetadata
}

type CodeFilesMetadata struct {
	GoCodeFiles   GoCodeFiles
	SQLCodeFiles  SQLCodeFiles
	YamlCodeFiles YamlCodeFiles
	JsCodeFiles   JsCodeFiles
	CssCodeFiles  CssCodeFiles
	HtmlCodeFiles HtmlCodeFiles
}

type JsCodeFiles struct {
	Files            []JsCodeFile
	DirectoryImports []string
}

type JsCodeFile struct {
	FileName  string
	Extension string
	Imports   []string
	Contents  string
}
type CssCodeFiles struct {
	Files []CssCodeFile
}

type CssCodeFile struct {
	FileName string
}
type HtmlCodeFiles struct {
	Files []HtmlCodeFile
}

type HtmlCodeFile struct {
	FileName string
}

type SQLCodeFiles struct {
	Files []SQLCodeFile
}

type SQLCodeFile struct {
	FileName string
	Contents string
}

type YamlCodeFiles struct {
	Files []YamlCodeFile
}

type YamlCodeFile struct {
	FileName string
	Contents string
}

type GoCodeFiles struct {
	Files            map[string]GoCodeFile
	DirectoryImports []string
}

type GoCodeFile struct {
	FileName    string
	PackageName string
	Imports     []string
	Functions   map[string]FunctionInfo
	Variables   map[string]GoVar  // Variable name and its type
	Constants   map[string]GoVar  // Constant name and its type
	Structs     map[string]string // Struct name and its definition
	Contents    string
}

type GoVar struct {
	Type    string
	Value   string
	Content string
}
