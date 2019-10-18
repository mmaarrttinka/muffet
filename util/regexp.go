package util

import "regexp"

// CompileRegexps compiles regular expressions of strings.
func CompileRegexps(ss []string) ([]*regexp.Regexp, error) {
	rs := make([]*regexp.Regexp, 0, len(ss))

	for _, s := range ss {
		r, err := regexp.Compile(s)

		if err != nil {
			return nil, err
		}

		rs = append(rs, r)
	}

	return rs, nil
}
