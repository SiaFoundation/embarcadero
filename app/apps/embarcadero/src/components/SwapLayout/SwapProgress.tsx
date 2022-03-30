import { Box, Flex, ProgressBar, Text } from '@siafoundation/design-system'
import { capitalize, kebabCase } from 'lodash'
import { useSwap } from '../../contexts/swap'
import { SwapStatus } from '../../lib/swapStatus'

const statusToProgress = {
  waitingForYouToAccept: 2,
  waitingForCounterpartyToAccept: 2,
  waitingForYouToFinish: 3,
  swapTransactionConfirmed: 4,
} as Record<SwapStatus, number>

export function SwapProgress() {
  const { status, hasDownloaded } = useSwap()

  const step = status && statusToProgress[status]

  const pending =
    status === 'swapTransactionPending' ||
    (status === 'waitingForCounterpartyToFinish' && hasDownloaded)

  return (
    <Flex direction="column" gap="3" css={{ width: '100%' }}>
      <Text size="20" weight="semibold">
        {capitalize(kebabCase(status).split('-').join(' '))}
      </Text>
      {step && (
        <Box css={{ width: '100%' }}>
          <ProgressBar key={step} value={step * 25} variant="gradient" />
        </Box>
      )}
      {pending && (
        <Box css={{ width: '100%' }}>
          <ProgressBar variant="gray" />
        </Box>
      )}
    </Flex>
  )
}
