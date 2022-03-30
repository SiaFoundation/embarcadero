import { Button } from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function DownloadTxn() {
  const { downloadTxn } = useSwap()
  return (
    <Button
      onClick={() => downloadTxn()}
      size="3"
      variant="gray"
      css={{ width: '100%' }}
    >
      Download signed transaction
    </Button>
  )
}
