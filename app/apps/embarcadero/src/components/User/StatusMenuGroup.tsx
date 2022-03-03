import {
  Code,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRightSlot,
  Tooltip,
  Locked16,
  Unlocked16,
} from '@siafoundation/design-system'
import { NetworkStatus } from '../NetworkStatus'
import { useConnectivity } from '../../hooks/useConnectivity'

export function StatusMenuGroup() {
  const { siad, embd, wallet } = useConnectivity()

  return (
    <DropdownMenuGroup>
      <DropdownMenuLabel>Status</DropdownMenuLabel>
      <DropdownMenuItem disabled>
        Wallet
        <DropdownMenuRightSlot
          css={{
            '& > *, &:hover > *': {
              color: wallet ? '$hiContrast' : '$red10',
            },
          }}
        >
          <Tooltip content={wallet ? 'Unlocked' : 'Locked'}>
            {wallet ? <Unlocked16 /> : <Locked16 />}
          </Tooltip>
        </DropdownMenuRightSlot>
      </DropdownMenuItem>
      <DropdownMenuItem disabled>
        <Code variant="gray">embd</Code>
        <DropdownMenuRightSlot>
          <NetworkStatus
            variant={embd ? 'green' : 'red'}
            content={embd ? 'Connected' : 'Disconnected'}
          />
        </DropdownMenuRightSlot>
      </DropdownMenuItem>
      <DropdownMenuItem disabled>
        <Code variant="gray">siad</Code>
        <DropdownMenuRightSlot>
          <NetworkStatus
            variant={siad ? 'green' : 'red'}
            content={siad ? 'Connected' : 'Disconnected'}
          />
        </DropdownMenuRightSlot>
      </DropdownMenuItem>
    </DropdownMenuGroup>
  )
}
