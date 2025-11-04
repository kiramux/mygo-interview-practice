# Challenge 21: Binary Search Implementation

[🔗link](https://app.gointerview.dev/challenge/21)

## BinarySearch

采取左闭右开区间的写法。

## BinarySearchRecursive

由于测试中是按照左闭右闭区间来写的，这里写左闭右闭。

## FindInsertPosition

左闭右开区间写法， `left == right` 时循环终止。此时这个边界就是插入位置。

### 为什么在 `left == right` 时循环终止

考虑边界情况：

当 left = right - 1 时：

> left = 5, right = 6
> mid = 5 + (6-5)>>1 = 5 + 0 = 5

- 如果选择 `left = mid + 1 = 6`，则新状态：`left = 6, right = 6`
- 如果选择 `right = mid = 5`，则新状态：`left = 5, right = 5`

两种情况下都有 left = right，循环终止。

### 为什么这个边界就是插入位置

`left == right` 时：

- `left` 左边所有元素都小于 target
- `right` 右边所有元素都大于等于 target

因此在 left 位置插入 target：

- target > 左边所有元素 ✓
- target ≤ 右边所有元素 ✓
- 数组保持有序 ✓

``` go
func FindInsertPosition(arr []int, target int) int {
	// 左闭右开，终止时 left 和 right 等价
	left := 0
	right := len(arr)
	for left < right {
		mid := left + (right-left)>>1
		if target == arr[mid] {
			right = mid
		} else if target < arr[mid] {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}
```
