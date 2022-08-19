package bear

// FmtPrettyPrint tells the error to format its Error()
func FmtPrettyPrint(on bool) ErrOption {
	return func(e *Error) {
		e.prettyPrint = on
	}
}

// FmtNoStack turns off the stack trace for Error()
func FmtNoStack(on bool) ErrOption {
	return func(e *Error) {
		e.noStack = on
	}
}

// FmtNoParents turns off the parents for Error()
func FmtNoParents(on bool) ErrOption {
	return func(e *Error) {
		e.noParents = on
	}
}

// FmtNoMsg turns off the message for Error()
func FmtNoMsg(on bool) ErrOption {
	return func(e *Error) {
		e.noMsg = on
	}
}

// FmtNoID turns off the error id for Error()
func FmtNoID(on bool) ErrOption {
	return func(e *Error) {
		e.noID = on
	}
}
