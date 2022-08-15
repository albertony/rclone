// Package path complements standard library package with same name, adding
// some additional utility routines for manipulating slash-separated paths.
//
// The path package should only be used for paths separated by forward slashes,
// such as the paths in URLs. This package does not deal with Windows paths with
// drive letters or backslashes; to manipulate operating system paths, use the
// path/filepath package.
package path

import (
	"path"
	"strings"
)

// Create concatenates any number of path/url elements, inserting a slash between.
//
// The slash is inserted unconditionally. This function is simply an alias for
// [strings.Join] with slash ("/") as hard-coded delimiter:
//
//	strings.Join(elems, "/")
//
// The result is identical to a simple concatenation:
//
//	a + "/" + b + "/" + c
//
// This makes it suitable for building a path from basic name elements not
// containing slashes. It can also be used to join paths containing slash, if it is
// known that no elements (except possibly the last) end with a slash, and no
// elements (except possibly the first) starts with a slash. For cases where this
// is not true, the result will contain two or three consecutive slashes in the
// joint. If it it known which of the elements do have leading and/or trailing
// slash, then a tailored concatenation statement would be needed:
//
//	a + b + "/" + c
//
// The function [JoinOnly] can then instead be used, as it will handle any such
// cases. Also, it handles the case when it is not known whether the elements
// contain leading/trailing slashes, as it checks and inserts slash only when
// necessary. There is very little difference in performance, if any at all, so you
// may in fact simply use [JoinOnly] for all of the above cases.
//
// Unlike what is the case with the standard library function [path.Join], the
// result will not be "cleaned". Standard library functions [path.Clean] or
// [path.Join] must be used for this, see documentation of [JoinOnly] for more
// details.
//
// The performance of this function should be close to optimal, except if hand-
// writing an implementation using low-level pointer operations etc. Performance
// is directly comparable to a corresponding single-statement concatenation. This
// function is implemented as a single call to [strings.Join], which is a simple
// utility function based on the [strings.Builder] type, which in turn is a
// performance oriented wrapper over a byte slice buffer. Assuming standard
// compiler optimizations are enabled, the function call to this function will be
// inlined and also most of the calls that [strings.Join] does to [strings.Builder]
// functions. It will perform only a single memory allocation, of minimal size,
// identical to what you get with a plain single-statement concatenation.
func Create(elem ...string) string {
	return strings.Join(elem, "/")
}

// JoinOnly concatenates any number of path/url elements, while ensuring slash
// separated. Empty elements are ignored.
//
// The slash is inserted conditionally, unlike [Create] but more similar to the
// standard library function [path.Join] (though there are other differences, as
// described below). The result when called with two non-zero length arguments is
// the the same as you would get from:
//
//	strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(b, "/")
//
// And for three non-zero length arguments arguments:
//
//	strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(b, "/"), "/") + "/" + strings.TrimPrefix(c, "/")
//
// If the first element ends with a slash and/or the second element starts with a
// slash, the joined result will only have a single slash between the elements.
// This makes it suitable for building a path from elements where it is not known
// if they contain leading/trailing slashes. If it is known that none of the
// elements start or end with a slash, then [Create] could be used instead. It
// inserts a slash unconditionally. Though, there is very little difference in
// performance, if any at all, so you may in fact simply use [JoinOnly] for these
// cases as well.
//
// Unlike what is the case with the standard library function [path.Join], the
// result will not be "cleaned". This means multiple consecutive slashes in
// either of the elements will be kept, and elements such as "." or ".." will
// not be treated any specially, i.e. not resolved as relative path references.
// If you need this kind of functionality, you can call standard library function
// [path.Clean] on the result, or use [path.Join] from the standard library
// instead. The performance is slightly better when using the combination of
// [JoinOnly] and [path.Clean] compared to [path.Join] from standard library.
// Note that unless you check for empty results from [JoinOnly], there may be
// one difference between these two alternatives: When arguments to [path.Join]
// ends up as empty string, it returns an empty string, but if [path.Clean] is
// called with an empty string as argument, it returns ".".
//
// One important difference from standard library [path.Join], is that [JoinOnly]
// does not break URLs! When working with complete urls, [path.Join] (and
// [path.Clean]) cannot be used because it will "clean" the two slashes following
// the scheme component, e.g. change "https://" into "https:/". Also it will remove
// any trailing path separators, which may make a difference in some cases. You
// could use [net.URL], but is slower and more cumbersome to use for simple
// concatenations (at least until go1.19 where a new function [url.JoinPath] was
// introduced).
//
// The performance of this function should be close to optimal, except if hand-
// writing an implementation using low-level pointer operations etc. Performance
// is directly comparable to a corresponding single-statement concatenation, and
// also to [Create], which effectively just an alias for [strings.Join]. The
// implementation is based on the [strings.Builder] type, which is a performance
// oriented wrapper over a byte slice buffer, same as used by the [strings.Join]
// function internally. Assuming standard compiler optimizations are enabled, the
// calls to strings.Builder functions will be inlined. There will only be a single
// memory allocation, of minimal size, identical what you get with a plain
// single-statement concatenation.
func JoinOnly(elem ...string) string {
	// Alternative simple size estimation, that overestimates instead of finding exact required size:
	/*
		size := 0
		for i := 0; i < len(elem); i++ {
			size += len(elem[i])
		}
		if size == 0 {
			return ""
		}
		size += len(elem) - 1 // Allocate space for worst case, i.e. that a separator must be inserted between each element
	*/
	size := 0
	endslash := false
	for i := 0; i < len(elem); i++ {
		if len(elem[i]) > 0 {
			if elem[i][0] == '/' {
				if endslash {
					size += len(elem[i]) - 1
				} else {
					size += len(elem[i])
				}
			} else {
				if !endslash && size > 0 {
					size++
				}
				size += len(elem[i])
			}
			endslash = elem[i][len(elem[i])-1] == '/'
		}
	}
	var b strings.Builder
	b.Grow(size)

	endslash = false
	for i := 0; i < len(elem); i++ {
		if len(elem[i]) > 0 {
			if elem[i][0] == '/' {
				if endslash {
					b.WriteString(elem[i][1:])
				} else {
					b.WriteString(elem[i])
				}
			} else {
				if !endslash && b.Len() > 0 {
					b.WriteRune('/')
				}
				b.WriteString(elem[i])
			}
			endslash = elem[i][len(elem[i])-1] == '/'
		}
	}

	// TEMP: FOR TESTING!
	/*
		if size != b.Len() {
			fmt.Printf("WRONG SIZE: Expected %v but was %v", size, b.Len())
		} else {
			fmt.Printf("OK SIZE: %v", size)
		}
	*/

	return b.String()
}

// Join concatenates any number of path elements, while ensuring slash separated.
// Empty elements are ignored. The result is cleaned.
//
// This is simply an alias for standard library [path.Join], and since the
// compiler will inline when standard compiler optimizations are enabled, there
// will be no performance difference whichever you use.
//
// For joining elements that are already known to be clean, you will get better
// performance by using the simpler functions [JoinOnly] or [Create] instead.
//
// The result from this function will be cleaned. This means multiple consecutive
// slashes in either of the elements will be kept, and elements such as "." or ".."
// will be resolved as relative path references, using purely lexical processing.
// Se standard library documentation on [path.Clean] for more details
// (https://pkg.go.dev/path#Clean).
//
// Be aware if consider using this on URLs: The clean process performed by
// [path.Clean], called internally from [path.Join], will "clean" the two slashes
// following the scheme component, e.g. change "https://" into "https:/". Also it
// will remove any trailing path separators, which may make a difference in some
// cases. Use [JoinOnly] instead to avoid these issues, if you can do without the
// other parts of the cleaning.
func Join(elem ...string) string {
	return path.Join(elem...)
}
