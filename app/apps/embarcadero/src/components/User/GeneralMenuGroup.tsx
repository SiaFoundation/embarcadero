import {
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuRightSlot,
  LogoDiscord16,
  Notebook16,
  LicenseGlobal16,
  Information16,
  Link,
} from '@siafoundation/design-system'
import { useDialog } from '../../contexts/dialog'

export function GeneralMenuGroup() {
  const { openDialog } = useDialog()

  return (
    <DropdownMenuGroup>
      <DropdownMenuLabel>General</DropdownMenuLabel>
      <DropdownMenuItem>
        <Link
          href="https://github.com/SiaFoundation/embarcadero"
          target="_blank"
          css={{
            display: 'flex',
            width: '100%',
            alignItems: 'center',
          }}
        >
          About
          <DropdownMenuRightSlot>
            <Information16 />
          </DropdownMenuRightSlot>
        </Link>
      </DropdownMenuItem>
      <DropdownMenuItem>
        <Link
          href="https://discord.gg/sia"
          target="_blank"
          css={{
            display: 'flex',
            width: '100%',
            alignItems: 'center',
          }}
        >
          Discord
          <DropdownMenuRightSlot>
            <LogoDiscord16 />
          </DropdownMenuRightSlot>
        </Link>
      </DropdownMenuItem>
      <DropdownMenuItem>
        <Link
          href="https://support.sia.tech"
          target="_blank"
          css={{
            display: 'flex',
            width: '100%',
            alignItems: 'center',
          }}
        >
          Docs
          <DropdownMenuRightSlot>
            <Notebook16 />
          </DropdownMenuRightSlot>
        </Link>
      </DropdownMenuItem>
      <DropdownMenuItem onSelect={() => openDialog('privacy')}>
        Privacy
        <DropdownMenuRightSlot>
          <LicenseGlobal16 />
        </DropdownMenuRightSlot>
      </DropdownMenuItem>
    </DropdownMenuGroup>
  )
}
