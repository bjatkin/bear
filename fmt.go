package bear

// FmtPrettyPrint tells the error to format its error
func FmtPrettyPrint(on bool) ErrOption {
	return func(e *Error) {
		e.prettyPrint = on
	}
}

// FmtNoStack turns off the stack trace for the error
func FmtNoStack(on bool) ErrOption {
	return func(e *Error) {
		e.noStack = on
	}
}

// FmtNoParents turns off the parents for error
func FmtNoParents(on bool) ErrOption {
	return func(e *Error) {
		e.noParents = on
	}
}

// FmtNoMsg turns off the message for error
func FmtNoMsg(on bool) ErrOption {
	return func(e *Error) {
		e.noMsg = on
	}
}
