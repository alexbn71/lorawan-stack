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

import React from 'react'

import { Route, Switch } from 'react-router-dom'
import bind from 'autobind-decorator'

import IntlHelmet from '../../../lib/components/intl-helmet'
import { withEnv, EnvProvider } from '../../../lib/components/env'
import ErrorView from '../../../lib/components/error-view'
import { pageDataSelector } from '../../../lib/selectors/env'
import SideNavigation from '../../../components/navigation/side'
import Header from '../../containers/header'
import Footer from '../../../components/footer'
import Landing from '../landing'
import Login from '../login'
import FullViewError from '../error'

import style from './app.styl'

@withEnv
@bind
export default class ConsoleApp extends React.Component {
  render () {
    const {
      env,
    } = this.props

    const pageData = pageDataSelector(env)

    if (pageData && pageData.error) {
      return <FullViewError error={pageData.error} />
    }

    return (
      <EnvProvider env={env}>
        <ErrorView ErrorComponent={FullViewError}>
          <div className={style.app}>
            <IntlHelmet
              titleTemplate={`%s - ${env.site_title ? `${env.site_title} - ` : ''}${env.site_name}`}
              defaultTitle={env.site_name}
            />
            <div id="modal-container" />
            <Header className={style.header} />
            <main className={style.main}>
              <div>
                <SideNavigation />
              </div>
              <div className={style.content}>
                <Switch>
                  {/* routes for registration, privacy policy, other public pages */}
                  <Route path={`${env.app_root}/login`} component={Login} />
                  <Route path={env.app_root} component={Landing} />
                </Switch>
              </div>
            </main>
            <Footer className={style.footer} />
          </div>
        </ErrorView>
      </EnvProvider>
    )
  }
}
