// Package utils provides utility functions.
package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-homedir"
)

// RemoveFrontmatter removes the front matter header of a markdown file.
func RemoveFrontmatter(content []byte) []byte {
	if frontmatterBoundaries := detectFrontmatter(content); frontmatterBoundaries[0] == 0 {
		return content[frontmatterBoundaries[1]:]
	}
	return content
}

var yamlPattern = regexp.MustCompile(`(?m)^---\r?\n(\s*\r?\n)?`)

func detectFrontmatter(c []byte) []int {
	if matches := yamlPattern.FindAllIndex(c, 2); len(matches) > 1 {
		return []int{matches[0][0], matches[1][1]}
	}
	return []int{-1, -1}
}

// ExpandPath expands tilde and all environment variables from the given path.
func ExpandPath(path string) string {
	s, err := homedir.Expand(path)
	if err == nil {
		return os.ExpandEnv(s)
	}
	return os.ExpandEnv(path)
}

// WrapCodeBlock wraps a string in a code block with the given language.
func WrapCodeBlock(s, language string) string {
	return "```" + language + "\n" + s + "```"
}

var markdownExtensions = []string{
	".md", ".mdown", ".mkdn", ".mkd", ".markdown",
}

// Unicode superscript and subscript mappings
var superscripts = map[rune]rune{
	'0': '\u2070', // ⁰
	'1': '\u00B9', // ¹
	'2': '\u00B2', // ²
	'3': '\u00B3', // ³
	'4': '\u2074', // ⁴
	'5': '\u2075', // ⁵
	'6': '\u2076', // ⁶
	'7': '\u2077', // ⁷
	'8': '\u2078', // ⁸
	'9': '\u2079', // ⁹
	'+': '\u207A', // ⁺
	'-': '\u207B', // ⁻
	'=': '\u207C', // ⁼
	'(': '\u207D', // ⁽
	')': '\u207E', // ⁾
	'n': '\u207F', // ⁿ
}

var subscripts = map[rune]rune{
	'0': '\u2080', // ₀
	'1': '\u2081', // ₁
	'2': '\u2082', // ₂
	'3': '\u2083', // ₃
	'4': '\u2084', // ₄
	'5': '\u2085', // ₅
	'6': '\u2086', // ₆
	'7': '\u2087', // ₇
	'8': '\u2088', // ₈
	'9': '\u2089', // ₉
	'+': '\u208A', // ₊
	'-': '\u208B', // ₋
	'=': '\u208C', // ₌
	'(': '\u208D', // ₍
	')': '\u208E', // ₎
}

var supPattern = regexp.MustCompile(`<sup[^>]*>([^<]*)</sup>`)
var subPattern = regexp.MustCompile(`<sub[^>]*>([^<]*)</sub>`)

// ProcessSuperscript converts <sup> and <sub> HTML tags to Unicode characters.
// This works around Glamour's HTML sanitization which strips these tags.
func ProcessSuperscript(markdown string) string {
	// Process superscript tags
	markdown = supPattern.ReplaceAllStringFunc(markdown, func(match string) string {
		// Extract content between tags
		content := supPattern.FindStringSubmatch(match)
		if len(content) < 2 {
			return match
		}
		// Convert each character to superscript
		var result strings.Builder
		for _, r := range content[1] {
			if sup, ok := superscripts[r]; ok {
				result.WriteRune(sup)
			} else {
				result.WriteRune(r)
			}
		}
		return result.String()
	})

	// Process subscript tags
	markdown = subPattern.ReplaceAllStringFunc(markdown, func(match string) string {
		content := subPattern.FindStringSubmatch(match)
		if len(content) < 2 {
			return match
		}
		var result strings.Builder
		for _, r := range content[1] {
			if sub, ok := subscripts[r]; ok {
				result.WriteRune(sub)
			} else {
				result.WriteRune(r)
			}
		}
		return result.String()
	})

	return markdown
}

// IsMarkdownFile returns whether the filename has a markdown extension.
func IsMarkdownFile(filename string) bool {
	ext := filepath.Ext(filename)

	if ext == "" {
		// By default, assume it's a markdown file.
		return true
	}

	for _, v := range markdownExtensions {
		if strings.EqualFold(ext, v) {
			return true
		}
	}

	// Has an extension but not markdown
	// so assume this is a code file.
	return false
}

// GlamourStyle returns a glamour.TermRendererOption based on the given style.
func GlamourStyle(style string, isCode bool) glamour.TermRendererOption {
	var styleConfig ansi.StyleConfig

	switch style {
	case styles.AutoStyle:
		if lipgloss.HasDarkBackground() {
			styleConfig = styles.DarkStyleConfig
		} else {
			styleConfig = styles.LightStyleConfig
		}
	case styles.DarkStyle:
		styleConfig = styles.DarkStyleConfig
	case styles.LightStyle:
		styleConfig = styles.LightStyleConfig
	case styles.PinkStyle:
		styleConfig = styles.PinkStyleConfig
	case styles.NoTTYStyle:
		styleConfig = styles.NoTTYStyleConfig
	case styles.DraculaStyle:
		styleConfig = styles.DraculaStyleConfig
	case styles.TokyoNightStyle:
		styleConfig = styles.DraculaStyleConfig
	default:
		return glamour.WithStylesFromJSONFile(style)
	}

	// If we are rendering a pure code block, we need to modify the style to
	// remove the indentation.
	if isCode {
		var margin uint
		styleConfig.CodeBlock.Margin = &margin
	}

	return glamour.WithStyles(styleConfig)
}
