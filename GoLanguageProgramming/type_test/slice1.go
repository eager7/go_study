package main

import (
    "fmt"
)

func main(){
    var myArray [10] int = [10] int {1,2,3,4,5,6,7,8,9,10}

    var mySlice []int = myArray[:5]
    fmt.Println("Elements of myArray:")
    for _,v := range myArray {
        println(v,"")
    }

    println("Elements of mySlice:")
    for _,v := range mySlice{
        println(v)
    }
    println("")
}
