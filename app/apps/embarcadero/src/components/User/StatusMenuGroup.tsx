import {
  Code,
  DropdownMenuGroup,
  DropdownMenuLabel,
  Tooltip,
  Locked16,
  Unlocked16,
  Flex,
  Text,
  Paragraph,
} from '@siafoundation/design-system'
import { NetworkStatus } from '../NetworkStatus'
import { useConnectivity } from '../../hooks/useConnectivity'

export function StatusMenuGroup() {
  const { siad, embd, wallet } = useConnectivity()

  return (
    <DropdownMenuGroup>
      <DropdownMenuLabel>Status</DropdownMenuLabel>
      <Flex justify="between" gap="3" css={{ padding: '$1 $2' }}>
        <Tooltip
          content={
            wallet ? (
              <Flex css={{ padding: '$1' }}>
                <Text>Unlocked</Text>
              </Flex>
            ) : (
              <Flex css={{ padding: '$1', maxWidth: '400px' }}>
                <Paragraph size="1">
                  Locked - unlock <Code>siad</Code> wallet to use Embarcadero.
                </Paragraph>
              </Flex>
            )
          }
        >
          <Flex direction="column" gap="1" align="center">
            <Flex
              css={{
                color: wallet ? '$green10' : '$red10',
              }}
            >
              {wallet ? <Unlocked16 /> : <Locked16 />}
            </Flex>
            <Code variant="gray">wallet</Code>
          </Flex>
        </Tooltip>
        <Tooltip
          content={
            embd ? (
              <Flex css={{ padding: '$1' }}>
                <Text>Connected</Text>
              </Flex>
            ) : (
              <Flex css={{ padding: '$1' }}>
                <Text>Disconnected</Text>
              </Flex>
            )
          }
        >
          <Flex direction="column" gap="1" align="center">
            <NetworkStatus variant={embd ? 'green' : 'red'} />
            <Code variant="gray">embd</Code>
          </Flex>
        </Tooltip>
        <Tooltip
          content={
            siad ? (
              <Flex css={{ padding: '$1' }}>
                <Text>Connected</Text>
              </Flex>
            ) : (
              <Flex css={{ padding: '$1', maxWidth: '400px' }}>
                <Paragraph size="1">
                  Disconnected - you may need to configure <Code>embc</Code>{' '}
                  with the correct <Code>siad</Code> address. Do this with the{' '}
                  <Code>--siad</Code> flag.
                </Paragraph>
              </Flex>
            )
          }
        >
          <Flex direction="column" gap="1" align="center">
            <NetworkStatus variant={siad ? 'green' : 'red'} />
            <Code variant="gray">siad</Code>
          </Flex>
        </Tooltip>
      </Flex>
    </DropdownMenuGroup>
  )
}
