import { ThemeProvider } from '@siafoundation/design-system'
import { Redirect, Route, Switch } from 'react-router-dom'
import { Layout } from './components/Layout'
import { Home } from './pages/Home'
import { CreateSwap } from './pages/CreateSwap'
import { InputSwap } from './pages/InputSwap'
import { StepSwap } from './pages/StepSwap'
import { routes } from './routes'
import { SwapProvider } from './hooks/useSwap'

export function App() {
  return (
    <ThemeProvider>
      <SwapProvider>
        <Layout>
          <Switch>
            <Route path={routes.home} exact component={Home} />
            <Route path={routes.create} component={CreateSwap} />
            <Route path={routes.input} component={InputSwap} />
            <Route path={routes.step} component={StepSwap} />
            <Redirect from="*" to="/" />
          </Switch>
        </Layout>
      </SwapProvider>
    </ThemeProvider>
  )
}

export default App
