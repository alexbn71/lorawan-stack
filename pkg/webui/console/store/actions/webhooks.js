// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { createRequestActions } from './lib'

export const GET_WEBHOOKS_LIST_BASE = 'GET_WEBHOOKS_LIST'
export const [{
  request: GET_WEBHOOKS_LIST,
  success: GET_WEBHOOKS_LIST_SUCCESS,
  failure: GET_WEBHOOKS_LIST_FAILURE,
}, {
  request: getWebhooksList,
  success: getWebhooksListSuccess,
  failure: getWebhooksListFailure,
}] = createRequestActions(GET_WEBHOOKS_LIST_BASE,
  appId => ({ appId })
)
