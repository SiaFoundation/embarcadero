import { Box, Flex, Label, Text, TextField } from '@siafoundation/design-system'
import { useSiaStatsNetworkStatus } from '@siafoundation/sia-react'
import { useSettings } from '../hooks/useSettings'
import { useHasBalance } from '../hooks/useHasBalance'
import BigNumber from 'bignumber.js'
import { useCallback } from 'react'

type Props = {
  disabled?: boolean
  isOffer?: boolean
  currency: 'SF' | 'SC'
  type: 'integer' | 'decimal'
  value?: BigNumber
  onChange?: (value?: BigNumber) => void
}

export function Input({
  currency,
  isOffer,
  type,
  disabled = false,
  value,
  onChange,
}: Props) {
  const { settings } = useSettings()
  const { data } = useSiaStatsNetworkStatus({
    disabled: !settings.siaStats,
  })
  const scPrice = settings.siaStats && data?.coin_price_USD

  const hasAvailableBalance = useHasBalance({
    value,
    isOffer,
    currency,
  })

  const handleChange = useCallback(
    (e) => {
      if (!onChange) {
        return
      }

      const { value } = e.target
      const num = value ? new BigNumber(value) : undefined
      onChange(num)
    },
    [onChange]
  )

  const usdValue =
    currency === 'SC' && Number(value) && scPrice
      ? scPrice * Number(value)
      : null

  return (
    <Flex
      direction="column"
      gap="2"
      css={{
        backgroundColor: '$gray3',
        border: '2px solid',
        borderColor: !hasAvailableBalance ? '$red8' : 'transparent',
        borderRadius: '$2',
        padding: '$2',
        transition: 'border 50ms linear',
        '&:hover': !disabled && {
          borderColor: !hasAvailableBalance ? '$red9' : '$primary6',
        },
      }}
    >
      <Label css={{ color: '$gray10' }}>
        {isOffer ? 'You are offering' : 'You are receiving'}
      </Label>
      <Box
        css={{
          position: 'relative',
        }}
      >
        <TextField
          disabled={disabled}
          size="3"
          variant="totalGhost"
          noSpin
          type="number"
          value={value !== undefined ? value.toString() : ''}
          onChange={handleChange}
          placeholder={type === 'integer' ? '0' : '0.0'}
          css={{
            '&:disabled': {
              color: '$hiContrast',
            },
          }}
        />
        <Flex
          align="center"
          css={{
            position: 'absolute',
            top: 0,
            right: '5px',
            backgroundColor: '$gray4',
            borderRadius: '$2',
            border: '1px solid $gray5',
            height: '100%',
            padding: '$1 $2',
          }}
        >
          <Text css={{ fontWeight: 'bolder' }}>{currency}</Text>
        </Flex>
      </Box>
      <Flex justify="between" css={{ pr: '$1' }}>
        {!hasAvailableBalance && <Text>Insufficient funds</Text>}
        {hasAvailableBalance && usdValue && (
          <Text>${usdValue.toLocaleString()} USD</Text>
        )}
        {usdValue && <Text css={{ color: '$gray10' }}>+ 5 SC fee</Text>}
      </Flex>
    </Flex>
  )
}
