import { Flex } from '@siafoundation/design-system'
import { Message } from '../components/Message'
import { SwapDropzone } from '../components/SwapDropzone'

export function LoadExistingSwap() {
  return (
    <Flex direction="column" align="center" gap="3">
      <Message
        message={`
          Retrieve a swap transaction file from your counterparty and open it to begin.
      `}
      />
      <SwapDropzone />
    </Flex>
  )
}
