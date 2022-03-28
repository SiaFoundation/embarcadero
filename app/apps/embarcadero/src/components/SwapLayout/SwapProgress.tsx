import { Box, Flex, ProgressBar, Text } from '@siafoundation/design-system'
import { capitalize, kebabCase } from 'lodash'
import { useSwap } from '../../contexts/swap'
import { SwapStatus } from '../../lib/swapStatus'

const statusToProgress = {
  createANewSwap: undefined,
  openASwap: undefined,
  waitingForYouToAccept: 2,
  waitingForCounterpartyToAccept: 2,
  waitingForYouToFinish: 3,
  waitingForCounterpartyToFinish: 0,
  swapTransactionPending: 0,
  swapTransactionConfirmed: 4,
} as Record<SwapStatus, number | undefined>

export function SwapProgress() {
  const { status } = useSwap()

  const step = status === undefined ? undefined : statusToProgress[status]

  return (
    <Flex direction="column" gap="3" css={{ width: '100%' }}>
      <Text size="20" weight="semibold">
        {capitalize(kebabCase(status).split('-').join(' '))}
      </Text>
      {step !== undefined && (
        <Box css={{ width: '100%' }}>
          <ProgressBar
            key={step}
            value={step ? step * 25 : undefined}
            variant={step ? 'gradient' : 'gray'}
          />
        </Box>
      )}
    </Flex>
  )
}
