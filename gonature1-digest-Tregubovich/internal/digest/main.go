package digest

import (
	"math/cmplx"
	"math/rand/v2"
	"strings"
	"unsafe"
)

// GetCharByIndex returns the i-th character from the given string.
func GetCharByIndex(str string, idx int) rune {
	runeIdx := 0
	for _, r := range str {
		if runeIdx == idx {
			return r
		}
		runeIdx++
	}
	panic("index out of range")
}

// GetStringBySliceOfIndexes returns a string formed by concatenating specific characters from the input string based
// on the provided indexes.
func GetStringBySliceOfIndexes(str string, indexes []int) string {
	var s strings.Builder
	s.Grow(len(indexes))
	strRune := []rune(str)
	for i := range indexes {
		s.WriteRune(strRune[indexes[i]])
	}
	return s.String()
}

// ShiftPointer shifts the given pointer by the specified number of bytes using unsafe.Add.
func ShiftPointer(pointer **int, shift int) {
	*pointer = (*int)(unsafe.Add(unsafe.Pointer(*pointer), shift))
}

// IsComplexEqual compares two complex numbers and determines if they are equal.
func IsComplexEqual(a, b complex128) bool {
	return a == b || cmplx.Abs(a-b) < 1e-6
}

// GetRootsOfQuadraticEquation returns two roots of a quadratic equation ax^2 + bx + c = 0.
func GetRootsOfQuadraticEquation(a, b, c float64) (complex128, complex128) {
	aСmplx := complex(a, 0)
	bСmplx := complex(b, 0)
	cСmplx := complex(c, 0)
	d := bСmplx*bСmplx - 4*aСmplx*cСmplx
	return (-bСmplx + cmplx.Sqrt(d)) / 2 * aСmplx, (-bСmplx - cmplx.Sqrt(d)) / 2 * aСmplx
}

// Sort sorts in-place the given slice of integers in ascending order.
func Sort(source []int) {
	QuickSort(source, 0, len(source)-1)
}

// Я решил сортировать через QuickSort, так как его проще написать без лишней аллокации, чем MergeSort
func QuickSort(source []int, l, r int) {
	if l < r {
		mid := Partition(source, l, r)
		QuickSort(source, l, mid-1)
		QuickSort(source, mid+1, r)
	}
}

func Partition(source []int, l, r int) int {
	idx := rand.IntN(r-l+1) + l
	source[idx], source[r] = source[r], source[idx]
	m := source[r]
	i := l - 1
	for j := l; j <= r-1; j++ {
		if source[j] < m {
			i++
			source[i], source[j] = source[j], source[i]
		}
	}
	source[i+1], source[r] = source[r], source[i+1]
	return i + 1
}

// ReverseSliceOne in-place reverses the order of elements in the given slice.
func ReverseSliceOne(s []int) {
	n := len(s)
	for i := 0; i < n/2; i++ {
		s[i], s[n-i-1] = s[n-i-1], s[i]
	}
}

// ReverseSliceTwo returns a new slice of integers with elements in reverse order compared to the input slice.
// The original slice remains unmodified.
func ReverseSliceTwo(s []int) []int {
	sCopy := make([]int, len(s))
	copy(sCopy, s)
	ReverseSliceOne(sCopy)
	return sCopy
}

// // SwapPointers swaps the values of two pointers.
func SwapPointers(a, b *int) {
	*a, *b = *b, *a
}

// IsSliceEqual compares two slices of integers and returns true if they contain the same elements in the same order.
func IsSliceEqual(a, b []int) bool {
	n := len(a)
	if len(b) != n {
		return false
	}
	for i := range n {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// DeleteByIndex deletes the element at the specified index from the slice and returns a new slice.
// The original slice remains unmodified.
func DeleteByIndex(s []int, idx int) []int {
	if idx < 0 || idx >= len(s) {
		panic("Index out of range")
	}
	sNew := []int{}
	for i := range s {
		if i != idx {
			sNew = append(sNew, s[i])
		}
	}
	return sNew
}
