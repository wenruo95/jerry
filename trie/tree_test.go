/*================================================================
*   Copyright (C) 2021. All rights reserved.
*
*   file : tree_test.go
*   coder: zemanzeng
*   date : 2021-06-17 20:49:55
*   desc : 前缀树测试用例
*
================================================================*/

package trie

import (
	"encoding/json"
	"math/rand"
	"sort"
	"strings"
	"testing"
	"time"
)

func Test_TreeInsert(t *testing.T) {
	type InsertResult struct {
		str    string
		succ   bool
		length int
		list   []string
	}

	str2result := []*InsertResult{
		&InsertResult{"abcdefg", true, 1, []string{}},
		&InsertResult{"abcred", true, 2, []string{}},
		&InsertResult{"abcded", true, 3, []string{}},
		&InsertResult{"abcd", true, 2, []string{"abcdefg", "abcded"}},
		&InsertResult{"abacds", true, 3, []string{}},
		&InsertResult{"adasfds", true, 4, []string{}},
		&InsertResult{"g2341", true, 5, []string{}},
		&InsertResult{"中华小当家", true, 6, []string{}},
		&InsertResult{"中华", true, 6, []string{"中华小当家"}},
		&InsertResult{"abcefds", true, 7, []string{}},
		&InsertResult{"abcdfds", false, 7, []string{"abcd"}},
		&InsertResult{"中华小当家2", false, 7, []string{"中华"}},
	}

	result := []string{
		"abcred",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华",
		"abcefds",
	}

	tree := NewTree()
	for index, result := range str2result {
		str := result.str
		succ, list := tree.Insert([]rune(str))
		if succ && !tree.Exist([]rune(str)) {
			t.Errorf("insert succ but not exist in tree. index:%v str:%v", index, str)
		}
		if result.succ != succ || result.length != tree.Len() || !listEqual(result.list, list) {
			t.Errorf("insert result error. index:%v str:%v result:%+v succ:%v list:%+v len:%v",
				index, str, result, succ, list, tree.Len())
		}
	}

	if !listEqual(tree.Root.List(), result) {
		t.Errorf("tree result list not equal. tree:%+v result:%+v", tree.Root.List(), result)
	}
}

func Test_TreeDelete(t *testing.T) {
	list := []string{
		"abcdefg",
		"abcred",
		"abcded",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华小当家",
		"中华",
		"abcefds",
		"abcdfds",
		"中华小当家2",
	}
	tree := NewTree()
	for _, str := range list {
		tree.Insert([]rune(str))
	}

	notExist := []string{
		"abcdefg",
		"abcded",
		"中华小当家",
		"abcdfds",
		"中华小当家2",
		"hello world",
	}
	for index, str := range notExist {
		if tree.Delete([]rune(str)) {
			t.Errorf("delete a not exist item index:%v str:%v", index, str)
		}
	}

	result := []string{
		"abcred",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华",
		"abcefds",
	}
	for index, str := range result {
		if !tree.Delete([]rune(str)) {
			t.Errorf("delete item failed index:%v str:%v", index, str)
		}
		if tree.Len() != len(result)-index-1 {
			t.Errorf("item length not matched expect:%v tree:%v index:%v str:%v",
				len(result)-index-1, tree.Len(), index, str)
		}
		if !listEqual(tree.Root.List(), result[index+1:]) {
			t.Errorf("tree list not equal expect:%v tree:%v index:%v str:%v",
				result[index+1:], tree.Root.List(), index, str)

		}
	}

}

func Test_TreeMatch(t *testing.T) {
	list := []string{
		"崩殂",
		"存亡",
		"存亡之际",
		"侍",
		"abcdefg",
		"abcred",
		"abcded",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华小当家",
		"中华",
		"abcefds",
		"abcdfds",
		"中华小当家2",
	}
	tree := NewTree()
	for _, str := range list {
		tree.Insert([]rune(str))
	}

	type MatchResult struct {
		Words     string
		SubString []string
	}

	results := []*MatchResult{
		&MatchResult{
			Words:     `先帝创业未半而中道崩殂，今天下三分，益州疲弊，此诚危急存亡之秋也。然侍卫之臣不懈于内，忠志之士忘身于外者，盖追先帝之殊遇，欲报之于陛下也。诚宜开张圣听，以光先帝遗德，恢弘志士之气，不宜妄自菲薄，引喻失义，以塞忠谏之路也。宫中府中，俱为一体，陟罚臧否，不宜异同。若有作奸犯科及为忠善者，宜付有司论其刑赏，以昭陛下平明之理，不宜偏私，使内外异法也。侍中、侍郎郭攸之、费祎、董允等，此皆良实，志虑忠纯，是以先帝简拔以遗陛下。愚以为宫中之事，事无大小，悉以咨之，然后施行，必能裨补阙漏，有所广益。将军向宠，性行淑均，晓畅军事，试用于昔日，先帝称之曰能，是以众议举宠为督。愚以为营中之事，悉以咨之，必能使行阵和睦，优劣得所。亲贤臣，远小人，此先汉所以兴隆也；亲小人，远贤臣，此后汉所以倾颓也。先帝在时，每与臣论此事，未尝不叹息痛恨于桓、灵也。侍中、尚书、长史、参军，此悉贞良死节之臣，愿陛下亲之信之，则汉室之隆，可计日而待也。臣本布衣，躬耕于南阳，苟全性命于乱世，不求闻达于诸侯。先帝不以臣卑鄙，猥自枉屈，三顾臣于草庐之中，咨臣以当世之事，由是感激，遂许先帝以驱驰。后值倾覆，受任于败军之际，奉命于危难之间，尔来二十有一年矣。先帝知臣谨慎，故临崩寄臣以大事也。受命以来，夙夜忧叹，恐托付不效，以伤先帝之明，故五月渡泸，深入不毛。今南方已定，兵甲已足，当奖率三军，北定中原，庶竭驽钝，攘除奸凶，兴复汉室，还于旧都。此臣所以报先帝而忠陛下之职分也。至于斟酌损益，进尽忠言，则攸之、祎、允之任也。愿陛下托臣以讨贼兴复之效，不效，则治臣之罪，以告先帝之灵。若无兴德之言，则责攸之、祎、允等之慢，以彰其咎；陛下亦宜自谋，以咨诹善道，察纳雅言，深追先帝遗诏，臣不胜受恩感激。今当远离，临表涕零，不知所言。`,
			SubString: []string{"崩殂", "存亡", "侍", "侍", "侍", "侍"},
		},

		&MatchResult{
			Words:     "abcdefghg2341ijklmnopqrstu《a中华小aba\ncdshello》vwxyz",
			SubString: []string{"abcd", "g2341", "中华", "abacds"},
		},
		&MatchResult{
			Words:     "abc",
			SubString: []string{},
		},
		&MatchResult{
			Words:     "a中b华",
			SubString: []string{},
		},
		&MatchResult{
			Words:     "中\t\r\n华小当家",
			SubString: []string{"中华"},
		},
		&MatchResult{
			Words:     "中华",
			SubString: []string{"中华"},
		},
	}

	igList := []rune{'\t', '\n', '\r', ' '}
	ignores := make(map[rune]struct{})
	for _, ig := range igList {
		ignores[ig] = struct{}{}
	}

	for index, result := range results {
		subs := tree.Match([]rune(result.Words), ignores)
		if !listEqual(result.SubString, subs) {
			t.Errorf("index:%v not matched expect:%+v result:%+v", index, result.SubString, subs)
		}
	}

}

func Test_TreeIndexs(t *testing.T) {
	list := []string{
		"崩殂",
		"存亡",
		"存亡之际",
		"侍",
		"abcdefg",
		"abcred",
		"abcded",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华小当家",
		"中华",
		"abcefds",
		"abcdfds",
		"中华小当家2",
	}
	tree := NewTree()
	for _, str := range list {
		tree.Insert([]rune(str))
	}

	type IndexsResult struct {
		Words  string
		Indexs []*RuneIndex
	}

	results := []*IndexsResult{
		&IndexsResult{
			Words: `先帝创业未半而中道崩殂，今天下三分，益州疲弊，此诚危急存亡之秋也。然侍卫之臣不懈于内，忠志之士忘身于外者，盖追先帝之殊遇，欲报之于陛下也。诚宜开张圣听，以光先帝遗德，恢弘志士之气，不宜妄自菲薄，引喻失义，以塞忠谏之路也。宫中府中，俱为一体，陟罚臧否，不宜异同。若有作奸犯科及为忠善者，宜付有司论其刑赏，以昭陛下平明之理，不宜偏私，使内外异法也。侍中、侍郎郭攸之、费祎、董允等，此皆良实，志虑忠纯，是以先帝简拔以遗陛下。愚以为宫中之事，事无大小，悉以咨之，然后施行，必能裨补阙漏，有所广益。将军向宠，性行淑均，晓畅军事，试用于昔日，先帝称之曰能，是以众议举宠为督。愚以为营中之事，悉以咨之，必能使行阵和睦，优劣得所。亲贤臣，远小人，此先汉所以兴隆也；亲小人，远贤臣，此后汉所以倾颓也。先帝在时，每与臣论此事，未尝不叹息痛恨于桓、灵也。侍中、尚书、长史、参军，此悉贞良死节之臣，愿陛下亲之信之，则汉室之隆，可计日而待也。臣本布衣，躬耕于南阳，苟全性命于乱世，不求闻达于诸侯。先帝不以臣卑鄙，猥自枉屈，三顾臣于草庐之中，咨臣以当世之事，由是感激，遂许先帝以驱驰。后值倾覆，受任于败军之际，奉命于危难之间，尔来二十有一年矣。先帝知臣谨慎，故临崩寄臣以大事也。受命以来，夙夜忧叹，恐托付不效，以伤先帝之明，故五月渡泸，深入不毛。今南方已定，兵甲已足，当奖率三军，北定中原，庶竭驽钝，攘除奸凶，兴复汉室，还于旧都。此臣所以报先帝而忠陛下之职分也。至于斟酌损益，进尽忠言，则攸之、祎、允之任也。愿陛下托臣以讨贼兴复之效，不效，则治臣之罪，以告先帝之灵。若无兴德之言，则责攸之、祎、允等之慢，以彰其咎；陛下亦宜自谋，以咨诹善道，察纳雅言，深追先帝遗诏，臣不胜受恩感激。今当远离，临表涕零，不知所言。`,
			Indexs: []*RuneIndex{
				&RuneIndex{Index: 9, Length: 2, Words: "崩殂"},
				&RuneIndex{Index: 27, Length: 2, Words: "存亡"},
				&RuneIndex{Index: 34, Length: 1, Words: "侍"},
				&RuneIndex{Index: 172, Length: 1, Words: "侍"},
				&RuneIndex{Index: 175, Length: 1, Words: "侍"},
				&RuneIndex{Index: 366, Length: 1, Words: "侍"},
			},
		},

		&IndexsResult{
			Words: "abcdefghg2341ijklmnopqrstu《a中华小aba\ncdshello》vwxyz",
			Indexs: []*RuneIndex{
				&RuneIndex{Index: 0, Length: 4, Words: "abcd"},
				&RuneIndex{Index: 8, Length: 5, Words: "g2341"},
				&RuneIndex{Index: 28, Length: 2, Words: "中华"},
				&RuneIndex{Index: 31, Length: 7, Words: "abacds"},
			},
		},
		&IndexsResult{
			Words:  "abc",
			Indexs: []*RuneIndex{},
		},
		&IndexsResult{
			Words:  "a中b华",
			Indexs: []*RuneIndex{},
		},
		&IndexsResult{
			Words: "中\t\r\n 华小当家",
			Indexs: []*RuneIndex{
				&RuneIndex{Index: 0, Length: 6, Words: "中华"},
			},
		},
		&IndexsResult{
			Words: "中华",
			Indexs: []*RuneIndex{
				&RuneIndex{Index: 0, Length: 2, Words: "中华"},
			},
		},
	}

	igList := []rune{'\t', '\n', '\r', ' '}
	ignores := make(map[rune]struct{})
	for _, ig := range igList {
		ignores[ig] = struct{}{}
	}

	for index, result := range results {
		indexs := tree.Indexs([]rune(result.Words), ignores)
		if !runIndexsEqual(result.Indexs, indexs, ignores) {
			t.Errorf("index:%v not matched expect:%+v result:%+v runes:%v", index, toStr(result.Indexs), toStr(indexs), []rune(result.Words))
		}
	}

}

func listEqual(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	var index int
	for index = 0; index < len(a) && a[index] == b[index]; index++ {
	}
	return index == len(a)
}

func runIndexsEqual(result []*RuneIndex, indexs []*RuneIndex, ignores map[rune]struct{}) bool {
	if len(result) != len(indexs) {
		return false
	}
	for i := 0; i < len(result); i++ {
		if result[i].Index != indexs[i].Index || result[i].Length != indexs[i].Length {
			return false
		}

		var s string = indexs[i].Words
		for ig := range ignores {
			s = strings.Replace(s, string(ig), "", -1)
		}
		if s != result[i].Words {
			return false
		}
	}
	return true
}

func toStr(i interface{}) string {
	if buff, err := json.Marshal(i); err == nil {
		return string(buff)
	}
	return ""
}

func Benchmark_Match(b *testing.B) {
	list := []string{
		"abcdefg",
		"abcred",
		"abcded",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华小当家",
		"中华",
		"abcefds",
		"abcdfds",
		"中华小当家2",
	}
	tree := NewTree()
	for _, str := range list {
		tree.Insert([]rune(str))
	}

	rand.Seed(time.Now().UnixNano())
	words := []rune(randStr(100000))
	for i := 0; i < b.N; i++ {
		tree.Match(words, nil)
	}

}

func Benchmark_Indexs(b *testing.B) {
	list := []string{
		"abcdefg",
		"abcred",
		"abcded",
		"abcd",
		"abacds",
		"adasfds",
		"g2341",
		"中华小当家",
		"中华",
		"abcefds",
		"abcdfds",
		"中华小当家2",
	}
	tree := NewTree()
	for _, str := range list {
		tree.Insert([]rune(str))
	}

	rand.Seed(time.Now().UnixNano())
	words := []rune(randStr(100000))
	for i := 0; i < b.N; i++ {
		tree.Indexs(words, nil)
	}

}

func randStr(count int) string {
	var b strings.Builder
	b.Grow(count)

	rs := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	for i := 0; i < count; i++ {
		index := rand.Int() % len(rs)
		if _, err := b.WriteRune(rs[index]); err != nil {
			break
		}
	}
	return b.String()
}
