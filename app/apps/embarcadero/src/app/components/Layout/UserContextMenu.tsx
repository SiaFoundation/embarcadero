import {
  DotsHorizontalIcon,
  ExclamationTriangleIcon,
} from '@radix-ui/react-icons'
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@siafoundation/design-system'
import { StatusMenu } from './StatusMenu'
import { GeneralMenu } from './GeneralMenu'
import { ThemeMenu } from './ThemeMenu'
import { useConnectivity } from '../../hooks/useConnectivity'

type Props = React.ComponentProps<typeof Button>

export function UserContextMenu(props: Props) {
  const { all } = useConnectivity()

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button {...props} css={{ color: all ? '$hiContrast' : '$red10' }}>
          {all ? <DotsHorizontalIcon /> : <ExclamationTriangleIcon />}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <StatusMenu />
        <DropdownMenuSeparator />
        <GeneralMenu />
        <DropdownMenuSeparator />
        <ThemeMenu />
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
