import { ReviewAccept } from './ReviewAccept'
import { ReviewFinish } from './ReviewFinish'
import { WaitingAccept } from './WaitingAccept'
import { WaitingFinish } from './WaitingFinish'
import { TransactionComplete } from './TransactionComplete'
import { SwapStatusRemote } from '../../lib/swapStatus'
import { useSwap } from '../../contexts/swap'
import { Redirect } from 'react-router-dom'
import { swapStatusToRoute, routes } from '../../routes'
import { usePathParams } from '../../hooks/usePathParams'

const componentMap: Record<SwapStatusRemote, () => JSX.Element> = {
  waitingForYouToAccept: ReviewAccept,
  waitingForCounterpartyToAccept: WaitingAccept,
  waitingForCounterpartyToFinish: WaitingFinish,
  waitingForYouToFinish: ReviewFinish,
  transactionComplete: TransactionComplete,
}

export function SwapStep() {
  const { status } = useSwap()
  const { route: currentRoute } = usePathParams()

  if (!status) {
    return null
  }

  const nextRoute = swapStatusToRoute[status]

  if (currentRoute !== nextRoute) {
    return <Redirect to={routes[nextRoute]} />
  }

  const Component = status ? componentMap[status as SwapStatusRemote] : null

  if (!Component) {
    return null
  }

  return <Component />
}
