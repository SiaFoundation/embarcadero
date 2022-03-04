import { Code, Dropzone, Flex } from '@siafoundation/design-system'
import { Fragment } from 'react'
import { useSwap } from '../contexts/swap'
import { Message } from './Message'

export function SwapDropzone() {
  const { fileReadError, loadTxnFromFile } = useSwap()

  return (
    <Flex direction="column" gap="3">
      {fileReadError && <Message variant="red" message={fileReadError} />}
      <Dropzone
        title={
          <Fragment>
            Drop your <Code>embc_txn.json</Code> here or click to open the file
            picker.{' '}
          </Fragment>
        }
        onFiles={(files) => loadTxnFromFile(files[0])}
      />
    </Flex>
  )
}
