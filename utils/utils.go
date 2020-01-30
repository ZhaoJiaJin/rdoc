package utils

import(
    "math/rand"
    "time"
)

func init(){
    rand.Seed(time.Now().UnixNano())
}

// Randint generate a random integer
func Randint()int{
    return rand.Int()
}
