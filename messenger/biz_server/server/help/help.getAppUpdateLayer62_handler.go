/*
 *  Copyright (c) 2018-present,  NebulaChat Studio (https://nebula.chat).
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package help

import (
    "github.com/golang/glog"
    "golang.org/x/net/context"
    "github.com/libregram/chatengine/pkg/grpc_util"
    "github.com/libregram/chatengine/pkg/logger"
    "github.com/libregram/chatengine/mtproto"
)

// help.getAppUpdate#ae2de196 = help.AppUpdate;
func (s *HelpServiceImpl) HelpGetAppUpdateLayer62(ctx context.Context, request *mtproto.TLHelpGetAppUpdateLayer62) (*mtproto.Help_AppUpdate, error) {
    md := grpc_util.RpcMetadataFromIncoming(ctx)
    glog.Infof("help.getAppUpdate#ae2de196 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

    // TODO(@benqi): Impl HelpGetAppUpdate logic
    reply := &mtproto.TLHelpNoAppUpdate{Data2: &mtproto.Help_AppUpdate_Data{}}

    glog.Infof("help.getAppUpdate#ae2de196 - reply: %s\n", logger.JsonDebugData(reply))
    return reply.To_Help_AppUpdate(), nil
}
