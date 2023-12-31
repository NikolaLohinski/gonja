package config

// Config holds plexer and parser parameters
type Config struct {
	// The string marking the beginning of a block. Defaults to '{%'
	BlockStartString string
	// The string marking the end of a block. Defaults to '%}'.
	BlockEndString string
	// The string marking the beginning of a print controlStructure. Defaults to '{{'.
	VariableStartString string
	// The string marking the end of a print controlStructure. Defaults to '}}'.
	VariableEndString string
	// The string marking the beginning of a comment. Defaults to '{#'.
	CommentStartString string
	// The string marking the end of a comment. Defaults to '#}'.
	CommentEndString string
	// If given and a string, this will be used as prefix for line based controlStructures.
	// See also Line ControlStructures.
	LineControlStructurePrefix string
	// If given and a string, this will be used as prefix for line based comments.
	// See also Line ControlStructures.
	LineCommentPrefix string
	// If set to True the XML/HTML autoescaping feature is enabled by default.
	// For more details about autoescaping see Markup.
	// This can also be a callable that is passed the template name
	// and has to return True or False depending on autoescape should be enabled by default.
	AutoEscape bool
	// Whether to be strict about undefined attribute or item in an object and return error
	// or return a nil value on missing data and ignore it entirely
	StrictUndefined bool
}

func New() *Config {
	return &Config{
		BlockStartString:    "{%",
		BlockEndString:      "%}",
		VariableStartString: "{{",
		VariableEndString:   "}}",
		CommentStartString:  "{#",
		CommentEndString:    "#}",
		AutoEscape:          false,
		StrictUndefined:     false,
	}
}

func (cfg *Config) Inherit() *Config {
	return &Config{
		BlockStartString:    cfg.BlockStartString,
		BlockEndString:      cfg.BlockEndString,
		VariableStartString: cfg.VariableStartString,
		VariableEndString:   cfg.VariableEndString,
		CommentStartString:  cfg.CommentStartString,
		CommentEndString:    cfg.CommentEndString,
		AutoEscape:          cfg.AutoEscape,
		StrictUndefined:     cfg.StrictUndefined,
	}
}
