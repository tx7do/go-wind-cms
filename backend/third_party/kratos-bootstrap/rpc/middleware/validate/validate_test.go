package validate

import (
	"google.golang.org/protobuf/proto"
)

type testcase struct {
	name string
	req  proto.Message
	err  bool
}
