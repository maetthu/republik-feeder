package client

// Based on
// https://github.com/syntax-tree/mdast?tab=readme-ov-file
// but ommiting some

type MdAstType string

const (
	Blockquote     = "blockquote"
	Break          = "break"
	Code           = "code"
	Definition     = "definition"
	Emphasis       = "emphasis"
	Heading        = "heading"
	Html           = "html"
	Image          = "image"
	ImageReference = "imagereference" // not sure about spelling
	InlineCode     = "inlinecode"
	Link           = "link"
	LinkReference  = "linkreference" // not sure about spelling
	List           = "list"
	ListItem       = "listItem"
	Paragraph      = "paragraph"
	Root           = "root"
	Sub            = "sub" // added
	Sup            = "sup" // added
	Strong         = "strong"
	Span           = "span" // added
	Text           = "text"
	ThematicBreak  = "thematicBreak" // not sure about spelling
	Zone           = "zone"          // added
)

type MdAstNode struct {
	Type       string      `json:"type"`
	Identifier Identifier  `json:"identifier"`
	Children   []MdAstNode `json:"children"`
	Depth      int         `json:"depth"`
	Value      string      `json:"value"`
	Alt        string      `json:"alt"`
	Title      string      `json:"title"`
	Url        string      `json:"url"`
	Code       string      `json:"code"`
	Lang       string      `json:"lang"`
	// Meta       string      `json:"meta`
	Label   string `json:"label"`
	Ordered bool   `json:"ordered"`
	// start, spread for list
}
