import { Flex, Paragraph, RLink, Separator } from '@siafoundation/design-system'
import { useEffect } from 'react'
import { useSwap } from '../contexts/swap'
import { routes } from '../routes'

export function Home() {
  const { clearTransaction } = useSwap()
  useEffect(() => {
    clearTransaction()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <Flex direction="column" gap="2">
      <Flex direction="column" gap="4">
        <Paragraph>
          Welcome to Embarcadero, a tool for conducting escrowless SC/SF swaps.
        </Paragraph>
        <Flex direction="column" gap="4" align="start" css={{ mb: '$2' }}>
          <RLink to={routes.create}>Create a new swap transaction →</RLink>
          <RLink to={routes.input}>Load an existing swap transaction →</RLink>
        </Flex>
      </Flex>
      <Separator size="3" />
      <Flex direction="column" gap="4">
        <Paragraph size="1">
          Executing a swap is a three-part process. Here, we assume that the
          swappers, Alice and Bob, have established a communication channel and
          have negotiated the terms of the swap (in this case, 7SF for 10MS).
        </Paragraph>
        <Paragraph size="1">
          1. Alice begins by creating a transaction with two outputs: one worth
          7SF, and one worth 10MS. She adds inputs from her wallet worth 7SF,
          but does not sign anything yet. She sends this partially-completed
          transaction to Bob.
        </Paragraph>
        <Paragraph size="1">
          2. Bob reviews the transaction and confirms that its outputs are
          correct. He then adds inputs from his wallet worth 10MS, signs the
          transaction, and returns it to Alice.
        </Paragraph>
        <Paragraph size="1">
          3. Alice reviews the transaction (in case Bob did something
          malicious). Assuming all is well, she adds her signatures. The
          transaction is now complete, so Alice broadcasts it to the Sia network
          for inclusion in the next block.
        </Paragraph>
      </Flex>
    </Flex>
  )
}
