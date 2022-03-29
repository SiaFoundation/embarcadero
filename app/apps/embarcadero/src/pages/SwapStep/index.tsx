import { ReviewAccept } from './ReviewAccept'
import { ReviewFinish } from './ReviewFinish'
import { WaitingAccept } from './WaitingAccept'
import { WaitingFinish } from './WaitingFinish'
import { TxnConfirmed } from './TxnConfirmed'
import { SwapStatusRemote } from '../../lib/swapStatus'
import { useSwap } from '../../contexts/swap'
import { Redirect } from 'react-router-dom'
import { swapStatusToRoute, routes } from '../../routes'
import { usePathParams } from '../../hooks/usePathParams'
import { TxnPending } from './TxnPending'

const componentMap: Record<SwapStatusRemote, () => JSX.Element> = {
  waitingForYouToAccept: ReviewAccept,
  waitingForCounterpartyToAccept: WaitingAccept,
  waitingForYouToFinish: ReviewFinish,
  waitingForCounterpartyToFinish: WaitingFinish,
  swapTransactionPending: TxnPending,
  swapTransactionConfirmed: TxnConfirmed,
}

export function SwapStep() {
  const { isValidating, status } = useSwap()
  const { route: currentRoute } = usePathParams()

  if (!status) {
    return <Redirect to={routes.home} />
  }

  const nextRoute = swapStatusToRoute[status]

  if (!isValidating && currentRoute !== nextRoute) {
    return <Redirect to={routes[nextRoute]} />
  }

  const Component = status ? componentMap[status as SwapStatusRemote] : null

  if (!Component) {
    return null
  }

  return <Component />
}
