import { Box, Container, Flex, Panel } from '@siafoundation/design-system'
import React from 'react'
import { SwapNav } from './SwapNav'
import { SwapProgress } from './SwapProgress'

type Props = {
  children: React.ReactNode
}

export function SwapLayout({ children }: Props) {
  return (
    <Container
      size="1"
      css={{
        py: '20px',
      }}
    >
      <Panel
        css={{
          backgroundColor: '$loContrast',
          borderRadius: '$3',
          width: '100%',
        }}
      >
        <Flex
          direction="column"
          gap="2"
          justify="center"
          css={{
            padding: '$1-5 $2 $2 $2',
          }}
        >
          <SwapNav />
          <SwapProgress />
          <Box>{children}</Box>
        </Flex>
      </Panel>
    </Container>
  )
}
