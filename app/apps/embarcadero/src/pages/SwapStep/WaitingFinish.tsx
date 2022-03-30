import { Flex } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { DownloadTxn } from '../../components/DownloadTxn'
import { Message } from '../../components/Message'
import { ErrorMessageTxn } from '../../components/ErrorMessageTxn'

export function WaitingFinish() {
  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex direction="column" align="center" gap="1-5">
        <ErrorMessageTxn />
        <Message
          message={`
            To finalize the swap, download the transaction file and share it
            with your counterparty. Your counterparty can then sign and broadcast the completed transaction.
          `}
        />
        <DownloadTxn />
        <Message
          message={`
            Once the counterparty finalizes the swap transaction, Embarcadero will be able to confirm it on the network.
          `}
        />
      </Flex>
    </Flex>
  )
}
