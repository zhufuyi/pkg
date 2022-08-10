package errcode

import (
	"fmt"
	"testing"
)

func TestToGRPCError(t *testing.T) {
	err := ToGRPCError(InvalidParam)
	fmt.Println(err)

}

func TestToGRPCStatus(t *testing.T) {
	status := ToGRPCStatus(InvalidParam, "invoker listUser error")
	fmt.Println(status.Err())
	fmt.Println(status.Code())
}
