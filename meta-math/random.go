package metamath

import (
	"math"
	"math/rand"
	metaerror "meta/meta-error"
	"meta/queue"
	"time"
)

var (
	metaRand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

var RandomCountMax = 100
var RandomElementNumMax = 1000

// RandomWeightedData 带权重的数据结构
type RandomWeightedData[T any] struct {
	Data   T
	Weight int
}

// GetRandomInt 获取一个[min, max]范围内的随机整数
// 注意多线程访问时会竞争锁，因此效率较低
func GetRandomInt(min, max int) int {
	return metaRand.Intn(max-min+1) + min
}

// GetRandomInt 获取一个(min, max)范围内的随机浮点数
func GetRandomFloat64(min, max float64) float64 {
	return metaRand.Float64()*(max-min) + min
}

// Shuffle 对切片进行乱序
func Shuffle[T any](slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := GetRandomInt(0, i)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// RandomSlice 从切片中随机选择count个元素
func RandomSlice[T any](slice []T, count int) []T {
	if count > len(slice) {
		count = len(slice)
	}
	Shuffle(slice)
	return slice[:count]
}

// RandomRange 从范围[min, max]中随机选择count个不重复的数字
// 如果有效的不足，则会尽可能返回
func RandomRange(min, max, count int, allowDuplicate bool) ([]int, error) {
	var result []int
	if max < min || count < 1 {
		return result, nil
	}
	if allowDuplicate {
		for i := 0; i < count; i++ {
			result = append(result, GetRandomInt(min, max))
		}
		return result, nil
	}
	var resultTemp []int
	for i := min; i <= max; i++ {
		resultTemp = append(resultTemp, i)
	}
	return RandomSlice(resultTemp, count), nil
}

// RandomWeight 从带权重的数据中随机选择count个元素
// 如果有效的不足，则返回-1以及错误
// 成功随机则返回所选的元素索引
func RandomWeightOneIndex[T any](data []RandomWeightedData[T]) (int, error) {
	if len(data) == 0 {
		return -1, metaerror.New("RandomWeightOneIndex no data")
	}
	weightSum := 0
	for _, d := range data {
		if d.Weight > 0 {
			weightSum += d.Weight
		}
	}
	if weightSum == 0 {
		return -1, metaerror.New("RandomWeightOneIndex weightSum is 0")
	}
	randomWeight := GetRandomInt(1, weightSum)
	for i, d := range data {
		if d.Weight > 0 {
			randomWeight -= d.Weight
			if randomWeight <= 0 {
				return i, nil
			}
		}
	}
	return -1, metaerror.New("unexpected error in random selection")
}

// RandomWeight 从带权重的数据中随机选择count个元素
// 如果有效的不足，则会尽可能返回
// Alone 表示每一次选择是否都是独立的随机（有可能重复）
func RandomWeight[T any](data []RandomWeightedData[T], count int, alone bool) ([]T, error) {
	if count < 1 || count > RandomCountMax {
		return nil, metaerror.New("RandomWeight count out of range [1, %d]", RandomCountMax)
	}
	var result []T
	if alone || count == 1 {
		for i := 0; i < count; i++ {
			index, err := RandomWeightOneIndex(data)
			if err != nil {
				return nil, err
			}
			result = append(result, data[index].Data)
		}
		return result, nil
	}
	if count > len(data) {
		count = len(data)
	}

	type WeightKey struct {
		Index int
		Key   float64
	}

	pq := queue.NewPriorityQueue(
		func(a, b WeightKey) bool {
			return a.Key > b.Key
		},
	)

	weightSum := float64(0)
	weightLimit := float64(0)

	for i, d := range data {
		if pq.Len() < count {
			randRate := GetRandomFloat64(0, 1)
			weight := math.Pow(randRate, 1/float64(d.Weight))
			pq.Add(WeightKey{Index: i, Key: weight})
			continue
		}
		if weightSum == 0 {
			minWeightKey := pq.TopHighestPriority().Key
			weightLimit = Log(minWeightKey, GetRandomFloat64(0, 1))
		}
		if weightSum+float64(d.Weight) < weightLimit {
			weightSum += float64(d.Weight)
			continue
		}

		weightSum = 0

		minWeightKey := pq.TopHighestPriority().Key
		minWeightKeyPow := math.Pow(minWeightKey, float64(d.Weight))
		iWeightKey := math.Pow(GetRandomFloat64(minWeightKeyPow, 1), 1/float64(d.Weight))
		pq.Pop()
		pq.Add(WeightKey{Index: i, Key: iWeightKey})
	}

	for !pq.IsEmpty() {
		index := pq.TopHighestPriority().Index
		if index >= 0 && index < len(data) {
			result = append(result, data[index].Data)
		}
		pq.Pop()
	}
	return result, nil
}
