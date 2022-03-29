import {
  Box,
  Button,
  ControlGroup,
  Copy24,
  copyToClipboard,
  IconButton,
  Tooltip,
} from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function DownloadTransaction() {
  const { id, downloadTxn } = useSwap()
  return (
    <ControlGroup css={{ width: '100%' }}>
      <Button
        onClick={() => downloadTxn()}
        size="3"
        variant="gray"
        css={{ flex: 1 }}
      >
        Download signed transaction
      </Button>
      <Tooltip content="Copy transaction ID">
        <IconButton
          variant="gray"
          size="3"
          css={{ borderRadius: '$2' }}
          onClick={() => id && copyToClipboard(id, 'transaction ID')}
        >
          <Copy24 />
        </IconButton>
      </Tooltip>
    </ControlGroup>
  )
}
