import { useSwap } from '../contexts/swap'
import { useHistory, useLocation } from 'react-router-dom'
import { routes } from '../routes'

export function useRouteToStep() {
  const { isValidating, status } = useSwap()
  const location = useLocation()
  const history = useHistory()

  if (isValidating) {
    return
  }

  const route = status && routes[status]

  if (route && !location.pathname.includes(route)) {
    console.log(status, route)
    history.push(route)
  }
}
