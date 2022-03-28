import {
  Box,
  Flex,
  Heading,
  Home24,
  IconButton,
  RRLink,
  Text,
} from '@siafoundation/design-system'
import { AdvancedSwapMenu } from './AdvancedSwapMenu'

export function SwapNav() {
  return (
    <Flex gap="1" justify="between" align="center">
      <RRLink to="/" css={{ textDecoration: 'none' }}>
        <Flex gap="1" align="center">
          <Home24 />
          <Text
            size="20"
            weight="semibold"
            css={{ position: 'relative', top: '2px' }}
          >
            Swap
          </Text>
        </Flex>
      </RRLink>
      <Box css={{ position: 'relative', right: '-$0-5', top: '-1px' }}>
        <AdvancedSwapMenu />
      </Box>
    </Flex>
  )
}
