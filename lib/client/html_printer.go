package client

import (
	"fmt"
)

// Convert the ArticleResponse into an HTML string.
func (r ArticleResponse) ToHTML() string {
	return fmt.Sprintf(`<div class="article"> %s </div>`, r.Article.toHTML())
}
func (a Article) toHTML() string {
	return a.Content.toHTML()
}
func (c Content) toHTML() string {
	res := ""
	for _, child := range c.Children {
		res += child.toHTML() + "\n"
	}
	return res
}

func (n MdAstNode) RecurseOnChildren() string {
	res := ""
	for _, child := range n.Children {
		res += child.toHTML() + "\n"
	}
	return res
}

func (n MdAstNode) toHTML() string {
	if n.Identifier != "" {
		switch n.Identifier {
		case Center:
			return fmt.Sprintf(`<div class="center"> %s </p>`, n.RecurseOnChildren())
		case Figure:
			return fmt.Sprintf(`<div class="figure"> %s </p>`, n.RecurseOnChildren())
		case Title:
			return fmt.Sprintf(`<div class="title"> %s </p>`, n.RecurseOnChildren())
		}
	}
	switch n.Type {
	case Blockquote:
		return fmt.Sprintf(`<blockquote> %s </blockquote>`, n.RecurseOnChildren())
	case Break:
		return `<br>`
	case Code:
		return fmt.Sprintf(`<pre><code> %s </code></pre>`, n.RecurseOnChildren())
	// case Definition
	case Emphasis:
		return fmt.Sprintf(`<em> %s </em>`, n.RecurseOnChildren())
	case Heading:
		d := n.Depth
		if d < 1 || d > 6 {
			fmt.Println("[WARNING]: Heading depth not in range")
			d = 1
		}
		return fmt.Sprintf(`<h%d> %s </h%d>`, n.Depth, n.RecurseOnChildren(), n.Depth)
	// Sounds scary security wise
	// case Html:
	// 	return n.Value
	case Image:
		return fmt.Sprintf(`<img title="%s" alt="%s" src="%s" />`, n.Title, n.Alt, n.Url)
	// case ImageReference
	case InlineCode:
		return fmt.Sprintf(`<code> %s </code>`, n.RecurseOnChildren())
	case Link:
		// TODO do we need to rewrite relative links?
		return fmt.Sprintf(`<a href="%s"> %s </a>`, n.Url, n.RecurseOnChildren())
	// case LinkReference
	case List:
		if n.Ordered {
			return fmt.Sprintf(`<ol> %s </ol>`, n.RecurseOnChildren())
		} else {
			return fmt.Sprintf(`<ul> %s </ul>`, n.RecurseOnChildren())
		}
	case ListItem:
		return fmt.Sprintf(`<li> %s </li>`, n.RecurseOnChildren())
	case Paragraph:
		return fmt.Sprintf(`<p> %s </p>`, n.RecurseOnChildren())

	case Strong:
		return fmt.Sprintf(`<strong> %s </strong>`, n.RecurseOnChildren())
	case Sub:
		return fmt.Sprintf(`<sub> %s </sub>`, n.RecurseOnChildren())
	case Sup:
		return fmt.Sprintf(`<sup> %s </sup>`, n.RecurseOnChildren())
	case Span:
		return fmt.Sprintf(`<span> %s </span>`, n.RecurseOnChildren())
	case Text:
		return n.Value
	case Zone:
		return fmt.Sprintf(`<div class="%s">%s</div>`, n.Identifier, n.RecurseOnChildren())
	case ThematicBreak:
		return fmt.Sprintf(`<hr/> %s`, n.RecurseOnChildren())
	// TODO footnotes?
	default:
		fmt.Printf("[WARNING]: Unsupported element: %s \n", n.Type)
		return fmt.Sprintf(`<div=class="error">Unsupported element: %s </div> %s`, n.Type, n.RecurseOnChildren())
	}
}
