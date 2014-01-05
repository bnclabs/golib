package strings

import (
	"fmt"
	"testing"
)

func TestFormatMultiColumn(t *testing.T) {
	words := []string{
		"is", "formed", "by", "the", "participant", "nodes.", "Each", "node",
		"is", "identified", "by", "a", "number", "or",
		"node", "ID.", "The", "node", "ID", "serves", "not", "only", "as",
		"identification,", "but", "the", "Kademlia",
		"algorithm", "uses", "the", "node", "ID", "to", "locate", "values",
		"(usually", "file", "hashes", "or",
		"keywords).", "In", "fact,", "the", "node", "ID", "provides", "a",
		"direct", "map", "to", "file", "hashes", "and",
		"that", "node", "stores", "information", "on", "where", "to", "obtain",
		"the", "file", "or", "resource.",
	}

	lines := FormatMultiColumn(words, 100)
	for _, line := range lines {
		fmt.Printf("%s\n", line)
	}
}

//func BenchmarkFormatMultiColumn(b *testing.B) {
//    words := []string{
//        "is", "formed", "by", "the", "participant", "nodes.", "Each", "node",
//        "is", "identified", "by", "a", "number", "or",
//        "node", "ID.", "The", "node", "ID", "serves", "not", "only", "as",
//        "identification,", "but", "the", "Kademlia",
//        "algorithm", "uses", "the", "node", "ID", "to", "locate", "values",
//        "(usually", "file", "hashes", "or",
//        "keywords).", "In", "fact,", "the", "node", "ID", "provides", "a",
//        "direct", "map", "to", "file", "hashes", "and",
//        "that", "node", "stores", "information", "on", "where", "to", "obtain",
//        "the", "file", "or", "resource.",
//    }
//
//    b.ResetTimer()
//    for i := 0; i < b.N; i++ {
//        lines := FormatMultiColumn(words, 100)
//    }
//}
