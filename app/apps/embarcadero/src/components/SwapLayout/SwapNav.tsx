import { Flex, RLink } from '@siafoundation/design-system'
import { AdvancedSwapMenu } from './AdvancedSwapMenu'

export function SwapNav() {
  return (
    <Flex gap="1" justify="between" align="center">
      <RLink to="/" css={{ fontSize: '$6' }}>
        Swap
      </RLink>
      <AdvancedSwapMenu />
    </Flex>
  )
}
