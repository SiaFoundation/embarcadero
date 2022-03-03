import { Code, Dropzone, Paragraph } from '@siafoundation/design-system'
import { useSwap } from '../contexts/swap'

export function SwapDropzone() {
  const { loadTransactionFromFile } = useSwap()

  return (
    <Dropzone
      title={
        <Paragraph>
          Drop your <Code>transaction.txt</Code> here or click to open the file
          picker.{' '}
        </Paragraph>
      }
      onFiles={(files) => loadTransactionFromFile(files[0])}
    />
  )
}
