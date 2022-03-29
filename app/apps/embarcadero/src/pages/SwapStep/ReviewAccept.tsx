import {
  Button,
  Flex,
  NextOutline16,
  NextOutline32,
  Number_132,
} from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { useSwap } from '../../contexts/swap'
import { Fragment } from 'react'
import { Message } from '../../components/Message'
import { useTxnHasBalance } from '../../hooks/useTxnHasBalance'
import { useConnectivity } from '../../hooks/useConnectivity'
import { ErrorMessageConn } from '../../components/ErrorMessageConn'
import { ErrorMessageTxn } from '../../components/ErrorMessageTxn'

export function ReviewAccept() {
  const { signTxn } = useSwap()
  const { all } = useConnectivity()
  const hasBalance = useTxnHasBalance()

  const readyToSign = all && hasBalance

  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex direction="column" align="center" gap="1-5">
        <Fragment>
          <Message
            message={`
            Accept and sign the transaction to continue. After this, the counterparty can complete the transaction
        `}
          />
          <ErrorMessageTxn />
          <ErrorMessageConn />
          <Button
            size="3"
            variant="accent"
            css={{ width: '100%' }}
            disabled={!readyToSign}
            onClick={() => signTxn('accept')}
          >
            Accept and sign transaction
          </Button>
        </Fragment>
      </Flex>
    </Flex>
  )
}
