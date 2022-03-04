import { Flex, Number_132, Number_232 } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { DownloadTransaction } from '../../components/DownloadTransaction'
import { Message } from '../../components/Message'
import { SwapDropzone } from '../../components/SwapDropzone'

export function WaitingFinish() {
  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex
        direction="column"
        align="center"
        gap="3"
        css={{ overflow: 'hidden', width: '100%' }}
      >
        <Message
          icon={<Number_132 />}
          message={`
            To finalize the swap, download the transaction file and share it
            with your counterparty. Your counterparty can then sign and broadcast the completed transaction.
          `}
        />
        <DownloadTransaction />
        <Message
          icon={<Number_232 />}
          message={`
            Retrieve the signed transaction file from your counterparty and open it to view the completed transaction.
        `}
        />
        <SwapDropzone />
      </Flex>
    </Flex>
  )
}