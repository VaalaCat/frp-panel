package utils

type Equaler[E any] interface {
	Equal(E) bool
}

// DiffResult 结构体存储了两个切片比较后的差异结果。
type DiffResult[T any] struct {
	NotInArr1 []T // 存在于 arr2 中但不在 arr1 中的元素 (arr2 - arr1)
	NotInArr2 []T // 存在于 arr1 中但不在 arr2 中的元素 (arr1 - arr2)
	InBoth    []T // 同时存在于 arr1 和 arr2 中的元素 (arr1 ∩ arr2)
}

func Diff[T Equaler[T]](arr1, arr2 []T) DiffResult[T] {
	result := DiffResult[T]{
		NotInArr1: make([]T, 0),
		NotInArr2: make([]T, 0),
		InBoth:    make([]T, 0),
	}

	matchedInArr2 := make([]bool, len(arr2))

	for _, item1 := range arr1 {
		foundInArr2 := false
		for i, item2 := range arr2 {
			if !matchedInArr2[i] && item1.Equal(item2) {
				result.InBoth = append(result.InBoth, item1)
				matchedInArr2[i] = true
				foundInArr2 = true
				break
			}
		}
		if !foundInArr2 {
			result.NotInArr2 = append(result.NotInArr2, item1)
		}
	}

	for i, item2 := range arr2 {
		if !matchedInArr2[i] {
			result.NotInArr1 = append(result.NotInArr1, item2)
		}
	}

	return result
}
