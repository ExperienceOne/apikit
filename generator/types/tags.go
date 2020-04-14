package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/ExperienceOne/apikit/generator/openapi"

	"github.com/pkg/errors"
)

type RegexValidator struct {
	Tag   string
	Regex string
}

func generateValidationTag(required bool, tags []string) string {

	if required {
		return strings.Join(tags, ",")
	} else {

		if len(tags) == 0 {
			return ""
		}

		tag, tags := generateDeltaTag(tags)
		if len(tags) > 0 {
			tag += "," + strings.Join(tags, ",")
		}

		return tag
	}
}

func generateDeltaTag(tags []string) (string, []string) {

	var tag string
	if tags, tag = FilterTag("regex", tags); tag != "" {
		return fmt.Sprintf("omitempty,%s", tag), tags
	} else if tags, tag = FilterTag("dive", tags); tag != "" {
		return fmt.Sprintf("omitempty,gt=0,%s", tag), tags
	} else if tags, tag = FilterTag("email", tags); tag != "" {
		return fmt.Sprintf("omitempty,%s", tag), tags
	} else {
		return "omitempty", tags
	}
}

func FilterTag(tagSub string, tags []string) ([]string, string) {

	var foundTag string
	newTags := make([]string, 0)
	for _, tag := range tags {
		if strings.Contains(tag, tagSub) {
			foundTag = tag
			continue
		}
		newTags = append(newTags, tag)
	}
	return newTags, foundTag
}

func generateIntegerRestriction(min, max *float64, exclusiveMinimum, exclusiveMaximun bool) []string {

	tags := make([]string, 0)
	if min != nil && max != nil {
		if exclusiveMinimum {
			*min = *min + 1
		}

		if exclusiveMaximun {
			*max = *max - 1
		}

		tags = append(tags, "min="+strconv.FormatInt(int64(*min), 10))
		tags = append(tags, "max="+strconv.FormatInt(int64(*max), 10))
	} else if max != nil {
		if exclusiveMaximun {
			*max = *max - 1
		}

		tags = append(tags, "max="+strconv.FormatInt(int64(*max), 10))
	} else if min != nil {
		if exclusiveMinimum {
			*min = *min + 1
		}

		tags = append(tags, "min="+strconv.FormatInt(int64(*min), 10))
	}
	return tags
}

var tags int32

func generateRegexRestriction(field, pattern string) (*RegexValidator, error) {

	if pattern != "" {
		_, err := regexp.Compile(pattern)
		if err != nil {
			return nil, errors.Wrapf(err, "couldn't compile regex (%v)", pattern)
		}
		n := atomic.AddInt32(&tags, 1)
		regexValidator := new(RegexValidator)
		regexValidator.Regex = pattern
		if field == "" {
			field = "regex" + strconv.Itoa(int(n))
		}
		regexValidator.Tag = field
		return regexValidator, nil
	}
	return nil, nil
}

func generateStringRestriction(maxLength, minLength *int64, format string) []string {

	tags := make([]string, 0)
	if minLength != nil && maxLength != nil {
		tags = append(tags, "min="+strconv.FormatInt(*minLength, 10))
		tags = append(tags, "max="+strconv.FormatInt(*maxLength, 10))
	} else if maxLength != nil {
		tags = append(tags, "max="+strconv.FormatInt(*maxLength, 10))
	} else if minLength != nil {
		tags = append(tags, "min="+strconv.FormatInt(*minLength, 10))
	}
	if format == openapi.EMAIL {
		tags = append(tags, "email")
	}
	return tags
}
