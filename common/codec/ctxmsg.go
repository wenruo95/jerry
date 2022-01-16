/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : ctxmsg.go
*   coder: zemanzeng
*   date : 2022-01-16 11:22:28
*   desc : ctx msg
*
================================================================*/

package codec

import "github.com/wenruo95/jerry/jlog"

type CtxMsg interface {
	Logger() jlog.Logger
}
