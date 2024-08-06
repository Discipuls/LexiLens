package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type ParsingError struct {
	message string
}

func (err ParsingError) Error() string {
	return "Parse error: " + err.message
}

func ParsedoOnlineGoogleDictionary(body []byte) (entry WordEntry, err error) {
	entry = WordEntry{}
	bodyReader := bytes.NewReader(body)
	parsedPage, err := html.Parse(bodyReader)
	if err != nil {
		parseErr := ParsingError{message: err.Error()}
		return entry, parseErr
	}

	var processPage func(*html.Node) error
	processPage = func(currentNode *html.Node) error {
		if currentNode.Type == html.ElementNode && currentNode.Data == "b" {
			parseNode(currentNode, &entry)
		}

		if currentNode.Type == html.ElementNode && currentNode.Data == "li" {
			parseNode(currentNode, &entry)
		}

		for childNode := currentNode.FirstChild; childNode != nil; childNode = childNode.NextSibling {
			currentErr := processPage(childNode)
			if currentErr != nil {
				fmt.Println(currentErr)
			}
		}
		return nil
	}

	processPage(parsedPage)
	return entry, nil
}

func parseNode(node *html.Node, entry *WordEntry) error {
	if node.Data == "b" {
		if speechPart := extractSpeechPart(node); speechPart != nil {
			err := entry.addSpeechPart(speechPart)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else if node.Data == "li" {
		if definition := extractDefinition(node); definition != nil {
			err := entry.addDefinition(definition)
			if err != nil {
				fmt.Println(err.Error())
			}
		} else if exampleNode := extractExampleNode(node); exampleNode != nil {
			err := entry.addWordUsageExample()
			if err != nil {
				fmt.Println(err.Error())
			}

			examplePieces := extractSentencePieces(exampleNode)

			err = entry.addWordExamplePieces(examplePieces)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func extractSpeechPart(node *html.Node) *SpeechPartEntry {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		for _, a := range node.Parent.Attr {
			if a.Key == "id" && a.Val == "forEmbed" {
				speechPart := SpeechPartEntry{SpeechPart: node.FirstChild.Data[0 : len(node.FirstChild.Data)-1]}
				return &speechPart
			}
		}
	}
	return nil
}

func extractDefinition(node *html.Node) *WordDefinition {
	if node.Parent.Parent.Data == "ol" && node.Parent.Parent.Parent.Data == "div" {
		for _, a := range node.Parent.Parent.Parent.Attr {
			if a.Key == "class" && a.Val == "std" {
				definition := WordDefinition{Definition: extractSentencePieces(node)}
				//definition := WordDefinition{}
				//definition := WordDefinition{Definition: node.FirstChild.Data[1 : len(node.FirstChild.Data)-1]}
				return &definition
			}
		}
	}
	return nil
}

// func extractDefinitionPieces(definitionNode *html.Node) []SentencePice {
// 	for definitionNode.FirstChild != nil {

// 	}
// }

func extractExampleNode(node *html.Node) *html.Node {
	if node.Parent.Data == "ul" && node.Parent.Parent.Data == "div" {
		for _, a := range node.Parent.Parent.Attr {
			if a.Key == "class" && a.Val == "std" {
				return node
			}
		}
	}
	return nil
}

type checkNode func(node *html.Node) bool

func extractChildNodes(parent *html.Node, checker checkNode) []*html.Node {
	res := make([]*html.Node, 0)
	for fc := parent.FirstChild; fc != nil; fc = fc.NextSibling {
		if checker(fc) {
			res = append(res, fc)
		}
	}
	return res
}

func extractSentencePieces(source *html.Node) []SentencePice {
	sentencePieces := make([]SentencePice, 0)

	for child := source.FirstChild; child != nil; child = child.NextSibling {
		if child.Data == "ul" || child.Data == "div" {
			continue
		}
		if child.Type == html.TextNode {
			sentencePieces = append(sentencePieces, SentencePice{Value: child.Data, ContainsMainWord: source.Data == "em"})
		} else if child.FirstChild != nil && child.Type == html.ElementNode {
			sentencePieces = append(sentencePieces, extractSentencePieces(child)...)
		}
	}
	return clearSentencePieces(sentencePieces)
}

func clearSentencePieces(pieces []SentencePice) []SentencePice {
	res := make([]SentencePice, 0)
	for _, piece := range pieces {
		piece.Value = strings.ReplaceAll(piece.Value, "\n", "")
		piece.Value = strings.ReplaceAll(piece.Value, "                                ", " ")
		for strings.Contains(piece.Value, "  ") {
			piece.Value = strings.ReplaceAll(piece.Value, "  ", " ")
		}
		if piece.Value != "" && piece.Value != " " {
			res = append(res, piece)
		}
	}
	return res
}