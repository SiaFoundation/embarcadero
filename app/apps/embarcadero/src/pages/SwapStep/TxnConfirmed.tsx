import { Flex, CheckmarkOutline16 } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { DownloadTxn } from '../../components/DownloadTxn'
import { Message } from '../../components/Message'

export function TxnConfirmed() {
  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex direction="column" align="center" gap="1-5">
        <Message
          variant="success"
          message={`
            The swap has been signed by both parties and is complete. Download the completed transaction.
          `}
        />
        <DownloadTxn />
      </Flex>
    </Flex>
  )
}
