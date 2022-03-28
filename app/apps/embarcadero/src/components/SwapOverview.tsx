import { ArrowDown16, Box, Flex } from '@siafoundation/design-system'
import { Input } from './Input'
import { useSwap } from '../contexts/swap'
import { Message } from './Message'
import { ToggleInputs } from './ToggleInputs'

export function SwapOverview() {
  const { offerSc, sc, sf, txn } = useSwap()

  const scInputs = txn?.siacoinInputs?.length || 0
  const sfInputs = txn?.siafundInputs?.length || 0
  const totalInputs = scInputs + sfInputs

  return (
    <Flex direction="column" gap="3" css={{ width: '100%' }}>
      <Flex direction="column" align="center" css={{ width: '100%' }}>
        <Box css={{ width: '100%', order: offerSc ? 1 : 3 }}>
          <Input
            currency="SC"
            tabIndex={offerSc ? 1 : 3}
            type="decimal"
            disabled
            value={sc}
            isOffer={offerSc}
          />
        </Box>
        <Box css={{ width: '100%', order: offerSc ? 3 : 1 }}>
          <Input
            currency="SF"
            tabIndex={offerSc ? 3 : 1}
            type="integer"
            disabled
            value={sf}
            isOffer={!offerSc}
          />
        </Box>
        <ToggleInputs disabled />
      </Flex>
      {totalInputs > 40 && (
        <Message
          variant="error"
          message={`
          Warning, this transactions has ${totalInputs} inputs, transactions with too many inputs may fail. Consider defragging your wallets if you run into issues.
        `}
        />
      )}
    </Flex>
  )
}
