import { Box, Flex, Heading, ProgressBar } from '@siafoundation/design-system'
import { capitalize, kebabCase } from 'lodash'
import { useSwap } from '../../hooks/useSwap'
import { SwapStatus } from '../../lib/swapStatus'

const stepMap = {
  creatingANewSwap: 0,
  loadingAnExistingSwap: 0,
  waitingForYouToAccept: 1,
  waitingForCounterpartyToAccept: 2,
  waitingForCounterpartyToFinish: 3,
  waitingForYouToFinish: 4,
} as Record<SwapStatus, number>

export function SwapProgress() {
  const { status } = useSwap()

  const step = status && stepMap[status]

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
