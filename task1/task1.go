package main

import (
	"fmt"
	"slices"
)

func main() {
	arr := []int{1, 5, 5, 2, 5, 2, 3, 3}
	onlyOnceNum(arr)
	robberyAndPlundering(arr)
	listNode1 := createList([]int{1, 2, 3})
	listNode2 := createList([]int{2, 3, 4})
	listNode3 := mergeTwoLists(listNode1, listNode2)
	printList(listNode3)

	arrList := permute([]int{1, 2, 3})
	fmt.Println(arrList)
	reverseString([]byte("hello"))
	fmt.Println(solution(5))

	removeDuplicates([]int{1, 2, 3, 5, 5})

	fmt.Println(merge([][]int{{1, 2, 3}, {2, 5, 6}, {7, 8, 9}}))
}

// 136. 只出现一次的数字
func onlyOnceNum(arr []int) {
	fmt.Println("136. 只出现一次的数字, 给出的数组是：", arr)
	numMap := make(map[int]int)
	for _, item := range arr {
		numMap[item]++
	}
	for key, value := range numMap {
		if value == 1 {
			fmt.Println("只出现一次的数字是：", key)
		}
	}
}

// 198. 打家劫舍
func robberyAndPlundering(arr []int) int {
	max1 := 0
	max2 := 0
	for i := 0; i < len(arr); i++ {
		if i%2 == 0 {
			max1 += arr[i]
		} else {
			max2 += arr[i]
		}
	}
	fmt.Println("198. 打家劫舍, 收到的数组是：", arr, ",max1=", max1, ",max2=", max2)
	return max(max1, max2)

}

// 21.合并两个有序链表
func mergeTwoLists(list1, list2 *ListNode) *ListNode {
	if list1 == nil {
		return list2
	}
	if list2 == nil {
		return list1
	}
	if list1.Val < list2.Val {
		list1.Next = mergeTwoLists(list1.Next, list2)
		return list1
	} else {
		list2.Next = mergeTwoLists(list1, list2.Next)
		return list2
	}
}

// 46.全排列
func permute(nums []int) [][]int {
	var result [][]int
	backtrack(&result, nums, 0)
	return result
}

func backtrack(result *[][]int, nums []int, start int) {
	if start == len(nums) {
		// 复制当前排列到结果中
		tmp := make([]int, len(nums))
		copy(tmp, nums)
		*result = append(*result, tmp)
		return
	}

	for i := start; i < len(nums); i++ {
		// 交换元素位置
		nums[start], nums[i] = nums[i], nums[start]
		// 递归处理下一个位置
		backtrack(result, nums, start+1)
		// 恢复交换（回溯）
		nums[start], nums[i] = nums[i], nums[start]
	}
}

// 344. 反转字符串
func reverseString(s []byte) {
	left := 0
	right := len(s) - 1
	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
	fmt.Println("344. 反转字符串, 反转后的数组是：", string(s))
}

// 69. x 的平方根
func solution(x int) int {
	result := 0
	for i := 0; i*i <= x; i++ {
		result = i
	}
	return result
}

// 26. 删除有序数组中的重复项
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	slow := 0
	for fast := 1; fast < len(nums); fast++ {
		if nums[fast] != nums[slow] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	fmt.Println("26. 删除有序数组中的重复项, 删除后的数组是：", nums[:slow+1])
	return slow + 1
}

// 56. 合并区间
func merge(intervals [][]int) (ans [][]int) {
	slices.SortFunc(intervals, func(p, q []int) int { return p[0] - q[0] }) // 按照左端点从小到大排序
	for _, p := range intervals {
		m := len(ans)
		if m > 0 && p[0] <= ans[m-1][1] { // 可以合并
			ans[m-1][1] = max(ans[m-1][1], p[1]) // 更新右端点最大值
		} else { // 不相交，无法合并
			ans = append(ans, p) // 新的合并区间
		}
	}
	return
}

// 430. 扁平化多级双向链表
func flatten(root *Node) *Node {
	if root == nil {
		return nil
	}
	cur := root
	for cur != nil {
		if cur.Child == nil {
			cur = cur.Next
			continue
		}
		// 存在子链表，递归
		next := cur.Next
		flattenChild := flatten(cur.Child)
		cur.Next = flattenChild
		flattenChild.Prev = cur
		cur.Child = nil // 不要忘记置空child！
		// 连接原来的next
		if next != nil {
			for cur.Next != nil {
				cur = cur.Next
			}
			cur.Next = next
			next.Prev = cur
		}
		cur = cur.Next
	}
	return root
}

// 729. 我的日程安排表 I
type MyCalendar struct {
	bookings [][]int
}

func Constructor() MyCalendar {
	return MyCalendar{}
}

func (mc *MyCalendar) Book(startTime int, endTime int) bool {
	for _, b := range mc.bookings {
		if b[0] < endTime && startTime < b[1] {
			return false
		}
	}
	mc.bookings = append(mc.bookings, []int{startTime, endTime})
	return true
}

type Node struct {
	Val   int
	Prev  *Node
	Next  *Node
	Child *Node
}

func createList(arr []int) *ListNode {
	if len(arr) == 0 {
		return nil
	}
	head := &ListNode{
		Val: arr[0],
	}
	cur := head
	for i := 1; i < len(arr); i++ {
		cur.Next = &ListNode{
			Val: arr[i],
		}
		cur = cur.Next
	}
	return head
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func printList(head *ListNode) {
	for head != nil {
		fmt.Print(head.Val, " -> ")
		head = head.Next
	}
	fmt.Println("nil")
}
