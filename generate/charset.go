package generate

// Charset constants define different sets of characters that can be used for generating random strings.
const (
	LatinUpperCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	LatinLowerCharset = "abcdefghijklmnopqrstuvwxyz"
	DigitsCharset     = "0123456789"
	SymbolsCharset    = "!\"#$%&'()*+,-./:;<=>?@[\\]^`{|}~"
	AllCharset        = LatinUpperCharset + LatinLowerCharset + DigitsCharset + SymbolsCharset
)
