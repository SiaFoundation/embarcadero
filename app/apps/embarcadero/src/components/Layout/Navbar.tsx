import {
  AppBar,
  Box,
  Container,
  Flex,
  Heading,
  Logo,
} from '@siafoundation/design-system'
import { User } from '../User'

export function Navbar() {
  return (
    <AppBar size="3" color="none" sticky>
      <Container size="4" css={{ position: 'relative' }}>
        <Flex align="center" gap="1">
          <Box css={{ position: 'relative', top: '-2px' }}>
            <Logo />
          </Box>
          <Heading
            css={{
              color: '$primary12',
              display: 'inline',
              fontWeight: '600',
            }}
          >
            Embarcadero
          </Heading>
          <Box css={{ flex: 1 }} />
          <User />
        </Flex>
      </Container>
    </AppBar>
  )
}
