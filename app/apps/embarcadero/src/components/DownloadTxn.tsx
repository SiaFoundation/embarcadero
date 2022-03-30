import { Button } from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'
import { SwapStatus } from '../lib/swapStatus'

const txnIsComplete = [
  'waitingForCounterpartyToFinish',
  'waitingForYouToFinish',
  'swapTransactionConfirmed',
  'swapTransactionPending',
] as SwapStatus[]

export function DownloadTxn() {
  const { status, downloadTxn } = useSwap()

  let message = 'Download incomplete transaction'

  if (status && txnIsComplete.includes(status)) {
    message = 'Download signed transaction'
  }

  return (
    <Button
      onClick={() => downloadTxn()}
      size="3"
      variant="gray"
      css={{ width: '100%' }}
    >
      {message}
    </Button>
  )
}
