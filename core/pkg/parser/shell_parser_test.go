package parser


import (
    "testing"
    "fmt"
    // "os"
)

// func TestDo(t *testing.T) {
//     f, err := os.Open("testdata/testfile1")
//     if err!=nil{
//         fmt.Println(err)
//     }
//     n,err:=os.Stat("testdata/testfile1")
//     if err!=nil{
//         fmt.Println(err)
//     }
//     b:=make([]byte,n.Size())
//     f.Read(b)
//     fmt.Println(do(b,""))
// }

func TestGetStartTime(t *testing.T) {
    fmt.Println(getStartTime("Script started on Mon, Oct  4, 2016 1:07:44 AM"))
}

func TestGetTimeOffset(t *testing.T){
    a:=getTimeOffset("0.005000 157")
    if a!=5000{
        t.Fail()
    }
    b:=getTimeOffset("113.819074 28")
    if b!=113819074{
        t.Fail()
    }
}
func TestGetCharOffset(t *testing.T){
    a:=getCharOffset("0.005000 157")
     if a!=157{
        t.Fail()
    }
    b:=getCharOffset("113.819074 28")
     if b!=28{
        t.Fail()
    }
    
}