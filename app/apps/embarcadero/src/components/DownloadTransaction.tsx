import { Box, Button, DocumentDownload16 } from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function DownloadTransaction() {
  const { downloadTxn } = useSwap()
  return (
    <Button
      onClick={() => downloadTxn()}
      size="3"
      variant="green"
      css={{ width: '100%' }}
    >
      Download transaction file
      <Box as="span" css={{ pl: '$1', lh: '1' }}>
        <DocumentDownload16 />
      </Box>
    </Button>
  )
}
