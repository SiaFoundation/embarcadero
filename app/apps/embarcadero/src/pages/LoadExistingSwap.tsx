import { Code, Dropzone, Flex, Text } from '@siafoundation/design-system'
import { DownloadTransaction } from '../components/DownloadTransaction'
import { Message } from '../components/Message'
import { useRouteToStep } from '../hooks/useRouteToStep'
import { useSwap } from '../contexts/swap'

export function LoadExistingSwap() {
  const { loadTransactionFromFile } = useSwap()

  useRouteToStep()

  return (
    <Flex direction="column" align="center" gap="3">
      <Message
        message={`
          Retrieve swap transaction file from your counterparty and open it to continue.
      `}
      />
      <Dropzone
        title={
          (
            <Text>
              Drop your <Code>transaction.txt</Code> here or click to open the
              file picker.{' '}
            </Text>
          ) as any
        }
        onFiles={(files) => loadTransactionFromFile(files[0])}
      />
    </Flex>
  )
}
