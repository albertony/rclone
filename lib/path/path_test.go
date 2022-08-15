package path

import (
	"fmt"
	stdpath "path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTwoPathElements(t *testing.T) {
	for _, test := range []struct {
		a, b, expectedCreate, expectedTrimConcat, expectedJoinOnly, expectedJoinClean, expectedStdJoin string
	}{
		{"", "", "/", "/", "", ".", ""}, // Note that stdpath.Clean returns "." on empty input, but stdpath.Join returns "" before calling stdpath.Clean in this case
		{"/", "", "//", "/", "/", "/", "/"},
		{"", "/", "//", "/", "/", "/", "/"},
		{"/", "/", "///", "/", "/", "/", "/"},
		{"//", "//", "/////", "///", "///", "/", "/"}, // JoinOnly differs from JoinClean and StdJoin cases in that at most a single slash is trimmed from both ends of a joint
		{"a", "", "a/", "a/", "a", "a", "a"},
		{"", "b", "/b", "/b", "b", "b", "b"},
		{"a", "b", "a/b", "a/b", "a/b", "a/b", "a/b"},
		{"a", "/", "a//", "a/", "a/", "a", "a"},
		{"/", "b", "//b", "/b", "/b", "/b", "/b"},
		{"a/", "/b", "a///b", "a/b", "a/b", "a/b", "a/b"},
		{"a/b", "c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"/a/", "/b/c/", "/a///b/c/", "/a/b/c/", "/a/b/c/", "/a/b/c", "/a/b/c"},
		{"//a//", "//b//c//", "//a/////b//c//", "//a///b//c//", "//a///b//c//", "/a/b/c", "/a/b/c"},
		{"a/b/", "../../../xyz", "a/b//../../../xyz", "a/b/../../../xyz", "a/b/../../../xyz", "../xyz", "../xyz"},
		{"https://", "example.com/", "https:///example.com/", "https://example.com/", "https://example.com/", "https:/example.com", "https:/example.com"},                                                                                                                                                               // Note: Join2 and path.Join replaces "https://" with "https:/" and removes trailing slash
		{"https://", "/example.com/", "https:////example.com/", "https://example.com/", "https://example.com/", "https:/example.com", "https:/example.com"},                                                                                                                                                             // Note: Join2 and path.Join replaces "https://" with "https:/" and removes trailing slash
		{"https://example.com/root/dir/", "/subdir2/sausage2/", "https://example.com/root/dir///subdir2/sausage2/", "https://example.com/root/dir/subdir2/sausage2/", "https://example.com/root/dir/subdir2/sausage2/", "https:/example.com/root/dir/subdir2/sausage2", "https:/example.com/root/dir/subdir2/sausage2"}, // Note: Join2 and path.Join replaces "https://" with "https:/" and removes trailing slash
	} {
		what := fmt.Sprintf("Concat(%q, %q)", test.a, test.b)
		got := test.a + "/" + test.b
		assert.Equal(t, test.expectedCreate, got, what)

		what = fmt.Sprintf("Create(%q, %q)", test.a, test.b)
		got = Create(test.a, test.b)
		assert.Equal(t, test.expectedCreate, got, what)

		what = fmt.Sprintf("TrimConcat(%q, %q)", test.a, test.b)
		got = strings.TrimSuffix(test.a, "/") + "/" + strings.TrimPrefix(test.b, "/")
		assert.Equal(t, test.expectedTrimConcat, got, what)

		what = fmt.Sprintf("JoinOnly(%q, %q)", test.a, test.b)
		got = JoinOnly(test.a, test.b)
		assert.Equal(t, test.expectedJoinOnly, got, what)

		what = fmt.Sprintf("path.Clean(JoinOnly(%q, %q)", test.a, test.b)
		got = stdpath.Clean(Join(test.a, test.b))
		assert.Equal(t, test.expectedJoinClean, got, what)

		what = fmt.Sprintf("Join(%q, %q)", test.a, test.b)
		got = Join(test.a, test.b)
		assert.Equal(t, test.expectedStdJoin, got, what)

		what = fmt.Sprintf("stdpath.Join(%q, %q)", test.a, test.b)
		got = stdpath.Join(test.a, test.b)
		assert.Equal(t, test.expectedStdJoin, got, what)
	}
}

func TestThreePathElements(t *testing.T) {
	for _, test := range []struct {
		a, b, c, expectedCreate, expectedTrimConcat, expectedJoinOnly, expectedJoinClean, expectedStdJoin string
	}{
		{"", "", "", "//", "//", "", ".", ""}, // Note that stdpath.Clean returns "." on empty input, but stdpath.Join returns "" before calling stdpath.Clean in this case
		{"/", "", "", "///", "//", "/", "/", "/"},
		{"", "/", "", "///", "//", "/", "/", "/"},
		{"", "", "/", "///", "//", "/", "/", "/"},
		{"/", "/", "", "////", "//", "/", "/", "/"},
		{"", "/", "/", "////", "//", "/", "/", "/"},
		{"/", "/", "/", "/////", "//", "/", "/", "/"},
		{"//", "//", "//", "////////", "////", "////", "/", "/"}, // JoinOnly differs from JoinClean and StdJoin cases in that at most a single slash is trimmed from both ends of a joint
		{"a", "", "", "a//", "a//", "a", "a", "a"},
		{"", "b", "", "/b/", "/b/", "b", "b", "b"},
		{"", "", "c", "//c", "//c", "c", "c", "c"},
		{"a", "b", "", "a/b/", "a/b/", "a/b", "a/b", "a/b"},
		{"", "b", "c", "/b/c", "/b/c", "b/c", "b/c", "b/c"},
		{"a", "b", "c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"/", "b", "c", "//b/c", "/b/c", "/b/c", "/b/c", "/b/c"},
		{"a", "/", "c", "a///c", "a//c", "a/c", "a/c", "a/c"},
		{"a", "b", "/", "a/b//", "a/b/", "a/b/", "a/b", "a/b"},
		{"a/", "/b/", "/c", "a///b///c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "", "a/b/c/", "a/b/c/", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a/b", "c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a", "b/c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a", "b", "c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"/a/", "/b/", "/c/d/", "/a///b///c/d/", "/a/b/c/d/", "/a/b/c/d/", "/a/b/c/d", "/a/b/c/d"},
		{"//a//", "//b//", "//c//d//", "//a/////b/////c//d//", "//a///b///c//d//", "//a///b///c//d//", "/a/b/c/d", "/a/b/c/d"},
		{"a/b", "../../../xyz", "1/./2/../3", "a/b/../../../xyz/1/./2/../3", "a/b/../../../xyz/1/./2/../3", "a/b/../../../xyz/1/./2/../3", "../xyz/1/3", "../xyz/1/3"},
		{"https://", "example.com", "/subdir2/sausage2/", "https:///example.com//subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https:/example.com/subdir2/sausage2", "https:/example.com/subdir2/sausage2"},                                                                                         // Note: Join2 and path.Join replaces "https://" with "https:/" and removes trailing slash
		{"https://", "/example.com/", "/subdir2/sausage2/", "https:////example.com///subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https:/example.com/subdir2/sausage2", "https:/example.com/subdir2/sausage2"},                                                                                     // Note: Join2 and path.Join replaces "https://" with "https:/" and removes trailing slash
		{"https://example.com/root/dir/", "/subdir2/sausage2/", "/docs/", "https://example.com/root/dir///subdir2/sausage2///docs/", "https://example.com/root/dir/subdir2/sausage2/docs/", "https://example.com/root/dir/subdir2/sausage2/docs/", "https:/example.com/root/dir/subdir2/sausage2/docs", "https:/example.com/root/dir/subdir2/sausage2/docs"}, // Note: Join2 and path.Join replaces "https://" with "https:/"
	} {
		what := fmt.Sprintf("Concat(%q, %q, %q)", test.a, test.b, test.c)
		got := test.a + "/" + test.b + "/" + test.c
		assert.Equal(t, test.expectedCreate, got, what)

		what = fmt.Sprintf("Create(%q, %q, %q)", test.a, test.b, test.c)
		got = Create(test.a, test.b, test.c)
		assert.Equal(t, test.expectedCreate, got, what)

		what = fmt.Sprintf("Trim/Concat(%q, %q, %q)", test.a, test.b, test.c)
		got = strings.TrimSuffix(test.a, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(test.b, "/"), "/") + "/" + strings.TrimPrefix(test.c, "/")
		assert.Equal(t, test.expectedTrimConcat, got, what)

		what = fmt.Sprintf("JoinOnly(%q, %q, %q)", test.a, test.b, test.c)
		got = JoinOnly(test.a, test.b, test.c)
		assert.Equal(t, test.expectedJoinOnly, got, what)

		what = fmt.Sprintf("Clean(JoinOnly(%q, %q, %q)", test.a, test.b, test.c)
		got = stdpath.Clean(Join(test.a, test.b, test.c))
		assert.Equal(t, test.expectedJoinClean, got, what)

		what = fmt.Sprintf("Join(%q, %q, %q)", test.a, test.b, test.c)
		got = Join(test.a, test.b, test.c)
		assert.Equal(t, test.expectedStdJoin, got, what)

		what = fmt.Sprintf("stdpath.Join(%q, %q, %q)", test.a, test.b, test.c)
		got = stdpath.Join(test.a, test.b, test.c)
		assert.Equal(t, test.expectedStdJoin, got, what)
	}
}

//
// Benchmarks
//

var result string

//
// Testing 2 argument variants with short names not containing slashes.
// The basic Concat and Create will produce clean result with single slash in
// the joints. The TrimPrefix/TrimSuffix will return the input argument unchanged.
//

func GetTwoShortNames() (string, string) {
	return "a", "b"
}

func BenchmarkTwoShortNamesConcat(b *testing.B) {
	// Expect:
	// - Fastest speed together with Create.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2
	}
	result = r
}

func BenchmarkTwoShortNamesCreate(b *testing.B) {
	// Expect:
	// - Fastest speed together with Concat.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op, same as Concat.
	// - When running with "-gcflags -m" the compiler is "inlining call to Create".
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2)
	}
	result = r
}

func BenchmarkTwoShortNamesTrimConcat(b *testing.B) {
	// Expect:
	// - Close to Create and Concat in speed.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op, same as Create and Concat.
	// - When running with "-gcflags -m" the compiler is inlining calls to
	//   strings.TrimSuffix and strings.TrimPrefix, and in those functions
	//   call to strings.HasSuffix and strings.HasPrefix.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(in2, "/")
	}
	result = r
}

func BenchmarkTwoShortNamesJoinOnly(b *testing.B) {
	// Expect:
	// - Same as TrimConcat in speed, close to Create and Concat.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op, same as TrimConcat, Create and Concat.
	// - When running with "-gcflags -m" the compiler is inlining calls to
	//   strings.Builder functions within the JoinOnly implementation.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2)
	}
	result = r
}

func BenchmarkTwoShortNamesCleanJoin(b *testing.B) {
	// Expect:
	// - Slower by maybe 50%-100% than JoinOnly.
	// - When running with -benchmem there will be 2 allocs/op and 6 B/op, twice of JoinOnly.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2))
	}
	result = r
}

func BenchmarkTwoShortNamesJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 2 allocs/op and 6 B/op.
	// - Slowest, even slower than Join2 (but not by much, and might end up faster in some runs).
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2)
	}
	result = r
}

func BenchmarkTwoShortNamesStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2)
	}
	result = r
}

//
// Testing 2 argument variants with url including the double slash separator
// following the scheme component, and with path elements that are "clean"
// so that we are only joining an element without ending slash with an element
// without starting slash.
// The Concat and Create will produce clean result with single slash in the joints.
// The TrimPrefix/TrimSuffix will return the input argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".
//

func GetTwoLongUrlPaths() (string, string) {
	return "https://example.com/root/dir", "subdir2/sausage2"
}

func BenchmarkTwoLongUrlPathsConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2
	}
	result = r
}

func BenchmarkTwoLongUrlPathsCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(in2, "/")
	}
	result = r
}

func BenchmarkTwoLongUrlPathsJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsCleanJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 192 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2))
	}
	result = r
}

func BenchmarkTwoLongUrlPathsJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 192 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2)
	}
	result = r
}

//
// Testing 2 argument variants with url including the double slash separator
// following the scheme component, and with path elements so that an element
// ending with a slash is joined with an element starting with a slash.
// The Concat and Create will produce triple slash in the joints.
// The TrimPrefix/TrimSuffix will return a slice and not only the input
// argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".

func GetTwoLongUrlPathsEndingWithExistingSlash() (string, string) {
	return "https://example.com/root/dir/", "/subdir2/sausage2/"
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(in2, "/")
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 48 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingCleanJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 192 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2))
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 192 B/op.
	// - Slowest, even slower than Join2 (but not by much, and might end up faster in some runs).
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2)
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2)
	}
	result = r
}

//
// Testing 3 argument variants with short names not containing slashes.
// The Concat and Create will produce clean result with single slash in the joints.
// The TrimPrefix/TrimSuffix will return the input argument unchanged.
//

func GetThreeShortNames() (string, string, string) {
	return "a", "b", "c"
}

func BenchmarkThreeShortNamesConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 5 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3
	}
	result = r
}

func BenchmarkThreeShortNamesCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 5 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeShortNamesTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 5 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(in3, "/")
	}
	result = r
}

func BenchmarkThreeShortNamesJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 5 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeShortNamesCleanJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 2 allocs/op and 10 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2, in3))
	}
	result = r
}

func BenchmarkThreeShortNamesJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 2 allocs/op and 10 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeShortNamesStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3)
	}
	result = r
}

//
// Testing 3 argument variants with url including the double slash separator
// following the scheme component, and with path elements that are "clean"
// so that we are only joining an element without ending slash with an element
// without starting slash.
// The Concat and Create will produce clean result with single slash in the joints.
// The TrimPrefix/TrimSuffix will return the input argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".
//

func GetThreeLongUrlPaths() (string, string, string) {
	return "https://example.com/root/dir", "subdir2/sausage2", "docs"
}

func BenchmarkThreeLongUrlPathsConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3
	}
	result = r
}

func BenchmarkThreeLongUrlPathsCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(in3, "/")
	}
	result = r
}
func BenchmarkThreeLongUrlPathsJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsCleanJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 256 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2, in3))
	}
	result = r
}

func BenchmarkThreeLongUrlPathsJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 256 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3)
	}
	result = r
}

//
// Testing 3 argument variants with url including the double slash separator
// following the scheme component, and with path elements so that an element
// ending with a slash is joined with an element starting with a slash.
// The Concat and Create will produce triple slash in the joints.
// The TrimPrefix/TrimSuffix will return a slice and not only the input
// argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".

func GetThreeLongUrlPathsEndingWithExistingSlash() (string, string, string) {
	return "https://example.com/root/dir/", "/subdir2/sausage2/", "/docs/"
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(in3, "/")
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 64 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashCleanJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 256 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(Join(in1, in2, in3))
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 256 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3)
	}
	result = r
}

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3)
	}
	result = r
}
