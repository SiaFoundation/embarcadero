import {
  Box,
  Flex,
  Label,
  Text,
  TextField,
  Tooltip,
} from '@siafoundation/design-system'
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
  tabIndex?: number
  onChange?: (value?: BigNumber) => void
}

export function Input({
  currency,
  isOffer,
  type,
  disabled = false,
  tabIndex,
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
    currency === 'SC' && value && scPrice ? scPrice * value.toNumber() : null

  return (
    <Flex
      direction="column"
      gap="1-5"
      css={{
        backgroundColor: '$brandGray3',
        border: '2px solid',
        borderColor: !hasAvailableBalance ? '$red8' : 'transparent',
        borderRadius: '$2',
        padding: '$1-5',
        transition: 'border 50ms linear',
        '&:hover': !disabled && {
          borderColor: !hasAvailableBalance ? '$red9' : '$brandAccent7',
        },
      }}
    >
      <Label css={{ color: '$brandGray10' }}>
        {isOffer ? 'You will send' : 'You will receive'}
      </Label>
      <Box
        css={{
          position: 'relative',
        }}
      >
        <TextField
          readOnly={disabled}
          tabIndex={tabIndex}
          size="3"
          variant="totalGhost"
          noSpin
          type="number"
          value={value !== undefined ? value.toString() : ''}
          onChange={handleChange}
          placeholder={type === 'integer' ? '0' : '0.0'}
          css={{
            px: 0,
            '&:disabled': {
              color: '$hiContrast',
            },
          }}
        />
        <Tooltip content={currency === 'SF' ? 'Siafunds' : 'Siacoins'}>
          <Flex
            align="center"
            css={{
              position: 'absolute',
              top: 0,
              right: '5px',
              backgroundColor: '$brandGray4',
              borderRadius: '$2',
              border: '1px solid $brandGray5',
              height: '100%',
              padding: '0 $1-5',
              ...(currency === 'SF' && {
                border: '1px dashed $orange8',
                backgroundColor: '$orange4',
              }),
            }}
          >
            <Text
              css={{
                fontWeight: 'bolder',
                ...(currency === 'SF' && {
                  color: '$orange12',
                }),
              }}
            >
              {currency}
            </Text>
          </Flex>
        </Tooltip>
      </Box>
      <Flex justify="between" css={{ pr: '$1' }}>
        {!hasAvailableBalance && <Text>Insufficient funds</Text>}
        {hasAvailableBalance && usdValue && (
          <Text>${usdValue.toLocaleString()} USD</Text>
        )}
        {usdValue && <Text>+ 5 SC fee</Text>}
      </Flex>
    </Flex>
  )
}
