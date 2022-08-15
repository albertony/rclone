// Package path complements standard library package with same name with
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
// The slash is inserted unconditionally. This makes it suitable for building a
// path from basic name elements not containing slashes, or joining a path element
// known to not end with a slash with a name or other path element known to not
// start with a slash.
//
// This function is simply calling [strings.Join] with slash ("/") as hard-coded
// delimiter:
//
//	strings.Join(elems, "/")
//
// The result is identical to a simple concatenation:
//
//	a + "/" + b + "/" + c
//
// If one element ends with a slash and/or the following element starts with a
// slash, there will two or three consecutive slashes in the joint. If it is known
// which of the elements do have a leading or trailing slash, then a tailored
// concatenation statement (a + b + "/" + c) will be perfect. If it is unknown
// whether the elements contain leading/trailing slashes, then you should
// use [JoinOnly] instead.
//
// The result will not be cleaned either, i.e. multiple consecutive slashes in
// either of the elements, or elements such as "." or "..", will be kept
// unresolved. Call [path.Clean] from the standard library on the result, or
// use [path.Join] from the standard library instead, if such processing is
// needed.
//
// The performance penalty is very minimal, it's directly comparable to a single-
// statement concatenation. This function is implemented as a single call to
// [strings.Join], which is a simple utility function based on the
// [strings.Builder] type, which in turn is a performance oriented wrapper over
// a byte slice buffer. Assuming standard compiler optimizations are enabled,
// the function call to this function will be inlined and also most of the calls
// that strings.Join does to strings.Builder functions. There is only a single
// memory allocation, of minimal size, identical to what you get with a plain
// single-statement concatenation.
func Create(elem ...string) string {
	return strings.Join(elem, "/")
}

// JoinOnly concatenates any number of path/url elements, while ensuring slash
// separated.
//
// The slash is inserted conditionally, unlike [Create] but more like the standard
// library [path.Join]. If the first element ends with a slash and/or the second
// element starts with a slash, the joined result will only have a single slash
// between the parts. But it does not perform clean on the result, like [path.Join]
// from the standard library does, nor does it parse it as an URL. This means if
// there are multiple consecutive slashes in either of the elements, then so will
// also the result, and also any "." or ".." elements will also be kept unresolved.
// Call [path.Clean] from the standard library on the result, or use [path.Join]
// from the standard library instead, if such processing is needed.
//
// The result when called with two non-zero length arguments is the the same as
// you would get from:
//
//	strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(b, "/")
//
// And for three non-zero length arguments arguments, it is the same as:
//
//	strings.TrimSuffix(a, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(b, "/"), "/") + "/" + strings.TrimPrefix(c, "/")
//
// The performance penalty is very minimal, it's directly comparable to [Create],
// which effectively just an alias for [strings.Join], and also to a basic single-
// statement concatenation. The implementation is based on the [strings.Builder]
// type, which is a performance oriented wrapper over a byte slice buffer, same
// as used by the [strings.Join] function internally. Assuming standard compiler
// optimizations are enabled, the calls to strings.Builder functions will be
// inlined. There is only a single memory allocation, of minimal size, identical
// to [Create], same as what you get with a plain single-statement concatenation.
//
// This avoids breaking URLs: When working with complete urls, [path.Join] cannot
// be used because it will "clean" the two slashes following the scheme component,
// e.g. change "https://" into "https:/", and also it will remove any trailing path
// separators, which may make a difference in some cases. You could use [net.URL],
// but is slower and more cumbersome to use for simple concatenations (at least
// until go1.19 where a new function url.JoinPath was introduced).
//
// If it is known that the elements does not have leading or trailing slashes,
// then you could use [Create] instead, which is a simpler implementation although
// with minimal performance improvement. For more general paths where clean may be
// required, you should consider [path.Join], as mentioned above.
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

// Join concatenates any number of path elements, while ensuring slash separated,
// ignoring empty elements and cleaning the result.
//
// This is simply an alias for standard library [path.Join], and since the
// compiler will inline it there will be no performance loss.
//
// For joining elements that are already known to be clean, you will get better
// performance by using the simpler functions [JoinOnly] or [Create] instead.
//
// Be careful if using this on URLs! It will "clean" the two slashes following the
// scheme component, e.g. change "https://" into "https:/", and also it will remove
// any trailing path separators, which may make a difference in some cases. Use
// [JoinOnly2] instead to avoid these issues, if you can do without the cleaning.
//
// Se standard library documentation on [path.Clean] for explanation
// of what cleaning does (https://pkg.go.dev/path#Clean).
func Join(elem ...string) string {
	return path.Join(elem...)
}
