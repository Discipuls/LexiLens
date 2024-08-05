package main

import (
	"testing"
)

func TestDodoOnlineGoogleDictionary(t *testing.T) {
	expectedWordEntries := make([]WordEntry, 0)
	expectedRipEntry := WordEntry{
		Word: "rip",
		SpeechParts: []SpeechPartEntry{
			{
				SpeechPart: "Verb",
				Definitions: []WordDefinition{
					{
						Definition: []SentencePice{{Value: "Tear or pull (something) quickly or forcibly away from something or someone", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "a fan tried to ", ContainsMainWord: false},
									{Value: "rip", ContainsMainWord: true},
									{Value: " his pants ", ContainsMainWord: false},
									{Value: "off", ContainsMainWord: false},
									{Value: " during a show", ContainsMainWord: false},
								},
							},
							{
								Pieces: []SentencePice{
									{Value: "countries ", ContainsMainWord: false},
									{Value: "ripped", ContainsMainWord: true},
									{Value: " apart by fighting", ContainsMainWord: false},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "Make a long tear or cut in", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "you've ", ContainsMainWord: false},
									{Value: "ripped", ContainsMainWord: true},
									{Value: " my jacket", ContainsMainWord: false},
								},
							},
							{
								Pieces: []SentencePice{
									{Value: "ripped", ContainsMainWord: true},
									{Value: " jeans", ContainsMainWord: false},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "Make (a hole) by force", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "the truck was struck by lightning and had a hole ", ContainsMainWord: false},
									{Value: "ripped", ContainsMainWord: true},
									{Value: " out of its roof", ContainsMainWord: false},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "Come violently apart; tear", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "he heard something ", ContainsMainWord: false},
									{Value: "rip", ContainsMainWord: true},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "Cut (wood) in the direction of the grain", ContainsMainWord: false}},
					},
					{
						Definition: []SentencePice{{Value: "Move forcefully and rapidly", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "fire ", ContainsMainWord: false},
									{Value: "ripped", ContainsMainWord: true},
									{Value: " through", ContainsMainWord: true},
									{Value: " her bungalow", ContainsMainWord: false},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "Use a program to copy (a sound sequence on a compact disc) on to a computer's hard drive", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "every Beatles song ever made, ", ContainsMainWord: false},
									{Value: "ripped", ContainsMainWord: true},
									{Value: " from my boxed set of CDs", ContainsMainWord: false},
								},
							},
						},
					},
				},
			},
			{
				SpeechPart: "Noun",
				Definitions: []WordDefinition{
					{
						Definition: []SentencePice{{Value: "A dissolute immoral person, esp. a man", ContainsMainWord: false}},
						Examples: []WordUsageExample{
							{
								Pieces: []SentencePice{
									{Value: "“Where is that old ", ContainsMainWord: false},
									{Value: "rip", ContainsMainWord: true},
									{Value: "?” a deep voice shouted", ContainsMainWord: false},
								},
							},
						},
					},
					{
						Definition: []SentencePice{{Value: "A mischievous person, esp. a child", ContainsMainWord: false}},
					},
					{
						Definition: []SentencePice{{Value: "A worthless horse", ContainsMainWord: false}},
					},
				},
			},
		},
	}

	expectedWordEntries = append(expectedWordEntries, expectedRipEntry)

	for _, expectedEntry := range expectedWordEntries {
		actualEntry := doOnlineGoogleDictionary(expectedEntry.Word)
		if expectedEntry.Word != actualEntry.Word {
			t.Fatalf("\nTest failed for word \"%s\": WordEntry.word doesn't match to expected \nExpected: \"%s\"\nGot: \"%s\"",
				expectedEntry.Word, expectedEntry.Word, actualEntry.Word)
		}

		if len(expectedEntry.SpeechParts) != len(actualEntry.SpeechParts) {
			t.Fatalf("\nTest failed for word \"%s\": length of WordEntry.speechParts doesn't match to expected \nExpected: %d\nGot: %d",
				expectedEntry.Word, len(expectedEntry.SpeechParts), len(actualEntry.SpeechParts))
		}

		for i := range expectedEntry.SpeechParts {
			if expectedEntry.SpeechParts[i].SpeechPart != actualEntry.SpeechParts[i].SpeechPart {
				t.Fatalf("\nTest failed for word \"%s\": WordEntry.SpeechParts[%d].speechPart doesn't match to expected"+
					"\nExpected: \"%s\"\nGot: \"%s\"", expectedEntry.Word, i,
					expectedEntry.SpeechParts[i].SpeechPart, actualEntry.SpeechParts[i].SpeechPart)
			}
		}

		for i := range expectedEntry.SpeechParts {
			expectedSpeechPart := expectedEntry.SpeechParts[i]
			actualSpeechPart := actualEntry.SpeechParts[i]
			if len(expectedSpeechPart.Definitions) != len(actualSpeechPart.Definitions) {
				t.Fatalf("\nTest failed for word \"%s\": length of WordEntry.speechParts[%d].definitions doesn't match to expected \nExpected: %d\nGot: %d",
					expectedEntry.Word, i, len(expectedSpeechPart.Definitions), len(actualSpeechPart.Definitions))
			}
		}

		for i := range expectedEntry.SpeechParts {
			expectedDefinitions := expectedEntry.SpeechParts[i].Definitions
			actualDefinitions := actualEntry.SpeechParts[i].Definitions
			for j := range expectedDefinitions {
				if expectedDefinitions[j].Definition[0] != actualDefinitions[j].Definition[0] {
					t.Fatalf("\nTest failed for word \"%s\": \nWordEntry.SpeechParts[%d].speechPart.definitions[%d].definition"+
						" doesn't match to expected"+
						"\nExpected: \"%s\"\nGot: \"%s\"", expectedEntry.Word, i, j,
						expectedDefinitions[j].Definition, actualDefinitions[j].Definition)
				}
			}
		}

		// expectedEntry.speechParts[0].definitions[1].examples =
		//  append(expectedEntry.speechParts[0].definitions[1].examples, WordUsageExample{})
		for i := range expectedEntry.SpeechParts {
			expectedDefinitions := expectedEntry.SpeechParts[i].Definitions
			actualDefinitions := actualEntry.SpeechParts[i].Definitions
			for j := range expectedDefinitions {
				if len(expectedDefinitions[j].Examples) != len(actualDefinitions[j].Examples) {
					t.Fatalf("\nTest failed for word \"%s\": \nlength of WordEntry.SpeechParts[%d].speechPart.definitions[%d].examples"+
						" doesn't match to expected"+
						"\nExpected: %d\nGot: %d", expectedEntry.Word, i, j,
						len(expectedDefinitions[j].Examples), len(actualDefinitions[j].Examples))
				}
			}
		}

		for i := range expectedEntry.SpeechParts {
			expectedDefinitions := expectedEntry.SpeechParts[i].Definitions
			actualDefinitions := actualEntry.SpeechParts[i].Definitions
			for j := range expectedDefinitions {
				expectedExamples := expectedDefinitions[j].Examples
				actualExamples := actualDefinitions[j].Examples
				for k := range expectedExamples {
					expectedExampleString := ""
					for _, examplePiece := range expectedExamples[k].Pieces {
						expectedExampleString += examplePiece.Value
					}

					actualExampleString := ""
					for _, examplePiece := range actualExamples[k].Pieces {
						actualExampleString += examplePiece.Value
					}
					if expectedExampleString != actualExampleString {
						t.Fatalf("\nTest failed for word \"%s\":"+
							" \nexample of WordEntry.SpeechParts[%d].speechPart.definitions[%d].examples[%d]"+
							" doesn't match to expected"+
							"\nExpected: \"%s\"\nGot: \"%s\"", expectedEntry.Word, i, j, k,
							expectedExampleString, actualExampleString)
					}
				}
			}
		}

	}
}
