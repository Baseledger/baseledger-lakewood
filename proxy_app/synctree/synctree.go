package synctree

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/imdario/mergo"
	uuid "github.com/kthomas/go.uuid"
)

type SyncTreeNode struct {
	SyncTreeNodeID string
	ParentNodeID   string
	Value          string
	IsLeaf         bool
	IsRoot         bool
	IsHash         bool
	IsCovered      bool
	Level          int
	Index          int
}

type BaseledgerSyncTree struct {
	RootProof string
	Nodes     []SyncTreeNode
}

func CreateFromBusinessObjectJson(businessObjectJson string, knowledgeLimiters []string) BaseledgerSyncTree {
	var result map[string]interface{}
	var leafIndex = 0
	json.Unmarshal([]byte(businessObjectJson), &result)
	//Flatten hierarchical structure into non-hierarchical leaf node structure
	FlattenOut, err := Flatten(result, nil)
	if err != nil {
		fmt.Println(err)

	}

	var LeafNodeSlice []SyncTreeNode
	//Create a leaf data strcuture for every leaf node
	for k, v := range FlattenOut {
		leaf := SyncTreeNode{}
		leaf.SyncTreeNodeID = uuid.NewV4().String()
		leaf.IsCovered = false
		leaf.IsHash = false
		leaf.IsLeaf = true
		leaf.IsRoot = false
		leaf.Value = k + ":" + fmt.Sprint(v) // this is the value
		leaf.Index = leafIndex
		leaf.Level = 0
		LeafNodeSlice = append(LeafNodeSlice, leaf)
		leafIndex++
	}

	//Each entry will be a leaf. We need 2^x leafs. Determine x
	var x = math.Max(1, math.Ceil(math.Log2(float64(len(LeafNodeSlice)))))
	//Now "fill" the dictionary up to 2^x entries
	for l := len(LeafNodeSlice); l < int(math.Pow(float64(2), x)); l++ {
		LeafNodeSlice = append(LeafNodeSlice, SyncTreeNode{SyncTreeNodeID: uuid.NewV4().String(), IsLeaf: true, Index: leafIndex})
		leafIndex++
	}

	syncTree := BaseledgerSyncTree{}
	//Now we build the tree out of the nodes. This means taking always two leafs and combining them by hashing their joint values. We do this recursively until we reached/built the root
	syncTree.Nodes = buildBONodesRecursive(LeafNodeSlice)

	//Set Root proof
	for _, rootnode := range syncTree.Nodes {
		if rootnode.IsRoot {
			syncTree.RootProof = rootnode.Value
			break
		}
	}

	//Limit knowledge
	if len(knowledgeLimiters) > 0 {
		for _, v := range syncTree.Nodes {
			if v.IsLeaf && isNodeKnowledgeLimited(v.Value, knowledgeLimiters) {
				v.IsCovered = true
				v.Value = ""
			}
		}
	}

	return syncTree
}

func GetBusinessObjectJson(syncTree BaseledgerSyncTree) string {
	jsonelements := make(map[string]interface{})
	for _, v := range syncTree.Nodes {
		if v.IsLeaf && v.Value != "" {
			var substrings []string = strings.SplitN(v.Value, ":", 2)
			if len(substrings) > 1 {
				jsonelements[substrings[0]] = substrings[1]
			}
		}
	}
	unflattenedOut, err := unflatten3(jsonelements)
	if err != nil {
		fmt.Println(err)
	}
	//Convert data structure into JSON string
	jsonstring, _ := json.Marshal(unflattenedOut)
	return string(jsonstring)
}

func VerifyHashMatch(blockchainProof string, existingBusinessObjectProof string, baseledgerSyncTreeJson string) bool {
	var ret bool = false
	//Level 0 check (Offchain Message Hash matches Blockchain stored Proof)
	ret = blockchainProof == existingBusinessObjectProof
	if ret {
		bpbo := BaseledgerSyncTree{}
		json.Unmarshal([]byte(baseledgerSyncTreeJson), &bpbo)

		//Level A check (Proofs match?)
		ret = existingBusinessObjectProof == bpbo.RootProof

		if ret {
			//Level B check (All intermediate proof calculcations match?)
			ret = verifyIntermediateHashes(bpbo)
		}
		if ret {
			var limitedKnowledgeNodes []int
			linq.From(bpbo.Nodes).Where(func(c interface{}) bool {
				return c.(SyncTreeNode).IsLeaf && c.(SyncTreeNode).IsCovered
			}).Select(func(c interface{}) interface{} {
				return c.(SyncTreeNode).Index
			}).ToSlice(&limitedKnowledgeNodes)

			if len(limitedKnowledgeNodes) > 0 {
				fmt.Println("Limited Knowledge in SyncTree! Indices:" + fmt.Sprint(limitedKnowledgeNodes))
			}

			//Level C check (Check existing (not masked) leaf nodes)
			ret = verifyLeafNodes(bpbo)
		}
	}
	return ret
}

func buildBONodesRecursive(nodes []SyncTreeNode) []SyncTreeNode {
	var ret []SyncTreeNode
	//Only one leaf left? We determined the root
	if len(nodes) <= 1 {
		return append(ret, nodes...)
	} else { //build leafs of higher level and dive into recursion
		var parentNodes []SyncTreeNode
		for i := 0; i < len(nodes); i += 2 {
			var stohash string = nodes[i].Value + "|" + nodes[i+1].Value
			//Create parent node
			parent := SyncTreeNode{}
			parent.SyncTreeNodeID = uuid.NewV4().String()
			parent.IsCovered = false
			parent.IsHash = true
			parent.IsLeaf = false
			parent.IsRoot = len(nodes) == 2
			parent.Value = createHash(stohash)
			parent.Index = i / 2
			parent.Level = nodes[i].Level + 1
			parentNodes = append(parentNodes, parent)

			nodes[i].ParentNodeID = parent.SyncTreeNodeID
			nodes[i+1].ParentNodeID = parent.SyncTreeNodeID
		}
		ret = append(ret, nodes...)
		ret = append(ret, buildBONodesRecursive(parentNodes)...)
	}
	return ret
}

func createHash(bo string) string {
	hash := md5.Sum([]byte(bo))
	return hex.EncodeToString(hash[:])
}

func verifyLeafNodes(bpbo BaseledgerSyncTree) bool {
	var leafnodes []SyncTreeNode
	var levelplusnodes []SyncTreeNode

	//using linq for go
	linq.From(bpbo.Nodes).Where(func(c interface{}) bool {
		return c.(SyncTreeNode).Level == 0
	}).OrderBy( // sort
		func(node interface{}) interface{} {
			return (node.(SyncTreeNode).Index)
		}).Select(func(c interface{}) interface{} {
		return c.(SyncTreeNode)
	}).ToSlice(&leafnodes)

	//using linq for go
	linq.From(bpbo.Nodes).Where(func(c interface{}) bool {
		return c.(SyncTreeNode).Level == 1
	}).OrderBy( // sort
		func(node interface{}) interface{} {
			return (node.(SyncTreeNode).Index)
		}).Select(func(c interface{}) interface{} {
		return c.(SyncTreeNode)
	}).ToSlice(&levelplusnodes)

	for n := 0; n < len(leafnodes); n += 2 {
		if !compareNodeHashes(leafnodes[n], leafnodes[n+1], levelplusnodes[n/2]) {
			return false
		}
	}
	return true
}

func verifyIntermediateHashes(bpbo BaseledgerSyncTree) bool {
	//Calculate all intermediate hashes
	maxlevel := 0
	for _, v := range bpbo.Nodes {
		if v.Level > maxlevel {
			maxlevel = v.Level
		}
	}

	for level := 1; level < maxlevel; level++ {
		var levelnodes []SyncTreeNode
		var levelplusnodes []SyncTreeNode

		//using linq for go
		linq.From(bpbo.Nodes).Where(func(c interface{}) bool {
			return c.(SyncTreeNode).Level == level
		}).OrderBy( // sort
			func(node interface{}) interface{} {
				return (node.(SyncTreeNode).Index)
			}).Select(func(c interface{}) interface{} {
			return c.(SyncTreeNode)
		}).ToSlice(&levelnodes)

		//using linq for go
		linq.From(bpbo.Nodes).Where(func(c interface{}) bool {
			return c.(SyncTreeNode).Level == level+1
		}).OrderBy( // sort
			func(node interface{}) interface{} {
				return (node.(SyncTreeNode).Index)
			}).Select(func(c interface{}) interface{} {
			return c.(SyncTreeNode)
		}).ToSlice(&levelplusnodes)

		for n := 0; n < len(levelnodes); n += 2 {
			if !compareNodeHashes(levelnodes[n], levelnodes[n+1], levelplusnodes[n/2]) {
				return false
			}
		}

	}
	return true
}

func compareNodeHashes(node1 SyncTreeNode, node2 SyncTreeNode, fatherNode SyncTreeNode) bool {
	ret := false

	//Create Hash
	var stohash string = node1.Value + "|" + node2.Value
	if fatherNode.Value == createHash(stohash) {
		ret = true
	}

	return ret
}

func isNodeKnowledgeLimited(nodeValue string, knowledgeLimiters []string) bool {
	if len(nodeValue) > 0 && len(knowledgeLimiters) > 0 {
		lastIndex := strings.LastIndex(nodeValue, ":")
		if lastIndex > -1 {
			nodeAttribute := nodeValue[lastIndex:]
			for _, v := range knowledgeLimiters {
				if v == nodeAttribute {
					return true
				}
			}
		}
	}
	return false
}

/// Flatten / Unflatten Methods below (can be put/packaged somewhere else if needed)
type Options struct {
	Delimiter string
	Safe      bool
	MaxDepth  int
}

// Flatten the map, it returns a map one level deep
// regardless of how nested the original map was.
// By default, the flatten has Delimiter = ".", and
// no limitation of MaxDepth
func Flatten(nested map[string]interface{}, opts *Options) (m map[string]interface{}, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}

	m, err = flatten("", 0, nested, opts)

	return
}

func flatten(prefix string, depth int, nested interface{}, opts *Options) (flatmap map[string]interface{}, err error) {
	flatmap = make(map[string]interface{})

	switch nested := nested.(type) {
	case map[string]interface{}:
		if opts.MaxDepth != 0 && depth >= opts.MaxDepth {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, map[string]interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for k, v := range nested {
			// create new key
			newKey := k
			if prefix != "" {
				newKey = prefix + opts.Delimiter + newKey
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	case []interface{}:
		if opts.Safe {
			flatmap[prefix] = nested
			return
		}
		if reflect.DeepEqual(nested, []interface{}{}) {
			flatmap[prefix] = nested
			return
		}
		for i, v := range nested {
			newKey := strconv.Itoa(i)
			if prefix != "" {
				// newKey = prefix + opts.Delimiter + newKey
				// we need to differentiate array indexes in key path
				newKey = prefix + "[" + newKey + "]"
			}
			fm1, fe := flatten(newKey, depth+1, v, opts)
			if fe != nil {
				err = fe
				return
			}
			update(flatmap, fm1)
		}
	default:
		flatmap[prefix] = nested
	}
	return
}

// update is the function that update to map with from
// example:
// to = {"hi": "there"}
// from = {"foo": "bar"}
// then, to = {"hi": "there", "foo": "bar"}
func update(to map[string]interface{}, from map[string]interface{}) {
	for kt, vt := range from {
		to[kt] = vt
	}
}

// Unflatten the map, it returns a nested map of a map
// By default, the flatten has Delimiter = "."
func Unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	if opts == nil {
		opts = &Options{
			Delimiter: ".",
		}
	}
	nested, err = unflatten(flat, opts)
	return
}

func unflatten(flat map[string]interface{}, opts *Options) (nested map[string]interface{}, err error) {
	nested = make(map[string]interface{})
	for k, v := range flat {
		temp := uf(k, v, opts).(map[string]interface{})
		err = mergo.Merge(&nested, temp)
		if err != nil {
			return
		}
	}

	return
}

func uf(k string, v interface{}, opts *Options) (n interface{}) {
	n = v

	keys := strings.Split(k, opts.Delimiter)

	for i := len(keys) - 1; i >= 0; i-- {
		temp := make(map[string]interface{})
		temp[keys[i]] = n
		n = temp
	}

	return
}

func unflatten2(flat map[string]interface{}) (map[string]interface{}, error) {
	unflat := map[string]interface{}{}

	for key, value := range flat {
		keyParts := strings.Split(key, ".")

		// Walk the keys until we get to a leaf node.
		m := unflat
		for i, k := range keyParts[:len(keyParts)-1] {
			v, exists := m[k]
			if !exists {
				newMap := map[string]interface{}{}
				m[k] = newMap
				m = newMap
				continue
			}

			innerMap, ok := v.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("key=%v is not an object", strings.Join(keyParts[0:i+1], "."))
			}
			m = innerMap
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%v already exists", key)
		}
		m[keyParts[len(keyParts)-1]] = value
	}

	return unflat, nil
}

// checks key path to see if format is array (eg foo[i]) and return key name and index i
func isArrayKeyPart(keyPart string) (bool, string, int) {
	lastChar := keyPart[len(keyPart)-1:]
	if lastChar != "]" {
		return false, "", -1
	}
	indexArray := keyPart[len(keyPart)-2 : len(keyPart)-1]
	arrayNum, _ := strconv.Atoi(indexArray)
	keyName := keyPart[:len(keyPart)-3]
	return true, keyName, arrayNum
}

func unflatten3(flat map[string]interface{}) (map[string]interface{}, error) {
	unflat := map[string]interface{}{}
	for key, value := range flat {
		keyParts := strings.Split(key, ".")
		// Walk the keys until we get to a leaf node.
		m := unflat
		for i, k := range keyParts[:len(keyParts)-1] {
			kName := k
			isArr, keyName, index := isArrayKeyPart(k)
			if isArr {
				kName = keyName
			}
			v, exists := m[kName]
			if !exists {
				newMap := map[string]interface{}{}
				// same as unflatten2, this means that in next iteration
				// newMap (empty map) will be used, but unflat[k] will be updated because of same ref
				if !isArr {
					m[k] = newMap
					m = newMap
					continue
				}

				// if its array key, we assing array full of empty maps, but keep reference to map at current index
				// which will be updated after
				var arr []interface{}
				for i = 0; i <= index; i++ {
					if i == index {
						arr = append(arr, newMap)
					} else {
						arr = append(arr, map[string]interface{}{})
					}
				}
				m[keyName] = arr
				m = newMap
				continue
			}

			if isArr {
				innerArr, _ := v.([]interface{})
				if index >= len(innerArr) {
					newMap := map[string]interface{}{}

					var arr []interface{}
					// making sure that only map at index is stored at m
					for i = 0; i <= index; i++ {
						if i < len(innerArr) {
							arr = append(arr, innerArr[i])
						} else {
							if i == index {
								arr = append(arr, newMap)
							} else {
								arr = append(arr, map[string]interface{}{})
							}
						}
					}
					m[keyName] = arr
					m = newMap

				} else {
					m = innerArr[index].(map[string]interface{})
				}
			} else {
				innerMap, ok := v.(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("key=%v is not an object", strings.Join(keyParts[0:i+1], "."))
				}
				m = innerMap
			}
		}

		leafKey := keyParts[len(keyParts)-1]
		if _, exists := m[leafKey]; exists {
			return nil, fmt.Errorf("key=%v already exists %v", leafKey, m)
		}
		m[leafKey] = value
	}

	return unflat, nil
}
