import { Button, Flex } from '@siafoundation/design-system'
import { SwapOverview } from '../components/SwapOverview'
import { useSwap } from '../contexts/swap'
import { Fragment } from 'react'
import { Message } from '../components/Message'
import { useRouteToStep } from '../hooks/useRouteToStep'
import { useProtectSwapRoute } from '../hooks/useProtectSwapRoute'

export function ReviewAccept() {
  const { signTransaction, transactionError } = useSwap()

  useRouteToStep()
  useProtectSwapRoute()

  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex
        direction="column"
        align="center"
        gap="3"
        css={{ overflow: 'hidden', width: '100%' }}
      >
        <Fragment>
          <Message
            message={`
            Accept and sign the transaction to continue. After this, the counterparty can complete the transaction
        `}
          />
          {transactionError && (
            <Message variant="red" message={'Error accepting transaction'} />
          )}
          <Button
            size="3"
            variant="green"
            css={{ width: '100%' }}
            onClick={() => signTransaction('accept')}
          >
            Accept and sign transaction
          </Button>
        </Fragment>
      </Flex>
    </Flex>
  )
}
