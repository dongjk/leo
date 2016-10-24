package parser


import (
    "testing"
    "fmt"
    "os"
)

func TestDo(t *testing.T) {
    f, err := os.Open("testdata/testfile1")
    if err!=nil{
        fmt.Println(err)
    }
    n,err:=os.Stat("testdata/testfile1")
    if err!=nil{
        fmt.Println(err)
    }
    b:=make([]byte,n.Size())
    f.Read(b)
    fmt.Println(do(b,""))
}

func TestGetStartTime(t *testing.T) {
    fmt.Println(getStartTime("Script started on Mon, Oct  4, 2016 1:07:44 AM"))
}