import { Flex } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { DownloadTxn } from '../../components/DownloadTxn'
import { Message } from '../../components/Message'

export function TxnPending() {
  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex direction="column" align="center" gap="1-5">
        <DownloadTxn />
        <Message
          variant="info"
          message={`
            The unconfirmed transaction was found in the transaction pool. Waiting for blockchain confirmation.
          `}
        />
      </Flex>
    </Flex>
  )
}
