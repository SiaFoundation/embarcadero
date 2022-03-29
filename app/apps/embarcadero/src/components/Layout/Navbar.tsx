import {
  AppBar,
  Box,
  Container,
  Flex,
  Heading,
  Logo,
  SimpleLogoIcon,
  Text,
} from '@siafoundation/design-system'
import { User } from '../User'

export function Navbar() {
  return (
    <AppBar size="2" color="none" sticky>
      <Container size="4">
        <Flex align="center" gap="2">
          <Flex
            align="center"
            gap="1-5"
            css={{
              position: 'relative',
              top: '-1px',
            }}
          >
            <Box
              css={{
                transform: 'scale(1.4)',
              }}
            >
              <Logo />
            </Box>
            <Heading
              size="1"
              css={{
                fontWeight: '600',
                display: 'none',
                '@bp1': {
                  display: 'block',
                },
              }}
            >
              Embarcadero
            </Heading>
          </Flex>
          <Box css={{ flex: 1 }} />
          <User />
        </Flex>
      </Container>
    </AppBar>
  )
}
