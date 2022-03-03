import { ThemeProvider } from '@siafoundation/design-system'
import { Redirect, Route, Switch } from 'react-router-dom'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { routes } from './routes'
import { SwapProvider } from './contexts/swap'
import { DialogProvider } from './contexts/dialog'
import { CreateNewSwap } from './pages/CreateNewSwap'
import { LoadExistingSwap } from './pages/LoadExistingSwap'
import { SwapStep } from './pages/SwapStep'

export function App() {
  return (
    <ThemeProvider>
      <SwapProvider>
        <DialogProvider>
          <Layout>
            <Switch>
              <Route path={routes.home} exact component={Home} />
              <Route path={routes.create} component={CreateNewSwap} />
              <Route path={routes.input} component={LoadExistingSwap} />
              <Route path={routes.swap} component={SwapStep} />
              <Redirect from="*" to="/" />
            </Switch>
          </Layout>
        </DialogProvider>
      </SwapProvider>
    </ThemeProvider>
  )
}

export default App
