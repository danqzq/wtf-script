package interpreter

// TODO: Implement a lexer for proper lexical analysis

func tokenize(expr string) []string {
	var tokens []string
	current := ""
	for _, ch := range expr {
		switch ch {
		case '+', '-', '*', '/', '(', ')':
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(ch))
		case ' ', '\t', '\n':
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
		default:
			current += string(ch)
		}
	}
	if current != "" {
		tokens = append(tokens, current)
	}
	return tokens
}
