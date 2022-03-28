import { Box, Button, DocumentDownload16 } from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function DownloadTransaction() {
  const { downloadTxn } = useSwap()
  return (
    <Button
      onClick={() => downloadTxn()}
      size="3"
      variant="gray"
      css={{ width: '100%' }}
    >
      Download signed transaction
      <Box as="span" css={{ pl: '$1', lh: '1' }}>
        <DocumentDownload16 />
      </Box>
    </Button>
  )
}
