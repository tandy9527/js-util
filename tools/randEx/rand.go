package randEx

// randEx 随机函数扩展
import (
	"fmt"
	"math/rand"
	"sort"
)

// RandChoice 随机返回一个元素
func RandChoice[T any](sequence []T) T {
	return sequence[rand.Intn(len(sequence))]
}

// RandChoiceByWeight  根据权重随机返回一个元素
func RandChoiceByWeight[T any](sequence []T, weights []int) T {
	n := len(sequence)
	if n == 0 || n != len(weights) {
		panic("Invalid input: sequence and weights must be same non-zero length.")
	}

	var sum int64
	for _, w := range weights {
		sum += int64(w)
	}
	if sum == 0 {
		panic("All weights are zero.")
	}

	dice := rand.Intn(int(sum))
	for i, w := range weights {
		dice -= int(w)
		if dice < 0 {
			return sequence[i]
		}
	}
	panic(fmt.Sprintf("RandChoiceByWeight logic error: dice=%d sum=%d", dice, sum))
}

// GetResultByGate 数组中的第一个元素为权重,第二个元素为总数
func GetResultByGate(gateInfo []int) bool {
	return gateInfo[0] >= rand.Intn(gateInfo[1])+1
}

// 根據權重array決定對應result, 會回傳是array中的第幾個以及對應result結果
func GetResultByWeight(awardList []int, weightList []int) (int, int) {
	prefixSum := make([]int, len(weightList))
	prefixSum[0] = weightList[0]
	for i := 1; i < len(weightList); i++ {
		prefixSum[i] = prefixSum[i-1] + weightList[i]
	}
	randNum := rand.Intn(prefixSum[len(prefixSum)-1]) + 1
	index := sort.SearchInts(prefixSum, randNum)
	return index, awardList[index]
}

// 根據權重array決定回傳是array中的第幾個index結果
func GetIndexByWeight(weightList []int) int {
	prefixSum := make([]int, len(weightList))
	prefixSum[0] = weightList[0]
	for i := 1; i < len(weightList); i++ {
		prefixSum[i] = prefixSum[i-1] + weightList[i]
	}
	randNum := rand.Intn(prefixSum[len(prefixSum)-1]) + 1
	index := sort.SearchInts(prefixSum, randNum)
	return index
}
