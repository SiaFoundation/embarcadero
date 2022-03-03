import { Box, Flex, Heading, ProgressBar } from '@siafoundation/design-system'
import { capitalize, kebabCase } from 'lodash'
import { useSwap } from '../../contexts/swap'
import { SwapStatus } from '../../lib/swapStatus'

const statusToStep = {
  creatingANewSwap: 0,
  loadingAnExistingSwap: 0,
  waitingForYouToAccept: 2,
  waitingForCounterpartyToAccept: 2,
  waitingForCounterpartyToFinish: 3,
  waitingForYouToFinish: 3,
  transactionComplete: 4,
} as Record<SwapStatus, number>

export function SwapProgress() {
  const { status } = useSwap()

  const step = status && statusToStep[status]

  return (
    <Flex direction="column" gap="3" css={{ width: '100%' }}>
      <Heading>{capitalize(kebabCase(status).split('-').join(' '))}</Heading>
      {step !== undefined && (
        <Box>
          <ProgressBar
            key={step}
            value={step ? step * 25 : undefined}
            variant="gradient"
          />
        </Box>
      )}
    </Flex>
  )
}
