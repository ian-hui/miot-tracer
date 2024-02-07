package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type ConsistentHash struct {
	nodes    map[int]string // 虚拟节点与实际节点的映射
	keys     []int          // 已排序的节点哈希值列表
	replicas int            // 每个节点的虚拟节点数量
}

// New 创建一个一致性哈希的实例
func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		nodes:    make(map[int]string),
		replicas: replicas,
	}
}

func Hash(key string) int {
	return int(crc32.ChecksumIEEE([]byte(key)))
}

// 添加节点
func (c *ConsistentHash) AddNode(node string) {
	for i := 0; i < c.replicas; i++ {
		hash := Hash(node + strconv.Itoa(i))
		c.nodes[hash] = node
		c.keys = append(c.keys, hash)
	}
	sort.Ints(c.keys) // 保持哈希环的有序性
}

// 移除节点
func (c *ConsistentHash) RemoveNode(node string) {
	for i := 0; i < c.replicas; i++ {
		hash := Hash(node + strconv.Itoa(i))
		delete(c.nodes, hash)
		// 在keys中找到并删除hash
		index := sort.SearchInts(c.keys, hash)
		if index < len(c.keys) && c.keys[index] == hash {
			c.keys = append(c.keys[:index], c.keys[index+1:]...)
		}
	}
}


// GetNode 返回给定键的哈希值最近的节点
func (c *ConsistentHash) GetNode(key string) string {
	if len(c.keys) == 0 {
		return ""
	}
	hash := Hash(key)
	idx := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hash
	})
	if idx == len(c.keys) {
		idx = 0
	}
	return c.nodes[c.keys[idx]]
}
