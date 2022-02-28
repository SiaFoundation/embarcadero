import { Button, Flex } from '@siafoundation/design-system'
import { SwapOverview } from '../../components/SwapOverview'
import { useSwap } from '../../hooks/useSwap'
import { Redirect, useHistory } from 'react-router-dom'
import { Fragment, useCallback, useState } from 'react'
import axios from 'axios'
import { Message } from '../../components/Message'
import { Share } from './Share'

export function StepSwap() {
  const { raw, status, signTransaction, transactionError } = useSwap()

  if (!raw) {
    return <Redirect to="/" />
  }

  return (
    <Flex direction="column" align="center" gap="3">
      <SwapOverview />
      <Flex
        direction="column"
        align="center"
        gap="3"
        css={{ overflow: 'hidden', width: '100%' }}
      >
        {status === 'waitingForYouToAccept' && (
          <Fragment>
            <Message
              message={`
            Accept and sign the transaction to continue. After this, the counterparty can complete the transaction
        `}
            />
            <Button
              size="3"
              variant="green"
              css={{ width: '100%' }}
              onClick={() => signTransaction('accept')}
            >
              Accept and sign transaction
            </Button>
          </Fragment>
        )}
        {status === 'waitingForCounterpartyToAccept' && <Share />}
        {status === 'waitingForYouToFinish' && (
          <Fragment>
            <Message
              message={`
              Accept and sign the transaction to continue. After this, the counterparty can complete the transaction.
            `}
            />
            {transactionError && (
              <Message variant="red" message={'Error completing transaction'} />
            )}
            <Button
              size="3"
              variant="green"
              css={{ width: '100%' }}
              onClick={() => signTransaction('finish')}
            >
              Finish and broadcast transaction
            </Button>
          </Fragment>
        )}
        {status === 'waitingForCounterpartyToFinish' && <Share />}
      </Flex>
    </Flex>
  )
}
