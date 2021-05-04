package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/minio/simdjson-go"
)

func printKey(iter simdjson.Iter, key string) (err error) {

	obj, tmp, elem := &simdjson.Object{}, &simdjson.Iter{}, simdjson.Element{}
	var count int64 = 0
	for {
		typ := iter.Advance()
		fmt.Println(typ)

		switch typ {
		case simdjson.TypeRoot:
			if typ, tmp, err = iter.Root(tmp); err != nil {
				fmt.Println("Failed at 1st stage")
				return
			}

			if typ == simdjson.TypeObject {
				if obj, err = tmp.Object(obj); err != nil {
					fmt.Println("Failed at 2nd stage")
					return
				}

				e := obj.FindKey(key, &elem)
				if e != nil && elem.Type == simdjson.TypeString {
					v, _ := elem.Iter.StringBytes()
					fmt.Println(string(v))
				}
			}

		default:
			count += 1
			fmt.Println("Processed", count)
		}
	}
}

func main() {
	if !simdjson.SupportedCPU() {
		log.Fatal("Unsupported CPU")
	}
	msg, err := ioutil.ReadFile("/Users/godfreybafana/data/golang/full.json")
	if err != nil {
		log.Fatalf("Failed to load file: %v", err)
	}

	parsed, err := simdjson.Parse(msg, nil)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	printKey(parsed.Iter(), "recipe")
}
