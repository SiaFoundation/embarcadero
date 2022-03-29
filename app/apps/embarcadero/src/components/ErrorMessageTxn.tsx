import { useSwap } from '../contexts/swap'
import { Message } from './Message'

export function ErrorMessageTxn() {
  const { txnError } = useSwap()

  if (!txnError) {
    return null
  }

  return <Message variant="error" message={txnError} />
}
