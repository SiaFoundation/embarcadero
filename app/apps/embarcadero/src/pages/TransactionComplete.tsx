import { Flex, CheckmarkOutline16 } from '@siafoundation/design-system'
import { SwapOverview } from '../components/SwapOverview'
import { DownloadTransaction } from '../components/DownloadTransaction'
import { Message } from '../components/Message'
import { useRouteToStep } from '../hooks/useRouteToStep'
import { useProtectSwapRoute } from '../hooks/useProtectSwapRoute'

export function TransactionComplete() {
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
        <Message
          icon={<CheckmarkOutline16 />}
          message={`
            The swap has been signed by both parties and is complete. Download the completed transaction.
          `}
        />
        <DownloadTransaction />
      </Flex>
    </Flex>
  )
}
