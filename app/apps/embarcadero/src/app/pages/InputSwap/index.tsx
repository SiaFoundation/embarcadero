import { Flex } from '@siafoundation/design-system'
import { Message } from '../../components/Message'
import { useSwap } from '../../hooks/useSwap'
import { Dropzone } from './Dropzone'

export function InputSwap() {
  const { loadTransactionFromFile } = useSwap()

  return (
    <Flex direction="column" align="center" gap="3">
      <Message
        message={`
          Retrieve a swap transaction file from your counterparty and open it to view details.
      `}
      />
      <Dropzone onFiles={(files) => loadTransactionFromFile(files[0])} />
    </Flex>
  )
}
