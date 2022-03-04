import { Button, Flex } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { useSwap } from '../../contexts/swap'
import { Fragment } from 'react'
import { Message } from '../../components/Message'
import { useTxnHasBalance } from '../../hooks/useTxnHasBalance'
import { useConnectivity } from '../../hooks/useConnectivity'
import { ErrorMessageConn } from '../../components/ErrorMessageConn'
import { ErrorMessageTxn } from '../../components/ErrorMessageTxn'

export function ReviewFinish() {
  const { signTxn } = useSwap()
  const { all } = useConnectivity()
  const hasBalance = useTxnHasBalance()

  const readyToSign = all && hasBalance

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
                Sign and broadcast the transaction to complete the swap.
              `}
          />
          <ErrorMessageTxn />
          <ErrorMessageConn />
          <Button
            size="3"
            variant="green"
            css={{ width: '100%' }}
            disabled={!readyToSign}
            onClick={() => signTxn('finish')}
          >
            Sign and broadcast transaction
          </Button>
        </Fragment>
      </Flex>
    </Flex>
  )
}