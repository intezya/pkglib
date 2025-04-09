package generate

// Charset constants define different sets of characters that can be used for generating random strings.
const (
	latinUpperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"                                           // Uppercase Latin letters
	latinLowerCharset = "abcdefghijklmnopqrstuvwxyz"                                           // Lowercase Latin letters
	digitsCharset     = "0123456789"                                                           // Digits from 0 to 9
	symbolsCharset    = "!\"#$%&'()*+,-./:;<=>?@[\\]^`{|}~"                                    // Symbols and special characters
	allCharset        = latinUpperCharset + latinLowerCharset + digitsCharset + symbolsCharset // All characters combined
)

// charset struct holds the different character sets for string generation.
type charset struct {
	LatinUpperCharset string // Uppercase Latin characters
	LatinLowerCharset string // Lowercase Latin characters
	DigitsCharset     string // Digits 0-9
	SymbolsCharset    string // Symbols and special characters
	AllCharset        string // All character sets combined
}

// Charset is a global variable holding the character sets, making them accessible for random string generation.
var Charset = charset{
	LatinUpperCharset: latinUpperCharset,
	LatinLowerCharset: latinLowerCharset,
	DigitsCharset:     digitsCharset,
	SymbolsCharset:    symbolsCharset,
	AllCharset:        allCharset,
}
