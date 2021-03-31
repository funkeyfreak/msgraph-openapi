package filter

import (
	"errors"
)

const (
	ErrMsgFilterFlagPath = "filter type tag only accepts []string as filterData"
	ErrMsgFilterFlagTag  = "filter type path only accepts []string as filterData"
)

var (
	ErrFilterFlagTypeInvalidType          = errors.New("invalid FilterFlagType")
	ErrFiltertypeCannotUseAllorNoneOnPath = errors.New("cannot use MatchType MatchTypeNoneOf or MatchtypeAllOf when matching with stratagy FilterFlagPath")
)

// This enum teslls us how to apply the filter. The types of matches are described below
type MatchType string

const (
	// Will only return items which match all patterns
	MatchTypeAllOf MatchType = "all"
	// Will return any item which match any of the patterns
	MatchTypeAnyOf MatchType = "any"
	// Will only return items which do not match any patterns
	MatchTypeNoneOf MatchType = "none"
	// Will return items which do not match any of the patterns
	MatchtypeInverseOf MatchType = "inverse"
)

func (f MatchType) IsValid() bool {
	switch f {
	case MatchTypeAllOf, MatchTypeAnyOf, MatchTypeNoneOf:
		return true
	}
	return false
}

// TODO: add regex filter
type FilterFlagType string

const (
	FilterFlagTag  FilterFlagType = "tag"
	FilterFlagPath FilterFlagType = "path"
)

func (f FilterFlagType) IsValid() bool {
	switch f {
	case FilterFlagPath, FilterFlagTag:
		return true
	}
	return false
}

func (f FilterFlagType) Error() string {
	switch f {
	case FilterFlagPath:
		return ErrMsgFilterFlagPath
	case FilterFlagTag:
		return ErrMsgFilterFlagTag
	}

	return "filter type does not accept the provided filterData type"
}

type FilterFlag struct {
	FilterType FilterFlagType
	MatchType  MatchType
	FilterData []string
}

func NewFilterFlag(filterType FilterFlagType, matchType MatchType, filterData []string) (*FilterFlag, error) {
	if !filterType.IsValid() {
		return nil, ErrFilterFlagTypeInvalidType
	}

	switch matchType {
	case MatchTypeNoneOf, MatchTypeAnyOf:
		if filterType == FilterFlagPath {
			return nil, ErrFiltertypeCannotUseAllorNoneOnPath
		}
	}

	return &FilterFlag{
		FilterType: filterType,
		FilterData: filterData,
		MatchType:  matchType,
	}, nil
}
