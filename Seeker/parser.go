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

			examplePieces := extractExamplePieces(exampleNode)

			err = entry.addWordExamplePieces(examplePieces)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return nil
}

func extractSpeechPart(node *html.Node) *SpeechPartEntry {
	if node.FirstChild.Type == html.TextNode {
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
				definition := WordDefinition{Definition: node.FirstChild.Data[1 : len(node.FirstChild.Data)-1]}
				return &definition
			}
		}
	}
	return nil
}

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

func extractExamplePieces(exampleNode *html.Node) []WordExamplePice {
	examplePieces := make([]WordExamplePice, 0)

	for exampleNode.FirstChild != nil {
		examplePieceNodes := make([]*html.Node, 0)
		if exampleNode.FirstChild.Type == html.TextNode {

			examplePieceNodes = append(examplePieceNodes, exampleNode)
		} else if exampleNode.FirstChild.Type == html.ElementNode &&
			exampleNode.FirstChild.FirstChild.Type == html.TextNode {

			examplePieceNodes = append(examplePieceNodes, exampleNode.FirstChild)
		} else if exampleNode.FirstChild.Type == html.ElementNode &&
			exampleNode.FirstChild.FirstChild.Type == html.ElementNode &&
			exampleNode.FirstChild.FirstChild.FirstChild.Type == html.TextNode {

			examplePieceNodes = append(examplePieceNodes, exampleNode.FirstChild.FirstChild)

			exampleNode.FirstChild.RemoveChild(exampleNode.FirstChild.FirstChild)
			if exampleNode.FirstChild.FirstChild != nil {
				examplePieceNodes = append(examplePieceNodes, exampleNode.FirstChild)
			}

		}
		for _, examplePieceNode := range examplePieceNodes {
			examplePiece := WordExamplePice{
				Value:            examplePieceNode.FirstChild.Data,
				ContainsMainWord: examplePieceNode.Data == "em"}

			examplePiece.Value = strings.ReplaceAll(examplePiece.Value, "\n", "")
			examplePiece.Value = strings.ReplaceAll(examplePiece.Value, "                                ", " ")
			for strings.Contains(examplePiece.Value, "  ") {
				examplePiece.Value = strings.ReplaceAll(examplePiece.Value, "  ", " ")
			}

			examplePieces = append(examplePieces, examplePiece)
		}
		exampleNode.RemoveChild(exampleNode.FirstChild)

	}
	return examplePieces
}
