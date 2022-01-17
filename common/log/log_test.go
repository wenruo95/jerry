/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : log_test.go
*   coder: zemanzeng
*   date : 2022-01-17 11:53:44
*   desc : test case
*
================================================================*/

package log

import (
	"context"
	"testing"
)

func TestAvailable(t *testing.T) {
	ctx := WithContextFields(context.Background(), "hello", "world", "test1")
	TraceContext(ctx, "test trace")
	ErrorContext(ctx, "test error")
	DebugContext(ctx, "test debug")

	ctx2 := WithContextFields(ctx, "hello2", "world2", "test2")
	TraceContext(ctx2, "test2 trace")
	ErrorContext(ctx2, "test2 error")
	DebugContext(ctx2, "test2 debug")

	Get(DefaultTag).Info("hello default info")

}
