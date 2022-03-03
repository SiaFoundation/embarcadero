import { ArrowDownIcon } from '@radix-ui/react-icons'
import { Box, Flex } from '@siafoundation/design-system'
import { Input } from './Input'
import { useSwap } from '../contexts/swap'
import { Message } from './Message'

export function SwapOverview() {
  const { offerSc, sc, sf, transaction } = useSwap()

  const scInputs = transaction?.siacoinInputs?.length || 0
  const sfInputs = transaction?.siafundInputs?.length || 0
  const totalInputs = scInputs + sfInputs

  return (
    <Flex direction="column" gap="3" css={{ width: '100%' }}>
      <Flex direction="column" align="center" css={{ width: '100%' }}>
        <Box css={{ width: '100%', order: offerSc ? 1 : 3 }}>
          <Input
            currency="SC"
            type="decimal"
            disabled
            value={sc?.toString()}
            isOffer={offerSc}
          />
        </Box>
        <Box css={{ width: '100%', order: offerSc ? 3 : 1 }}>
          <Input
            currency="SF"
            type="integer"
            disabled
            value={sf?.toString()}
            isOffer={!offerSc}
          />
        </Box>
        <Box css={{ height: '$2', zIndex: '$1', order: 2 }}>
          <Box
            css={{
              position: 'relative',
              top: '-15px',
              height: '40px',
              width: '40px',
              backgroundColor: '$loContrast',
              borderRadius: '15px',
            }}
          >
            <Flex
              align="center"
              justify="center"
              css={{
                backgroundColor: '$gray4',
                borderRadius: '$4',
                position: 'absolute',
                transform: 'translate(-50%, -50%)',
                left: '50%',
                top: '50%',
                height: '30px',
                width: '30px',
              }}
            >
              <ArrowDownIcon />
            </Flex>
          </Box>
        </Box>
      </Flex>
      {totalInputs > 40 && (
        <Message
          variant="red"
          message={`
          Warning, this transactions has ${totalInputs} inputs, transactions with too many inputs may fail. Consider defragging your wallets if you run into issues.
        `}
        />
      )}
    </Flex>
  )
}
