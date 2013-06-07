package recordio

import (
	"os"
	"testing"

	"github.com/arunjitsingh/gorecordio/recordio"
)

const (
	testFileName = "-test.rec"
)

func Test_WriteAndRead(t *testing.T) {
	defer os.Remove(testFileName)

	file, _ := os.Create(testFileName)
	recw := recordio.NewRecordWriter(file)
	n, err := recw.Write([]byte("Hello!"))
	if err != nil {
		t.Error("didn't work", n, err)
	}
	recw.Close()

	file, _ = os.Open(testFileName)
	recr := recordio.NewRecordReader(file)
	d, err := recr.ReadNext()
	if err != nil || string(d) != "Hello!" {
		t.Error("didn't work", d, err)
	}
	recr.Close()

}

func Test_WriteAndRead_Repeated(t *testing.T) {
	defer os.Remove(testFileName)

	values := []string{"One", "Two"}

	file, _ := os.Create(testFileName)
	recw := recordio.NewRecordWriter(file)
	for _, value := range values {
		n, err := recw.Write([]byte(value))
		if err != nil {
			t.Error("didn't work", n, err)
		}
	}
	recw.Close()

	file, _ = os.Open(testFileName)
	recr := recordio.NewRecordReader(file)
	for _, value := range values {
		d, err := recr.ReadNext()
		if err != nil || string(d) != value {
			t.Error("didn't work", d, err)
		}
	}
	recr.Close()
}
