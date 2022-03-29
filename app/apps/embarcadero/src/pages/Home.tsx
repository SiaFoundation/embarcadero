import {
  Button,
  ControlGroup,
  Flex,
  Paragraph,
  RRLink,
  RRLinkButton,
  Separator,
  Text,
} from '@siafoundation/design-system'
import { useEffect } from 'react'
import { useHistory } from 'react-router-dom'
import { useSwap } from '../contexts/swap'
import { routes } from '../routes'

export function Home() {
  const history = useHistory()
  const { resetTxn } = useSwap()
  useEffect(() => {
    resetTxn()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return (
    <Flex direction="column" gap="3">
      <Flex direction="column" gap="2">
        <Paragraph>
          Welcome to Embarcadero, a tool for conducting escrowless SC/SF swaps.
        </Paragraph>
        <ControlGroup>
          <Button
            size="2"
            css={{ flex: 1 }}
            onClick={() => history.push(routes.input)}
          >
            Open a swap →
          </Button>
          <Button
            size="2"
            variant="accent"
            css={{ flex: 1 }}
            onClick={() => history.push(routes.create)}
          >
            Create a new swap →
          </Button>
        </ControlGroup>
        <Separator size="100" pad="0" />
      </Flex>
      <Flex direction="column" gap="2-5">
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
