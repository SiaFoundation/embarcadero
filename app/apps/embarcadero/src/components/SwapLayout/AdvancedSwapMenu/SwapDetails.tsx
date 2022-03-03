import { Codeblock, Text } from '@siafoundation/design-system'
import { useSwap } from '../../../contexts/swap'

export function SwapDetails() {
  const { summary } = useSwap()

  if (!summary) {
    return (
      <Text css={{ p: '$3 $2', color: '$slate9' }}>
        Load a swap to view details
      </Text>
    )
  }

  return (
    <pre>
      <Codeblock css={{ overflow: 'auto' }}>
        {JSON.stringify(summary, null, 2)}
      </Codeblock>
    </pre>
  )
}
