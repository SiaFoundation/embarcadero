import { Flex, Number_132, Number_232 } from '@siafoundation/design-system'
import { SwapOverview } from '../components/SwapOverview'
import { DownloadTransaction } from '../components/DownloadTransaction'
import { Message } from '../components/Message'
import { useRouteToStep } from '../hooks/useRouteToStep'
import { useProtectSwapRoute } from '../hooks/useProtectSwapRoute'
import { SwapDropzone } from '../components/SwapDropzone'

export function WaitingFinish() {
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
          icon={<Number_132 />}
          message={`
            To finish the swap, download the transaction file and share it
            with your counterparty for signing.
          `}
        />
        <DownloadTransaction />
        <Message
          icon={<Number_232 />}
          message={`
            Retrieve the signed transaction file from your counterparty and open it to view the completed transaction.
        `}
        />
        <SwapDropzone />
      </Flex>
    </Flex>
  )
}
