reverse

``` diff
 func ReverseSlice(slice []int) []int {
 	if len(slice) == 0 {
 		return []int{}
 	}
 	result := make([]int, len(slice))
-	copy(result, slice)
-	left := 0
-	right := len(result) - 1
-	for left < right {
-		result[left], result[right] = result[right], result[left]
-		left++
-		right--
+	for i, v := range slice {
+		result[len(slice)-1-i] = v
 	}
 	return result
 }
```