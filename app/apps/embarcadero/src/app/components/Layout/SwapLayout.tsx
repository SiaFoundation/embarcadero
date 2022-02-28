import { Box, Container, Flex, Grid, Panel } from '@siafoundation/design-system'
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
        padding: '150px 0',
      }}
    >
      <Panel size="3">
        <Flex
          direction="column"
          gap="4"
          justify="center"
          css={{
            padding: '$3',
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
