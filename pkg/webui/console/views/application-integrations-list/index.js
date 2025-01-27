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

import React, { Component } from 'react'
import { Container, Row, Col } from 'react-grid-system'

import WebhooksTable from '../../containers/webhooks-table'

import IntlHelmet from '../../../lib/components/intl-helmet'
import sharedMessages from '../../../lib/shared-messages'

import PAGE_SIZES from '../../constants/page-sizes'

export default class ApplicationIntegrationsList extends Component {
  render () {
    const { appId } = this.props.match.params

    return (
      <Container>
        <Row>
          <IntlHelmet title={sharedMessages.integrations} />
          <Col sm={12}>
            <WebhooksTable
              pageSize={PAGE_SIZES.REGULAR}
              appId={appId}
            />
          </Col>
        </Row>
      </Container>
    )
  }
}
