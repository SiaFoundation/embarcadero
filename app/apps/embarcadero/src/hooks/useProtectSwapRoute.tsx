import { useSwap } from '../contexts/swap'
import { useHistory } from 'react-router-dom'
import { routes } from '../routes'

export function useProtectSwapRoute() {
  const { raw, isValidating } = useSwap()
  const history = useHistory()

  if (!isValidating && !raw) {
    history.push(routes.home)
  }
}
