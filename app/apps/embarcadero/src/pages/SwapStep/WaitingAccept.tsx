import {
  Flex,
  NextOutline16,
  Number_132,
  Number_232,
} from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { Message } from '../../components/Message'
import { SwapDropzone } from '../../components/SwapDropzone'
import { DownloadTxn } from '../../components/DownloadTxn'

export function WaitingAccept() {
  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex direction="column" align="center" gap="1-5">
        <Message
          message={`
            To proceed, download the transaction file and share it
            with your counterparty for signing.
          `}
        />
        <DownloadTxn />
        <Message
          message={`
            Retrieve the signed transaction file from your counterparty and open it to continue.
          `}
        />
        <SwapDropzone />
      </Flex>
    </Flex>
  )
}
