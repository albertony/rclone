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
		a, b, expectedCreate, expectedTrimConcat, expectedJoinOnly, expectedJoinOnlyClean, expectedStdJoin string
	}{
		{"", "", "/", "/", "", ".", ""}, // Note that stdpath.Clean returns "." on empty input, but stdpath.Join returns "" before calling stdpath.Clean in this case
		{"/", "", "//", "/", "/", "/", "/"},
		{"", "/", "//", "/", "/", "/", "/"},
		{"/", "/", "///", "/", "/", "/", "/"},
		{"//", "//", "/////", "///", "///", "/", "/"}, // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a", "", "a/", "a/", "a", "a", "a"},
		{"", "b", "/b", "/b", "b", "b", "b"},
		{"a", "b", "a/b", "a/b", "a/b", "a/b", "a/b"},
		{"a", "/", "a//", "a/", "a/", "a", "a"}, // JoinOnly differs from JoinClean and StdJoin cases (they remove trailing slash)
		{"/", "b", "//b", "/b", "/b", "/b", "/b"},
		{"a/", "/b", "a///b", "a/b", "a/b", "a/b", "a/b"},
		{"a/b", "c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"/a/", "/b/c/", "/a///b/c/", "/a/b/c/", "/a/b/c/", "/a/b/c", "/a/b/c"},                                   // JoinOnly differs from JoinClean and StdJoin cases (they remove trailing slash)
		{"a//", "b", "a///b", "a//b", "a//b", "a/b", "a/b"},                                                       // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a//", "/b", "a////b", "a//b", "a//b", "a/b", "a/b"},                                                     // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a//", "//b", "a/////b", "a///b", "a///b", "a/b", "a/b"},                                                 // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"//a//", "//b//c//", "//a/////b//c//", "//a///b//c//", "//a///b//c//", "/a/b/c", "/a/b/c"},               // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a/b/", "../../../xyz", "a/b//../../../xyz", "a/b/../../../xyz", "a/b/../../../xyz", "../xyz", "../xyz"}, // JoinOnly differs from JoinClean and StdJoin cases (they resolve "." and ".." paths)
		{"https://", "example.com/", "https:///example.com/", "https://example.com/", "https://example.com/", "https:/example.com", "https:/example.com"},                                                                                                                                                               // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
		{"https://", "/example.com/", "https:////example.com/", "https://example.com/", "https://example.com/", "https:/example.com", "https:/example.com"},                                                                                                                                                             // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
		{"https://example.com/root/dir/", "/subdir2/sausage2/", "https://example.com/root/dir///subdir2/sausage2/", "https://example.com/root/dir/subdir2/sausage2/", "https://example.com/root/dir/subdir2/sausage2/", "https:/example.com/root/dir/subdir2/sausage2", "https:/example.com/root/dir/subdir2/sausage2"}, // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
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
		got = stdpath.Clean(JoinOnly(test.a, test.b))
		assert.Equal(t, test.expectedJoinOnlyClean, got, what)

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
		a, b, c, expectedCreate, expectedTrimConcat, expectedJoinOnly, expectedJoinOnlyClean, expectedStdJoin string
	}{
		{"", "", "", "//", "//", "", ".", ""}, // Note that stdpath.Clean returns "." on empty input, but stdpath.Join returns "" before calling stdpath.Clean in this case
		{"/", "", "", "///", "//", "/", "/", "/"},
		{"", "/", "", "///", "//", "/", "/", "/"},
		{"", "", "/", "///", "//", "/", "/", "/"},
		{"/", "/", "", "////", "//", "/", "/", "/"},
		{"", "/", "/", "////", "//", "/", "/", "/"},
		{"/", "/", "/", "/////", "//", "/", "/", "/"},
		{"//", "//", "//", "////////", "////", "////", "/", "/"}, // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a", "", "", "a//", "a//", "a", "a", "a"},
		{"", "b", "", "/b/", "/b/", "b", "b", "b"},
		{"", "", "c", "//c", "//c", "c", "c", "c"},
		{"a", "b", "", "a/b/", "a/b/", "a/b", "a/b", "a/b"},
		{"", "b", "c", "/b/c", "/b/c", "b/c", "b/c", "b/c"},
		{"a", "b", "c", "a/b/c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"/", "b", "c", "//b/c", "/b/c", "/b/c", "/b/c", "/b/c"},
		{"a", "/", "c", "a///c", "a//c", "a/c", "a/c", "a/c"},
		{"a", "b", "/", "a/b//", "a/b/", "a/b/", "a/b", "a/b"}, // JoinOnly differs from JoinClean and StdJoin cases (they remove trailing slash)
		{"a/", "/b/", "/c", "a///b///c", "a/b/c", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "", "a/b/c/", "a/b/c/", "a/b/c", "a/b/c", "a/b/c"},
		{"a", "b/c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a/b", "c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a", "b/c", "d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"a", "b", "c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d", "a/b/c/d"},
		{"/a/", "/b/", "/c/d/", "/a///b///c/d/", "/a/b/c/d/", "/a/b/c/d/", "/a/b/c/d", "/a/b/c/d"},                             // JoinOnly differs from JoinClean and StdJoin cases (they remove trailing slash)
		{"a//", "b", "c", "a///b/c", "a//b/c", "a//b/c", "a/b/c", "a/b/c"},                                                     // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a//", "/b", "c", "a////b/c", "a//b/c", "a//b/c", "a/b/c", "a/b/c"},                                                   // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a//", "//b", "c", "a/////b/c", "a///b/c", "a///b/c", "a/b/c", "a/b/c"},                                               // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"//a//", "//b//", "//c//d//", "//a/////b/////c//d//", "//a///b///c//d//", "//a///b///c//d//", "/a/b/c/d", "/a/b/c/d"}, // JoinOnly differs from JoinClean and StdJoin cases (because it trims at most a single slash from both ends of a joint)
		{"a/b", "../../../xyz", "1/./2/../3", "a/b/../../../xyz/1/./2/../3", "a/b/../../../xyz/1/./2/../3", "a/b/../../../xyz/1/./2/../3", "../xyz/1/3", "../xyz/1/3"},                                                                                                                                                                                       // JoinOnly differs from JoinClean and StdJoin cases (they resolve "." and ".." paths)
		{"https://", "example.com", "/subdir2/sausage2/", "https:///example.com//subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https:/example.com/subdir2/sausage2", "https:/example.com/subdir2/sausage2"},                                                                                         // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
		{"https://", "/example.com/", "/subdir2/sausage2/", "https:////example.com///subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https://example.com/subdir2/sausage2/", "https:/example.com/subdir2/sausage2", "https:/example.com/subdir2/sausage2"},                                                                                     // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
		{"https://example.com/root/dir/", "/subdir2/sausage2/", "/docs/", "https://example.com/root/dir///subdir2/sausage2///docs/", "https://example.com/root/dir/subdir2/sausage2/docs/", "https://example.com/root/dir/subdir2/sausage2/docs/", "https:/example.com/root/dir/subdir2/sausage2/docs", "https:/example.com/root/dir/subdir2/sausage2/docs"}, // JoinOnly differs from JoinClean and StdJoin cases (they replace "https://" with "https:/" and remove trailing slash)
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
		got = stdpath.Clean(JoinOnly(test.a, test.b, test.c))
		assert.Equal(t, test.expectedJoinOnlyClean, got, what)

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
	// - Concat, Create, TrimConcat and JoinOnly have fastest speed,
	//   in most cases 70-80% improvement compared to path.Join, but
	//   in this specific case with two short arguments down to 40%
	//   improvement.
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
	// - Speed similar to the above.
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
	// - Speed similar to the above.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op, same as Create.
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
	// - Speed similar to the above.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op,
	//   same as TrimConcat, Create and Concat.
	// - When running with "-gcflags -m" the compiler is inlining calls to
	//   strings.Builder functions within the JoinOnly implementation.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2)
	}
	result = r
}

func BenchmarkTwoShortNamesJoinOnlyClean(b *testing.B) {
	// Expect:
	// - Slower than the above, but slightly faster than Join,
	//   around 10-20% improvement.
	// - When running with -benchmem there will be 1 allocs/op and 3 B/op,
	//   same as the above.
	in1, in2 := GetTwoShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2))
	}
	result = r
}

func BenchmarkTwoShortNamesJoin(b *testing.B) {
	// Expect:
	// - Slowest speed.
	// - When running with -benchmem there will be 2 allocs/op and 6 B/op.
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

func BenchmarkTwoLongUrlPathsJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 144 B/op.
	in1, in2 := GetTwoLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2))
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

func BenchmarkTwoLongUrlPathsEndingWithExistingJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 144 B/op.
	in1, in2 := GetTwoLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2))
	}
	result = r
}

func BenchmarkTwoLongUrlPathsEndingWithExistingSlashJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 192 B/op.
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

func BenchmarkThreeShortNamesJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 5 B/op.
	in1, in2, in3 := GetThreeShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3))
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

func BenchmarkThreeLongUrlPathsJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 192 B/op.
	in1, in2, in3 := GetThreeLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3))
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

func BenchmarkThreeLongUrlPathsEndingWithExistingSlashJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 192 B/op.
	in1, in2, in3 := GetThreeLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3))
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

//
// Testing 5 argument variants with short names not containing slashes.
// The Concat and Create will produce clean result with single slash in the joints.
// The TrimPrefix/TrimSuffix will return the input argument unchanged.
//

func GetFiveShortNames() (string, string, string, string, string) {
	return "a", "b", "c", "d", "e"
}

func BenchmarkFiveShortNamesConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 16 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3 + "/" + in4 + "/" + in5
	}
	result = r
}

func BenchmarkFiveShortNamesCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 16 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveShortNamesTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 16 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in3, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in4, "/"), "/") + "/" + strings.TrimPrefix(in5, "/")
	}
	result = r
}

func BenchmarkFiveShortNamesJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 16 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveShortNamesJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 16 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3, in4, in5))
	}
	result = r
}

func BenchmarkFiveShortNamesJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 2 allocs/op and 32 B/op.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveShortNamesStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3, in4, in5 := GetFiveShortNames()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3, in4, in5)
	}
	result = r
}

//
// Testing 5 argument variants with url including the double slash separator
// following the scheme component, and with path elements that are "clean"
// so that we are only joining an element without ending slash with an element
// without starting slash.
// The Concat and Create will produce clean result with single slash in the joints.
// The TrimPrefix/TrimSuffix will return the input argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".
//

func GetFiveLongUrlPaths() (string, string, string, string, string) {
	return "https://example.com/root/dir", "subdir2/sausage2", "docs", "user/manual", "public/master"
}

func BenchmarkFiveLongUrlPathsConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3 + "/" + in4 + "/" + in5
	}
	result = r
}

func BenchmarkFiveLongUrlPathsCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in3, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in4, "/"), "/") + "/" + strings.TrimPrefix(in5, "/")
	}
	result = r
}
func BenchmarkFiveLongUrlPathsJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 240 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3, in4, in5))
	}
	result = r
}

func BenchmarkFiveLongUrlPathsJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 320 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPaths()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3, in4, in5)
	}
	result = r
}

//
// Testing 5 argument variants with url including the double slash separator
// following the scheme component, and with path elements so that an element
// ending with a slash is joined with an element starting with a slash.
// The Concat and Create will produce triple slash in the joints.
// The TrimPrefix/TrimSuffix will return a slice and not only the input
// argument unchanged.
// Standard library Clean, and Join which calls Clean, will break the url
// by changing prefix "https://" to "http:/".

func GetFiveLongUrlPathsEndingWithExistingSlash() (string, string, string, string, string) {
	return "https://example.com/root/dir/", "/subdir2/sausage2/", "/docs/", "/user/manual/", "/public/master/"
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 96 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = in1 + "/" + in2 + "/" + in3 + "/" + in4 + "/" + in5
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashCreate(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 96 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Create(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashTrimConcat(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = strings.TrimSuffix(in1, "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in2, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in3, "/"), "/") + "/" + strings.TrimPrefix(strings.TrimSuffix(in4, "/"), "/") + "/" + strings.TrimPrefix(in5, "/")
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashJoinOnly(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 1 allocs/op and 80 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = JoinOnly(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashJoinOnlyClean(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 3 allocs/op and 240 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Clean(JoinOnly(in1, in2, in3, in4, in5))
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashJoin(b *testing.B) {
	// Expect:
	// - When running with -benchmem there will be 4 allocs/op and 368 B/op.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = Join(in1, in2, in3, in4, in5)
	}
	result = r
}

func BenchmarkFiveLongUrlPathsEndingWithExistingSlashStdJoin(b *testing.B) {
	// Expect:
	// - Identical to Join, which is just a simple alias that should be inlined.
	in1, in2, in3, in4, in5 := GetFiveLongUrlPathsEndingWithExistingSlash()
	var r string
	for i := 0; i < b.N; i++ {
		r = stdpath.Join(in1, in2, in3, in4, in5)
	}
	result = r
}
