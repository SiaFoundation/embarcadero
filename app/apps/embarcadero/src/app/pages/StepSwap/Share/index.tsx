import { ArrowRightIcon, DownloadIcon } from '@radix-ui/react-icons'
import { Box, Button, RLinkButton } from '@siafoundation/design-system'
import { Fragment } from 'react'
import { Message } from '../../../components/Message'
import { useSwap } from '../../../hooks/useSwap'
import { routes } from '../../../routes'

export function Share() {
  const { downloadTransaction } = useSwap()
  return (
    <Fragment>
      <Message
        message={`
          To proceed, download the partially completed transaction and share it
          with your counterparty for signing.
        `}
      />
      <Button
        onClick={() => downloadTransaction()}
        size="3"
        css={{ width: '100%' }}
      >
        Download transaction file
        <Box as="span" css={{ pl: '$1', lh: '1' }}>
          <DownloadIcon />
        </Box>
      </Button>
      <RLinkButton
        variant="green"
        size="3"
        to={routes.input}
        css={{
          width: '100%',
        }}
      >
        Load swap from counterparty
        <Box as="span" css={{ pl: '$1', lh: '1' }}>
          <ArrowRightIcon />
        </Box>
      </RLinkButton>
    </Fragment>
  )
}
