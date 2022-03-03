import { ThemeProvider } from '@siafoundation/design-system'
import { Redirect, Route, Switch } from 'react-router-dom'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { routes } from './routes'
import { SwapProvider } from './contexts/swap'
import { DialogProvider } from './contexts/dialog'
import { ReviewAccept } from './pages/ReviewAccept'
import { ReviewFinish } from './pages/ReviewFinish'
import { WaitingAccept } from './pages/WaitingAccept'
import { WaitingFinish } from './pages/WaitingFinish'
import { CreateNewSwap } from './pages/CreateNewSwap'
import { LoadExistingSwap } from './pages/LoadExistingSwap'

export function App() {
  return (
    <ThemeProvider>
      <SwapProvider>
        <DialogProvider>
          <Layout>
            <Switch>
              <Route path={routes.home} exact component={Home} />
              <Route path={routes.creatingANewSwap} component={CreateNewSwap} />
              <Route
                path={routes.loadingAnExistingSwap}
                component={LoadExistingSwap}
              />
              <Route
                path={routes.waitingForYouToAccept}
                component={ReviewAccept}
              />
              <Route
                path={routes.waitingForYouToFinish}
                component={ReviewFinish}
              />
              <Route
                path={routes.waitingForCounterpartyToAccept}
                component={WaitingAccept}
              />
              <Route
                path={routes.waitingForCounterpartyToFinish}
                component={WaitingFinish}
              />
              <Redirect from="*" to="/" />
            </Switch>
          </Layout>
        </DialogProvider>
      </SwapProvider>
    </ThemeProvider>
  )
}

export default App
