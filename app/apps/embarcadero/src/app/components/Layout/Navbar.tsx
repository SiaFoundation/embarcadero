import {
  AppBar,
  Box,
  Container,
  Flex,
  Heading,
  Logo,
  Text,
} from '@siafoundation/design-system'
import { toHumanReadable } from '@siafoundation/sia-js'
import { useSwap } from '../../hooks/useSwap'
import { Wallet } from './Wallet'

export function Navbar() {
  const { raw } = useSwap()
  return (
    <AppBar size="3" color="none" sticky>
      <Container size="4" css={{ position: 'relative' }}>
        <Flex align="center" gap="1" css={{}}>
          <Logo />
          <Heading
            css={{
              color: '$siaGreenA12',
              display: 'inline',
              // fontStyle: 'oblique',
              fontWeight: '600',
            }}
          >
            Embarcadero
          </Heading>
          <Text>
            {raw ? `${(raw.length / 1024 / 1024).toFixed(2)}mb` : 'none'}
          </Text>
          <Box css={{ flex: 1 }} />
          <Wallet />
        </Flex>
      </Container>
    </AppBar>
  )
}
