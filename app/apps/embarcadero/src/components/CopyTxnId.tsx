import {
  Copy24,
  copyToClipboard,
  IconButton,
  Tooltip,
} from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function CopyTxnId() {
  const { id } = useSwap()

  return (
    <Tooltip content="Copy transaction ID">
      <IconButton
        size="2"
        disabled={!id}
        onClick={() => id && copyToClipboard(id, 'transaction ID')}
      >
        <Copy24 />
      </IconButton>
    </Tooltip>
  )
}
