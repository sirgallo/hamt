package hamttests

import "testing"
import "fmt"

import "github.com/sirgallo/hamt/pkg/hamt"

func TestHamtOperations(t *testing.T) {
	hamt := hamt.NewHAMT()

	hamt.Insert("hello", "world")
	hamt.Insert("new", "wow!")
	hamt.Insert("again", "test!")
	hamt.Insert("woah", "random entry")
	hamt.Insert("key", "Saturday!")
	hamt.Insert("sup", "6")
	hamt.Insert("final", "the!")
	hamt.Insert("6", "wow!")
	hamt.Insert("asdfasdf", "add 10")
	hamt.Insert("asdfasdf", "123123")
	hamt.Insert("asd", "queue!")
	hamt.Insert("fasdf", "interesting")
	hamt.Insert("yup", "random again!")
	hamt.Insert("asdf", "hello")
	hamt.Insert("asdffasd", "uh oh!")
	hamt.Insert("fasdfasdfasdfasdf", "error message")
	hamt.Insert("fasdfasdf", "info!")
	hamt.Insert("woah", "done")

	fmt.Printf("bitmap of root after inserts: %032b\n", hamt.Root.BitMap)

	fmt.Println("hamt after inserts")
	hamt.PrintChildren()
	
	expectedBitMap := uint32(3975939716)
	t.Log("actual bitmap:", hamt.Root.BitMap, "expected bitmap:", expectedBitMap)
	if expectedBitMap != hamt.Root.BitMap {
		t.Error("actual bitmap does not match expected bitmap")
	}

	fmt.Println("retrieve values")

	val1 := hamt.Retrieve("hello")
	expVal1 :=  "world"
	t.Logf("actual: %s, expected: %s", val1, expVal1)
	if val1 != expVal1 {
		t.Error("val 1 does not match expected val 1")
	}

	val2 := hamt.Retrieve("new")
	expVal2 :=  "wow!"
	t.Logf("actual: %s, expected: %s", val2, expVal2)
	if val2 != expVal2 {
		t.Error("val 2 does not match expected val 2")
	}

	val3 := hamt.Retrieve("asdf")
	expVal3 := "hello"
	t.Logf("actual: %s, expected: %s", val3, expVal3)
	if val3 != expVal3 {
		t.Error("val 3 does not match expected val 3")
	}

	val4 := hamt.Retrieve("asdfasdf")
	expVal4 := "add 10"
	t.Logf("actual: %s, expected: %s", val4, expVal4)
	if val4 != expVal4 {
		t.Error("val 4 does not match expected val 4")
	}

	hamt.Delete("hello")
	hamt.Delete("yup")
	hamt.Delete("asdf")
	hamt.Delete("asdfasdf")
	hamt.Delete("new")
	hamt.Delete("6")

	fmt.Printf("bitmap of root after deletes: %032b\n", hamt.Root.BitMap)

	fmt.Println("hamt after deletes")
	hamt.PrintChildren()

	expectedRootBitmap := uint32(1614546564)
	t.Log("actual bitmap:", hamt.Root.BitMap, "expected bitmap:", expectedRootBitmap)
	if expectedRootBitmap != hamt.Root.BitMap {
		t.Error("actual bitmap does not match expected bitmap")
	}
}